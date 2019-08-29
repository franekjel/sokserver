package server

import (
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

	ch := make(chan bool) //to control threads amount
	for i := uint16(0); i < s.conf.Workers; i++ {
		ch <- true //new thread will start only if there is something in chan
	}

	for _, group := range task.InitialTests {
		checkTestGroup(ch, &group, sub)
	}
	//now we wait for all checking threads to end (eckTestGroups is invoked without "go")

	for _, group := range task.InitialTests {
		setPoints(&group, sub)
	}
}

func checkTestGroup(ch chan bool, group *tasks.TestGroup, sub *tasks.Submission) {
	for _, test := range group.Tests {
		<-ch
		go checkTest(ch, &test, sub)
	}

}

func checkTest(ch chan bool, test *tasks.Test, sub *tasks.Submission) {
	//TODO
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
