package rounds

import (
	"gopkg.in/yaml.v2"
	"time"

	"github.com/franekjel/sokserver/fs"
	. "github.com/franekjel/sokserver/logger"
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
}

//struct to parsing round.yml neccessary due to yaml date parsing issues
type roundParse struct {
	Name        string   `yaml:"name"`
	Tasks       []string `yaml:"tasks,flow"`
	Start       string   `yaml:"start_date"`
	End         string   `yaml:"end_date"`
	ResultsShow string   `yaml:"results_show_date"`
}

func (r *Round) verifyTasks(tasks map[string]*tasks.Task) {
	newTasks := make([]string, 0, len(r.Tasks))
	for _, task := range r.Tasks {
		if _, ok := tasks[task]; ok {
			newTasks = append(newTasks, task)
		} else {
			Log(ERR, "%s: missing task %s", r.fs.Path, task)
		}
	}
	r.Tasks = newTasks
}

func (r *Round) loadData(tasks map[string]*tasks.Task) {
	if !r.fs.FileExist("round.yml") {
		Log(FATAL, "Round settings missing! %s", r.fs.Path)
	}
	buff := r.fs.ReadFile("round.yml")
	var temp = new(roundParse)
	err := yaml.Unmarshal(buff, temp)
	if err != nil {
		Log(ERR, "%s", err.Error())
	}
	time.Now().Format("2006-01-02 15:04")
	r.Name = temp.Name
	r.Tasks = temp.Tasks //TODO: tasks existence check?
	r.Start, err = time.Parse("2006-01-02 15:04", temp.Start)
	if err != nil {
		Log(FATAL, "Wrong start date in round config in %s: %s", r.fs.Path, err.Error())
	}
	r.End, err = time.Parse("2006-01-02 15:04", temp.End)
	if err != nil {
		Log(FATAL, "Wrong end date in round config in %s: %s", r.fs.Path, err.Error())
	}
	r.ResultsShow, err = time.Parse("2006-01-02 15:04", temp.ResultsShow)
	if err != nil {
		Log(WARN, "Wrong or missing result show date in %s, using star date instead (results will be show immediately)", r.fs.Path)
		r.ResultsShow = r.Start
	}
	r.verifyTasks(tasks)
}

//LoadRound loads round in given folder
func LoadRound(fs *fs.Fs, tasks map[string]*tasks.Task) *Round {
	Log(INFO, "Loading round %s", fs.Path)
	round := new(Round)
	round.fs = fs
	round.loadData(tasks)
	Log(DEBUG, "%s %s %+q", round.Name, round.Start, round.Tasks)
	return round
}
