package server

import (
	"os/exec"
	"time"

	log "github.com/franekjel/sokserver/logger"

	"github.com/franekjel/sokserver/fs"
	"github.com/franekjel/sokserver/tasks"
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
			time.Sleep(time.Millisecond * 500) //maybe not best solution...
		}
	}
}

//check given submission submission - in many threads (specified in config - workers)
func (s *Server) check(buff []byte) {
	sub := tasks.LoadSubmission(buff)
	log.Info("Checking submission from %s, contest %s, round %s, task %s", sub.User, sub.Contest, sub.Round, sub.Task)
	task := s.tasks[sub.Task]

	s.compileCode(sub) //TODO check for errors

	ch := make(chan bool) //to control threads amount
	for i := uint16(0); i < s.conf.Workers; i++ {
		ch <- true //new thread will start only if there is something in chan
	}

	for _, group := range task.InitialTests {
		s.checkTestGroup(ch, &group, sub)
	}
	//now we wait for all checking threads to end (eckTestGroups is invoked without "go")

	for _, group := range task.InitialTests {
		setPoints(&group, sub)
	}
}

func (s *Server) checkTestGroup(ch chan bool, group *tasks.TestGroup, sub *tasks.Submission) {
	for _, test := range group.Tests {
		<-ch
		go s.checkTest(ch, &test, sub)
	}

}

//compiles solution code. TODO - gcc compile flags in config
func (s *Server) compileCode(sub *tasks.Submission) {
	cmd := exec.Command("g++", "--static")
	cmd.Output()
	//TODO check for compilation errors, if there are any stop further proceeding
}

//execute solution code on test. TODO - configurable sio2jail path
func (s *Server) checkTest(ch chan bool, test *tasks.Test, sub *tasks.Submission) {
	commandString := ""
	cmd := exec.Command(fs.Join(s.fs.Path, "sio2jail"), commandString)
	cmd.Output()
	//TODO - check output for status and time, set in submissions struct, put value to chan
}

//this function calcs and sets points for TestGroup
func setPoints(group *tasks.TestGroup, sub *tasks.Submission) {
	pointsFactor := 1.0
	for _, test := range group.Tests {
		if sub.Results[test.Name] == "MLE" || sub.Results[test.Name] == "RV" || sub.Results[test.Name] == "RE" || sub.Results[test.Name] == "TLE" {
			pointsFactor = 0
			break
		} else if float64(test.TimeLimit)/float64(sub.Time[test.Name]) > 0.5 {
			temp := 2 * (1.0 - float64(test.TimeLimit)/float64(sub.Time[test.Name]))
			if temp < pointsFactor {
				pointsFactor = temp
			}
		}
	}
	sub.Points[group.Name] = uint(pointsFactor * float64(group.Points))
}
