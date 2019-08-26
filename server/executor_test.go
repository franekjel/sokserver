package server

import (
	"os"
	"testing"
	"time"

	"github.com/franekjel/sokserver/contests"
	"github.com/franekjel/sokserver/fs"
	"github.com/franekjel/sokserver/rounds"
	"github.com/franekjel/sokserver/tasks"
	"github.com/franekjel/sokserver/users"
)

func TestCreateAccountOK(t *testing.T) {
	//simple call function in /tmp, should work
	var s Server
	s.fs = fs.Init("/tmp", "")
	s.fs.CreateDirectory("users")
	s.users = make(map[string]*users.User)
	com := Command{
		Login:    "testLogin",
		Password: "password",
		Command:  "create_account",
	}
	buff := s.createAccount(&com)
	if string(buff) != "status: ok\n" {
		t.Error(string(buff))
	}
	s.users["testLogin"].AddToGroup("aaa")
	s.users["testLogin"].AddToGroup("bbb")
	s.users["testLogin"].SaveData()
}

func TestCreateAccountDuplicated(t *testing.T) {
	//create account two times
	var s Server
	s.fs = fs.Init("/tmp", "")
	s.fs.CreateDirectory("users")
	s.users = make(map[string]*users.User)
	com := Command{
		Login:    "testLogin",
		Password: "password",
		Command:  "create_account",
	}
	buff := s.createAccount(&com)
	buff = s.createAccount(&com)
	if string(buff) != "status: Cannot create account, there is already user with this login\n" {
		t.Error(string(buff))
	}
}

func TestSubmit(t *testing.T) {
	s := prepareTestEnvironment()
	defer s.cleanTestEnvironment()

	com := Command{
		Login:    "testLogin",
		Password: "password",
		Command:  "submit",
		Contest:  "con1",
		Round:    "round1",
		Task:     "taska",
		Data: `include<stdio.h>
		int main(){
			printf("answer");
		}`,
	}
	buff := s.submit(&com)
	if string(buff) != "status: ok\n" {
		t.Error(string(buff))
	}
}

func TestContestRanking(t *testing.T) {
	s := prepareTestEnvironment()
	defer s.cleanTestEnvironment()

	com := Command{
		Login:    "testLogin",
		Password: "password",
		Command:  "contest_ranking",
		Contest:  "con1",
	}
	buff := s.getContestRanking(&com)
	if string(buff[0:10]) != "status: ok" {
		t.Error(string(buff))
	}
}

func prepareTestEnvironment() *Server {
	var s Server
	s.fs = fs.Init("/tmp", "")
	s.fs.CreateDirectory("queue")
	s.users = make(map[string]*users.User)
	s.users["testLogin"] = &users.User{ //adding user
		Groups: []string{"con1"},
	}
	s.contests = make(map[string]*contests.Contest)
	s.contests["con1"] = &contests.Contest{ //adding contest
		Rounds:  make(map[string]*rounds.Round),
		Ranking: map[string]uint{"testLogin": 120, "foo": 200, "bar": 45},
	}
	s.contests["con1"].Rounds["round1"] = &rounds.Round{Name: "round1", Start: time.Now(), End: time.Now().Add(time.Minute)} //adding round
	s.tasks = make(map[string]*tasks.Task)
	s.tasks["taska"] = &tasks.Task{Name: "taska", Statement: "Task Statement", StatementFileName: "task.txt"} //Statement normally is store in base64
	s.tasks["taskb"] = &tasks.Task{Name: "taskb"}
	s.tasks["taskc"] = &tasks.Task{Name: "taskc"}

	r := s.contests["con1"].Rounds["round1"]
	r.Tasks = []string{"taska", "taskb", "taskc"}
	s.contests["con1"].Rounds["round1"].Ranking = rounds.RoundRanking{
		Points: make([][]uint, 0, 3),
		Names:  make([]string, 3),
	}

	res := map[string]uint{ //map to get results, key is concatenation of login and task name
		"testLogintaska": 100,
		"testLogintaskc": 50,
		"footaskb":       60,
		"bartaskc":       20,
	}

	for i, user := range []string{"testLogin", "foo", "bar"} {
		r.Ranking.Points = append(r.Ranking.Points, make([]uint, 3))
		r.Ranking.Names[i] = user
		for j, task := range r.Tasks {
			r.Ranking.Points[i][j] = res[user+task]
		}
	}
	return &s
}

func (s *Server) cleanTestEnvironment() {
	p := s.fs.Path
	os.RemoveAll(p + "/queue")
	os.RemoveAll(p + "/users")
}

func TestRoundRanking(t *testing.T) {
	s := prepareTestEnvironment()
	defer s.cleanTestEnvironment()
	com := Command{
		Login:    "testLogin",
		Password: "password",
		Command:  "round_ranking",
		Contest:  "con1",
		Round:    "round1",
	}
	buff := s.getRoundRanking(&com)
	if string(buff[0:10]) != "status: ok" {
		t.Error(string(buff))
	}
}

func TestGetTask(t *testing.T) {
	s := prepareTestEnvironment()
	defer s.cleanTestEnvironment()
	com := Command{
		Login:    "testLogin",
		Password: "password",
		Command:  "get_task",
		Contest:  "con1",
		Round:    "round1",
		Task:     "taska",
	}
	buff := s.getTask(&com)
	if string(buff[0:10]) != "status: ok" {
		t.Error(string(buff))
	}
}
