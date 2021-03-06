package contests

import (
	"gopkg.in/yaml.v2"
	"time"

	"github.com/franekjel/sokserver/fs"
	log "github.com/franekjel/sokserver/logger"
	"github.com/franekjel/sokserver/tasks"
)

//Round has tasks, start time, end time and time when results will be show
type Round struct {
	Name        string
	Tasks       []string
	Start       time.Time
	End         time.Time
	ResultsShow time.Time
	fs          *fs.Fs
	Ranking     RoundRanking
}

//RoundRanking hold ranking as a two-dimensional array and additional slice contains columns description (tasks names).
//First columns is sum of row (labeled as "Sum")
type RoundRanking struct {
	Points [][]uint
	Names  []string
}

//struct to parsing round.yml neccessary due to yaml date parsing issues
type roundParse struct {
	Name        string   `yaml:"name"`
	Tasks       []string `yaml:"tasks,flow"`
	Start       string   `yaml:"start_date"`
	End         string   `yaml:"end_date"`
	ResultsShow string   `yaml:"results_show_date"`
}

//ListSubmissions is used in command list_submissions to (surprise!) list submissions
func (r *Round) ListSubmissions(login string, task string) []*tasks.Submission {
	if !r.fs.FileExist(fs.Join(login, task)) {
		return []*tasks.Submission{}
	}
	dir := fs.Init(r.fs.Path, fs.Join(login, task))
	subs := dir.ListFiles("")
	re := make([]*tasks.Submission, 0, len(subs))
	for _, file := range subs {
		sub := tasks.LoadSubmission(dir.ReadFile(file))
		if sub != nil {
			re = append(re, sub)
		}
	}
	return re
}

//GetSubmission return submission with given id
func (r *Round) GetSubmission(login string, task string, id string) *tasks.Submission {
	if !r.fs.FileExist(fs.Join(login, task)) {
		return nil
	}
	dir := fs.Init(r.fs.Path, fs.Join(login, task))
	if !dir.FileExist(id + ".yml") {
		return nil
	}
	sub := tasks.LoadSubmission(dir.ReadFile(id + ".yml"))
	return sub
}

func (r *Round) verifyTasks(tasks map[string]*tasks.Task) {
	newTasks := make([]string, 0, len(r.Tasks))
	for _, task := range r.Tasks {
		if _, ok := tasks[task]; ok {
			newTasks = append(newTasks, task)
		} else {
			log.Error("%s: missing task %s", r.fs.Path, task)
		}
	}
	r.Tasks = newTasks
}

func (r *Round) loadData(tasks map[string]*tasks.Task) {
	if !r.fs.FileExist("round.yml") {
		log.Fatal("Round settings missing! %s", r.fs.Path)
	}
	buff := r.fs.ReadFile("round.yml")
	var temp = new(roundParse)
	err := yaml.Unmarshal(buff, temp)
	if err != nil {
		log.Error("%s", err.Error())
	}
	time.Now().Format("2006-01-02 15:04")
	r.Name = temp.Name
	r.Tasks = temp.Tasks //TODO: tasks existence check?
	r.Start, err = time.Parse("2006-01-02 15:04", temp.Start)
	if err != nil {
		log.Fatal("Wrong start date in round config in %s: %s", r.fs.Path, err.Error())
	}
	r.End, err = time.Parse("2006-01-02 15:04", temp.End)
	if err != nil {
		log.Fatal("Wrong end date in round config in %s: %s", r.fs.Path, err.Error())
	}
	r.ResultsShow, err = time.Parse("2006-01-02 15:04", temp.ResultsShow)
	if err != nil {
		log.Warn("Wrong or missing result show date in %s, using start date instead (results will be show immediately)", r.fs.Path)
		r.ResultsShow = r.Start
	}
	r.verifyTasks(tasks)
}

//roundname/username/taskname.yml holds last user submissions. This function get results from this submission
func (r *Round) getResult(user, task string) uint {
	submission := r.fs.ReadFile(fs.Join(user, task+".yml"))
	return tasks.LoadSubmission(submission).Sum
}

func (r *Round) loadRanking() {
	dirs := r.fs.ListDirs(".")
	n := len(r.Tasks)
	r.Ranking = RoundRanking{
		Points: make([][]uint, 0, len(dirs)),
		Names:  make([]string, n),
	}
	for i, user := range dirs {
		r.Ranking.Points = append(r.Ranking.Points, make([]uint, n))
		for j, task := range r.Tasks {
			r.Ranking.Points[i][j] = r.getResult(user, task)
		}
	}
}

//LoadRound loads round in given folder
func LoadRound(fs *fs.Fs, tasks map[string]*tasks.Task) *Round {
	log.Info("Loading round %s", fs.Path)
	round := new(Round)
	round.fs = fs
	round.loadData(tasks)
	round.loadRanking()
	log.Debug("%s %s %+q", round.Name, round.Start, round.Tasks)
	return round
}
