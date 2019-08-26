package server

import (
	"strings"

	log "github.com/franekjel/sokserver/logger"

	"github.com/franekjel/sokserver/config"
	"github.com/franekjel/sokserver/contests"
	"github.com/franekjel/sokserver/fs"
	"github.com/franekjel/sokserver/tasks"
	"github.com/franekjel/sokserver/users"
)

//Server stores main SOK data
type Server struct {
	users    map[string]*users.User
	tasks    map[string]*tasks.Task
	contests map[string]*contests.Contest
	conf     *config.Config
	fs       *fs.Fs
}

func (s *Server) loadConfig() {
	var buff = new([]byte)
	if !s.fs.FileExist("sok.yml") {
		s.conf = config.MakeConfig(buff)
		log.Warn("sok.yml doesn't exists or is inaccesible")
		log.Debug(string(s.conf.GetConfig()))
		defer s.fs.WriteFile("sok.yml", string(s.conf.GetConfig()))
	} else {
		*buff = s.fs.ReadFile("sok.yml")
		s.conf = config.MakeConfig(buff)
	}
	log.Info("port: %d", s.conf.Port)
	log.Info("workers: %d", s.conf.Workers)
	log.Info("default memory limit: %d", s.conf.DefaultMemoryLimit)
	log.Info("default time limit: %d", s.conf.DefaultTimeLimit)
}

func (s *Server) loadUsers() {
	log.Info("---Loading users data")
	s.users = make(map[string]*users.User)
	if !s.fs.FileExist("users") {
		log.Warn("\"users\" directory doesn't exist, creating")
		s.fs.CreateDirectory("users")
		return //we can't parse any users
	}
	//at this point we are sure that "users" exists
	dir := fs.Init(s.fs.Path, "users")
	for _, login := range dir.ListFiles("") {
		if strings.HasSuffix(login, ".yml") {
			login = strings.TrimSuffix(login, ".yml")
			log.Info("Loading user %s", login)
			s.users[login] = users.LoadUser(dir, login)
		}
	}
}

func (s *Server) loadTasks() {
	log.Info("---Loading tasks")
	s.tasks = make(map[string]*tasks.Task)
	if !s.fs.FileExist("tasks") {
		log.Warn("\"tasks\" directory doesn't exist, creating")
		s.fs.CreateDirectory("tasks")
		return
	}
	dir := fs.Init(s.fs.Path, "tasks")
	for _, name := range dir.ListDirs("") {
		s.tasks[name] = tasks.LoadTask(fs.Init(dir.Path, name), s.conf, &name)
	}
}

func (s *Server) loadContests() {
	log.Info("---Loading contests")
	s.contests = make(map[string]*contests.Contest)
	if !s.fs.FileExist("contests") {
		log.Warn("\"contests\" directory doesn't exist, creating")
		s.fs.CreateDirectory("contests")
		return
	}
	dir := fs.Init(s.fs.Path, "contests")
	for _, name := range dir.ListDirs("") {
		s.contests[name] = contests.LoadContest(fs.Init(dir.Path, name), s.tasks)
	}
}

//InitServer initializes structures and starts listening
func InitServer(dir string) {
	var server Server
	server.fs = fs.Init(dir, "")
	server.loadConfig()
	server.loadUsers()
	server.loadTasks()
	server.loadContests()
	if !server.fs.FileExist("queue") {
		log.Warn("\"contests\" directory doesn't exist, creating")
		server.fs.CreateDirectory("queue")
	}
	ch := make(chan *connectionData)
	go server.startListening(ch)
	go server.CheckSubmissions(fs.Init(server.fs.Path, "queue"), ch)

	for { //main loop - execute users commands
		data := <-ch
		response := server.Execute(data.data)
		sendResponse(data.conn, response)
	}
}
