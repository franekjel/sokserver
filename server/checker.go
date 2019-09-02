package server

import (
	"bytes"
	"io"
	"os/exec"
	"strconv"
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
			fs.RemoveFile(files[0])
		} else {
			time.Sleep(time.Millisecond * 500) //maybe not best solution...
		}
	}
}

//check given submission submission - in many threads (specified in config - workers)
func (s *Server) check(buff []byte) {
	sub := tasks.LoadSubmission(buff)
	sub.Sum = 0
	log.Info("Checking submission from %s, contest %s, round %s, task %s", sub.User, sub.Contest, sub.Round, sub.Task)
	task := s.tasks[sub.Task]

	ok, err := compileCode(sub) //TODO check for errors
	if !ok {
		n := 200
		if len(err) < 200 {
			n = len(err)
		}
		sub.InitialStatus = err[:n] //since c++ compilation error may be very long, we store only beginning
		return
	}

	ch := make(chan bool) //to control threads amount
	for i := uint16(0); i < s.conf.Workers; i++ {
		ch <- true //new thread will start only if there is something in chan
	}

	for _, group := range task.InitialTests {
		s.checkTestGroup(ch, &group, sub)
	}
	//now we wait for all checking threads to end (checkTestGroups is invoked without "go")

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
func compileCode(sub *tasks.Submission) (bool, string) {
	cmd := exec.Command("g++", "-x", "c++", "--static", "-O2", "-o", "/tmp/exe", "-")
	stdin, _ := cmd.StdinPipe()
	io.WriteString(stdin, sub.Code)
	stdin.Close()
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		if err.Error() == "exit status 1" {
			return false, stderr.String()
		}
		log.Fatal("Cannot compile program: %s", err.Error())
	}
	return true, ""
}

//execute solution code on test. TODO - configurable sio2jail path
func (s *Server) checkTest(ch chan bool, test *tasks.Test, sub *tasks.Submission) {
	inDir := fs.Init(fs.CreatePath(s.fs.Path, "tasks", sub.Task, "in"), "")
	cmd := exec.Command(fs.Join(s.fs.Path, "sio2jail"),
		"-b", "/tmp/box:/:ro",
		"--memory-limit", strconv.FormatUint(uint64(test.MemoryLimit), 10),
		"--instruction-count-limit", strconv.FormatUint(uint64(test.TimeLimit*2*1000000), 10), //TODO - clean this magic number, this is time*2*10^9/1000
		"--rtimelimit", strconv.FormatUint(uint64(test.TimeLimit*5), 10)+"ms",
		"--output-limit", "64K",
		"-o", "oiaug",
		"/tmp/exe",
		"<"+fs.Join(inDir.Path, sub.Task+test.Name+".in"),
	)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	output, err := cmd.Output()
	if err != nil {
		log.Warn("Runtime error during cheking submission %s", sub.Id)
		sub.Results[test.Name] = "runtime error"
		return
	}

	results := bytes.Fields(stderr.Bytes())
	if string(results[0]) != "OK" {
		sub.Results[test.Name] = string(results[len(results)-1])
		return
	}
	outDir := fs.Init(fs.CreatePath(s.fs.Path, "tasks", sub.Task, "out"), "")
	if ok, diff := checkOutput(output, outDir.ReadFile(sub.Task+test.Name+".out")); !ok {
		sub.Results[test.Name] = diff
		return
	}

	sub.Results[test.Name] = "OK"
	temp, _ := strconv.ParseUint(string(results[2]), 10, 64)
	sub.Time[test.Name] = uint(temp)
}

func checkOutput(userOutput []byte, goodOutput []byte) (bool, string) {
	return true, ""
}

//this function calcs and sets points for TestGroup
func setPoints(group *tasks.TestGroup, sub *tasks.Submission) {
	pointsFactor := 1.0
	for _, test := range group.Tests {
		if sub.Results[test.Name] != "OK" {
			pointsFactor = 0
			break
		} else if float64(sub.Time[test.Name])/float64(test.TimeLimit) > 0.5 { //if program exceeds half of time points will decrase linearly (to 0 at time limit)
			cur := 2 * (1.0 - float64(sub.Time[test.Name])/float64(test.TimeLimit))
			if cur < pointsFactor {
				pointsFactor = cur
			}
		}
	}
	sub.Points[group.Name] = uint(pointsFactor * float64(group.Points))
}
