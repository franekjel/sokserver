package server

import (
	"time"

	"github.com/franekjel/sokserver/fs"
	//"github.com/franekjel/sokserver/submissions"
)

//CheckSubmissions read and check submissions from sok/queue
func (s *Server) CheckSubmissions(fs *fs.Fs, ch chan *connectionData) {
	for {
		if len(fs.ListFiles("")) != 0 {
			//there are submissions, get first and check
			files := fs.ListFiles("")
			buff := fs.ReadFile(files[0])
			s.check(buff)
		} else {
			time.Sleep(time.Second) //maybe not best solution...
		}
	}
}

func (s *Server) check(buff []byte) {
	//check given submission submission - in many threads (specified in config - workers)
	//sub := submissions.LoadSubmission(buff)
	//task := s.tasks[sub.Task]

}
