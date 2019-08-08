package server

import (
	"gopkg.in/yaml.v2"

	log "github.com/franekjel/sokserver/logger"
	//"github.com/franekjel/sokserver/users"
)

//Command holds user command data (unmarshallized from yaml)
type Command struct {
	Login    string `yaml:"login"`    //user login
	Password string `yaml:"password"` //user password - deleted after check
	Command  string `yaml:"command"`  //given command - described in protocol.md
	Contest  string `yaml:"contest"`  //contest name for commands that need it
	Round    string `yaml:"round"`    //round name for commands that need it
	Task     string `yaml:"task"`     //task name for commands that need it
	Code     string `yaml:"code"`     //solution code if it is submission
}

//Execute given command. Return response to the client
func (s *Server) Execute(buff []byte) []byte {
	var com Command
	err := yaml.Unmarshal(buff, com)
	if err != nil {
		log.Error("Error parsing command")
		return []byte{}
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
