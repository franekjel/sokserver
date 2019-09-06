package server

import (
	"regexp"
	"strconv"
	"time"

	"github.com/franekjel/sokserver/contests"
	"github.com/franekjel/sokserver/fs"
	"gopkg.in/yaml.v2"

	log "github.com/franekjel/sokserver/logger"
	"github.com/franekjel/sokserver/tasks"
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
	Status         string           `yaml:"status"`                       //status may be "ok" or contains error
	ContestRanking map[string]uint  `yaml:"contest_ranking,omitempty"`    //used in contest_ranking
	Tasks          []string         `yaml:"tasks,omitempty"`              //used in round_ranking
	Users          []string         `yaml:"users,omitempty"`              //used in round_ranking
	RoundRanking   [][]uint         `yaml:"round_ranking,omitempty,flow"` //used in round_ranking
	Filename       string           `yaml:"filename,omitempty"`           //used in get_task
	Data           string           `yaml:"data,omitempty"`               //used in get_task, encoded in base64
	Submissions    [][3]string      `yaml:"submissions,omitempty,flow"`   //used in list_submissions
	Submission     tasks.Submission `yaml:"submission,omitempty"`         //used in get_submission
	Contests       [][2]string      `yaml:"contests,omitempty,flow"`      //used in list_contests
}

//Execute given command. Return response to the client
func (s *Server) Execute(buff []byte) []byte {
	com := Command{}
	err := yaml.Unmarshal(buff, &com)
	if err != nil {
		log.Error("Error parsing command %s", err.Error())
		return returnStatus("Error parsing request - bad or currupted struture:")
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
	case "list_submissions":
		return s.listSubmissions(&com)
	case "get_submission":
		return s.getSubmission(&com)
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
	if ok, status := s.checkRound(com); !ok {
		return returnStatus(status)
	}
	if s.contests[com.Contest].Rounds[com.Round].End.Before(time.Now()) {
		return returnStatus("Round ended")
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
	sub := tasks.Submission{
		Id:      strconv.FormatInt(time.Now().UnixNano(), 16),
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
	queue.WriteFile(sub.Id+"_"+com.Login, string(buff))

	return returnStatus("ok")
}

//check if round has given task
func hasTask(round *contests.Round, task string) bool {
	flag := false
	for _, t := range round.Tasks {
		if t == task {
			flag = true
			break
		}
	}
	return flag
}

func (s *Server) getTask(com *Command) []byte {
	if ok, status := s.checkRound(com); !ok {
		return returnStatus(status)
	}

	if !hasTask(s.contests[com.Contest].Rounds[com.Round], com.Task) {
		return returnStatus("Bad task or round")
	}
	task := s.tasks[com.Task]
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
	if ok, status := s.checkRound(com); !ok {
		return returnStatus(status)
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

func (s *Server) listSubmissions(com *Command) []byte {
	if ok, status := s.checkRound(com); !ok {
		return returnStatus(status)
	}
	if !hasTask(s.contests[com.Contest].Rounds[com.Round], com.Task) {
		return returnStatus("Bad task or round")
	}
	subs := s.contests[com.Contest].Rounds[com.Round].ListSubmissions(com.Login, com.Task)
	subsList := make([][3]string, 0, len(subs))
	for _, sub := range subs {
		var temp [3]string
		temp[0] = sub.Id
		if s.contests[com.Contest].Rounds[com.Round].ResultsShow.After(time.Now()) { //if results are hidden
			temp[1] = sub.InitialStatus
			temp[2] = "0"
		} else {
			temp[1] = sub.FinalStatus
			temp[2] = strconv.FormatUint(uint64(sub.Sum), 10)
		}
		subsList = append(subsList, temp)
	}

	msg := ReturnMessage{
		Status:      "ok",
		Submissions: subsList,
	}
	buff, _ := yaml.Marshal(msg)
	return buff
}

func (s *Server) getSubmission(com *Command) []byte {
	if ok, status := s.checkRound(com); !ok {
		return returnStatus(status)
	}
	if !hasTask(s.contests[com.Contest].Rounds[com.Round], com.Task) {
		return returnStatus("Bad task or round")
	}
	sub := s.contests[com.Contest].Rounds[com.Round].GetSubmission(com.Login, com.Task, com.Data)
	if sub == nil {
		return returnStatus("Bad submission ID")
	}
	msg := ReturnMessage{
		Status:     "ok",
		Submission: *sub,
	}
	buff, _ := yaml.Marshal(msg)
	return buff
}

func (s *Server) checkRound(com *Command) (bool, string) {
	if !s.checkContest(com) {
		return false, "Contest doesn't exist or you don't have permissions"
	}
	if _, ok := s.contests[com.Contest].Rounds[com.Round]; !ok {
		return false, "Round doesn't exist"
	}
	if s.contests[com.Contest].Rounds[com.Round].Start.After(time.Now()) {
		return false, "Round has not yet started"
	}
	return true, ""
}

func (s *Server) listContests(com *Command) []byte {
	conList := make([][2]string, 0, len(s.contests))
	for id, con := range s.contests {
		var temp [2]string
		temp[0] = id
		temp[1] = con.Name
		conList = append(conList, temp)
	}
	msg := ReturnMessage{
		Status:   "ok",
		Contests: conList,
	}
	buff, _ := yaml.Marshal(msg)
	return buff
}
