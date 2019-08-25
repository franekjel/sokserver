package server

import (
	"github.com/franekjel/sokserver/fs"
	"gopkg.in/yaml.v2"
	"regexp"
	"strconv"
	"time"

	log "github.com/franekjel/sokserver/logger"
	"github.com/franekjel/sokserver/submissions"
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
	Status         string          `yaml:"status"`                       //status may be "ok" or contains error
	ContestRanking map[string]uint `yaml:"contest_ranking,omitempty"`    //used in contest_ranking
	Tasks          []string        `yaml:"tasks,omitempty"`              //used in round_ranking
	Users          []string        `yaml:"users,omitempty"`              //used in round_ranking
	RoundRanking   [][]uint        `yaml:"round_ranking,omitempty,flow"` //used in round_ranking
	Filename       string          `yaml:"filename,omitempty"`           //used in get_task
	Data           string          `yaml:"data,omitempty"`               //used in get_task, encoded in base64
}

//Execute given command. Return response to the client
func (s *Server) Execute(buff []byte) []byte {
	var com Command
	err := yaml.Unmarshal(buff, com)
	if err != nil {
		log.Error("Error parsing command")
		return returnStatus("Error parsing request - bad or currupted struture")
	}
	if com.Command == "create_account" { //one case when we don't check password
		return s.createAccount(&com)
	}

	if !s.verifyUser(&com) {
		return returnStatus("Bad login or password")
	}

	log.Info("User %s execute %s", com.Login, com.Command)

	switch com.Command {
	case "submit":
		return s.submit(&com)
	case "get_task":
		return s.getTask(&com)
	case "contest_ranking":
		return s.getContestRanking(&com)
	case "round_ranking":
		return s.getRoundRanking(&com)
	}

	return returnStatus("Bad command name")
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

func (s *Server) checkContest(com *Command) bool { //check if contest exists and user has acces to it
	if _, ok := s.contests[com.Contest]; !ok {
		return false
	}
	if !s.users[com.Login].CheckGroup(&com.Contest) {
		return false
	}
	return true
}

func returnStatus(msg string) []byte {
	buff, _ := yaml.Marshal(ReturnMessage{Status: msg})
	return buff
}

func (s *Server) createAccount(com *Command) []byte {
	if _, ok := s.users[com.Login]; ok {
		return returnStatus("Cannot create account, there is already user with this login")
	}
	if len(com.Login) < 5 {
		return returnStatus("Login too short. Should have at least 5 letters")
	}
	if len(com.Login) > 20 {
		return returnStatus("Login too long.")
	}
	if ok, _ := regexp.Match(`^\w+$`, []byte(com.Login)); !ok {
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

func (s *Server) submit(com *Command) []byte {
	if !s.checkContest(com) {
		return returnStatus("Contest doesn't exist or you don't have permissions")
	}
	if _, ok := s.contests[com.Contest].Rounds[com.Round]; !ok {
		return returnStatus("Round doesn't exist")
	}
	if s.contests[com.Contest].Rounds[com.Round].Start.After(time.Now()) {
		return returnStatus("Round has not yet started")
	}
	if s.contests[com.Contest].Rounds[com.Round].End.Before(time.Now()) {
		return returnStatus("Round has ended")
	}
	flag := false
	for _, t := range s.contests[com.Contest].Rounds[com.Round].Tasks {
		if t == com.Task {
			flag = true
			break
		}
	}
	if !flag {
		return returnStatus("Bad task or round")
	}
	sub := submissions.Submission{
		User:    com.Login,
		Task:    com.Task,
		Round:   com.Round,
		Contest: com.Contest,
		Code:    com.Data,
	}
	buff, err := yaml.Marshal(sub)
	if err != nil {
		return returnStatus("Unknown error in parsing submission")
	}
	queue := fs.Init(s.fs.Path, "queue")
	queue.WriteFile(strconv.FormatInt(time.Now().UnixNano(), 16)+"_"+com.Login, string(buff))

	return returnStatus("ok")
}

func (s *Server) getTask(com *Command) []byte {
	if !s.checkContest(com) {
		return returnStatus("Contest doesn't exist or you don't have permissions")
	}
	if _, ok := s.contests[com.Contest].Rounds[com.Round]; !ok {
		return returnStatus("Round doesn't exist")
	}
	if s.contests[com.Contest].Rounds[com.Round].Start.After(time.Now()) {
		return returnStatus("Round has not yet started")
	}
	flag := false
	var t string
	for _, t = range s.contests[com.Contest].Rounds[com.Round].Tasks {
		if t == com.Task {
			flag = true
			break
		}
	}
	if !flag {
		return returnStatus("Bad task or round")
	}
	task := s.tasks[t]
	msg := ReturnMessage{Status: "ok", Filename: task.StatementFileName, Data: task.Statement}
	buff, _ := yaml.Marshal(msg)
	return buff
}

func (s *Server) getContestRanking(com *Command) []byte {
	if !s.checkContest(com) {
		return returnStatus("Contest doesn't exist or you don't have permissions")
	}
	msg := ReturnMessage{Status: "ok", ContestRanking: s.contests[com.Contest].Ranking}
	buff, _ := yaml.Marshal(msg)
	return buff
}

func (s *Server) getRoundRanking(com *Command) []byte {
	if !s.checkContest(com) {
		return returnStatus("Contest doesn't exist or you don't have permissions")
	}
	if _, ok := s.contests[com.Contest].Rounds[com.Round]; !ok {
		return returnStatus("Round doesn't exist")
	}
	if s.contests[com.Contest].Rounds[com.Round].Start.After(time.Now()) {
		return returnStatus("Round has not yet started")
	}
	if s.contests[com.Contest].Rounds[com.Round].ResultsShow.After(time.Now()) {
		return returnStatus("Results will be shown at" + s.contests[com.Contest].Rounds[com.Round].ResultsShow.Format("2006-01-02 15:04"))
	}
	msg := ReturnMessage{
		Status:       "ok",
		Users:        s.contests[com.Contest].Rounds[com.Round].Ranking.Names,
		Tasks:        s.contests[com.Contest].Rounds[com.Round].Tasks,
		RoundRanking: s.contests[com.Contest].Rounds[com.Round].Ranking.Points,
	}
	buff, _ := yaml.Marshal(msg)
	return buff
}