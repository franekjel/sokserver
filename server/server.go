package server

import(
	. "github.com/franekjel/sokserver/logger"
	"github.com/franekjel/sokserver/users"
	"github.com/franekjel/sokserver/config"
	"github.com/franekjel/sokserver/fs"
)

//Server stores main SOK data
type Server struct {
	users map[string] users.User
	conf *config.Config
	fs *fs.Fs
}

func (s *Server) loadConfig(){
	var buff = new([]byte)
	if !s.fs.FileExist("sok.yml") {
		s.conf = config.MakeConfig(buff)
		Log(WARN, "sok.yml doesn't exists or is inaccesible")
		Log(DEBUG, string(s.conf.GetConfig()))
		defer s.fs.WriteFile("sok.yml", string(s.conf.GetConfig()))
	} else {
	buff = s.fs.ReadFile("sok.yml")
	s.conf = config.MakeConfig(buff)
	}
	Log(INFO, "port: %d", s.conf.Port)
	Log(INFO, "workers: %d", s.conf.Workers)
	Log(INFO, "default memory limit: %d", s.conf.DefaultMemoryLimit)
	Log(INFO, "default time limit: %d", s.conf.DefaultTimeLimit)
}

//InitServer initializes structures and starts listening
func InitServer(dir string) {
	var server Server
	server.fs=fs.Init(dir, "")
	server.loadConfig()

}

