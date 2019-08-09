package server

import (
	"github.com/franekjel/sokserver/fs"
	"gopkg.in/yaml.v2"
	"regexp"

	log "github.com/franekjel/sokserver/logger"
	"github.com/franekjel/sokserver/users"
)

//Command holds user command data (unmarshallized from yaml)
type Command struct {
	Login    string `yaml:"login"`    //user login
	Password string `yaml:"password"` //user password - deleted after check
	Command  string `yaml:"command"`  //given command - described in protocol.md
	Contest  string `yaml:"contest"`  //contest name for commands that need it
	Round    string `yaml:"round"`    //round name for commands that need it
	Task     string `yaml:"task"`     //task name for commands that need it
	Data     string `yaml:"data"`     //additional data like submission code (specified in docs)
}

//ReturnMessage send to client after execute command
type ReturnMessage struct {
	Status string `yaml:status` //status may be "ok" or contains error
}

//Execute given command. Return response to the client
func (s *Server) Execute(buff []byte) []byte {
	var com Command
	err := yaml.Unmarshal(buff, com)
	if err != nil {
		log.Error("Error parsing command")
		return []byte{}
	}
	if com.Command == "create_account" { //one case when we don't check password
		return s.createAccount(&com)
	}

	if !s.verifyUser(&com) {
		return []byte{}
	}
	log.Info("User %s execute %s", com.Login, com.Command)

	return []byte{}
}

func (s *Server) verifyUser(com *Command) bool {

	user, ok := s.users[com.Login]
	if !ok {
		log.Error("User %s doesn't exist", com.Login)
		return false
	}
	if !user.VerifyPassword([]byte(com.Password)) {
		log.Error("Bad password for user %s", com.Login)
		return false
	}

	return true
}

func returnStatus(msg string) []byte {
	buff, _ := yaml.Marshal(ReturnMessage{msg})
	return buff
}

func (s *Server) createAccount(com *Command) []byte {
	if _, ok := s.users[com.Login]; ok {
		return returnStatus("Cannot create account, there is already user with this login")
	}
	if ok, _ := regexp.Match(`^([a-z][A-Z][0-9])$`, []byte(com.Login)); !ok {
		return returnStatus("Login may contains only letters and numbers")
	}
	dir := fs.Fs{Path: fs.Join(s.fs.Path, "users")}
	user := users.AddUser(&dir, &com.Login, []byte(com.Password))
	if user == nil {
		return returnStatus("Error creating user (probably during hashing password)")
	}
	s.users[com.Login] = user
	return returnStatus("ok")
}
