package tasks

import (
	"gopkg.in/yaml.v2"
	"regexp"
	"strings"

	"github.com/franekjel/sokserver/config"
	"github.com/franekjel/sokserver/fs"
	. "github.com/franekjel/sokserver/logger"
)

type test struct {
	name        string
	timeLimit   uint
	memoryLimit uint
}

type testGroup struct {
	name   string
	tests  map[string]test
	points uint
}

type taskConfig struct {
	Title        string          `yaml:"title"`
	TimeLimit    uint            `yaml:"time_limit"`
	MemoryLimit  uint            `yaml:"memory_limit"`
	TimeLimits   map[string]uint `yaml:"time_limits,flow"`
	MemoryLimits map[string]uint `yaml:"memory_limits,flow"`
}

//Task - struct for holding task data and performing tests
type Task struct {
	Name               string
	fs                 *fs.Fs
	config             taskConfig
	defaultMemoryLimit uint
	defaultTimeLimit   uint
	initialTests       map[string]testGroup
	finalTests         map[string]testGroup
}

func (t *Task) listInputs() []string {
	if !t.fs.FileExist("in") {
		return nil
	}
	files := t.fs.ListFiles("in")
	var re = make([]string, 0, len(files))
	for _, file := range files {
		if strings.HasSuffix(file, ".in") {
			re = append(re, strings.TrimSuffix(file, ".in"))
		}
	}
	if len(re) == 0 {
		return nil
	}
	return re
}

func (t *Task) listOutputs() []string {
	if !t.fs.FileExist("out") {
		return nil
	}
	files := t.fs.ListFiles("out")
	var re = make([]string, 0, len(files))
	for _, file := range files {
		if strings.HasSuffix(file, ".out") {
			re = append(re, strings.TrimSuffix(file, ".out"))
		}
	}
	if len(re) == 0 {
		return nil
	}
	return re
}

func isInitialTest(s []string) bool {
	if s[1] == "ocen" || s[2] == "0" {
		return true
	}
	return false
}

func (t *Task) insertInitalTest(te *test, groupName *string) {
	group, ok := t.initialTests[*groupName]
	if ok {
		group.tests[te.name] = *te
	} else {
		t.initialTests[*groupName] = testGroup{
			name:  *groupName,
			tests: map[string]test{te.name: *te},
		}
	}
}

func (t *Task) insertFinalTest(te *test, groupName *string) {
	group, ok := t.finalTests[*groupName]
	if ok {
		group.tests[te.name] = *te
	} else {
		t.finalTests[*groupName] = testGroup{
			name:  *groupName,
			tests: map[string]test{te.name: *te},
		}
	}
}

func (t *Task) addTests() {
	inList := t.listInputs()
	outList := t.listOutputs()
	var out = map[string]bool{}
	for i := 0; i < len(outList); i++ {
		out[outList[i]] = true
	}
	for _, i := range inList {
		if _, ok := out[i]; !ok {
			Log(WARN, "Missing output for test %s", i)
		} else {
			reg := regexp.MustCompile(`^([a-z]*)(\d+)([a-z]*)$`)
			s := reg.FindStringSubmatch(i)
			if len(s) == 0 {
				Log(WARN, "Bad test name %s, skipping", i)
				continue
			}
			if isInitialTest(s) {
				t.insertInitalTest(&test{name: i}, &s[2])
			} else {
				t.insertFinalTest(&test{name: i}, &s[2])
			}
		}
	}
}

func (t *Task) parseConfig(globalConfig *config.Config) {
	t.config.MemoryLimit = globalConfig.DefaultMemoryLimit
	t.config.TimeLimit = globalConfig.DefaultTimeLimit
	t.config.Title = t.Name
	if !t.fs.FileExist("config.yml") {
		return
	}
	s := t.fs.ReadFile("config.yml")
	yaml.Unmarshal(s, &t.config)
}

func (t *Task) setLimits(tests *map[string]testGroup) {
	for _, i := range *tests {
		Log(DEBUG, "Test group: %s:", i.name)
		for _, j := range i.tests {

			j.memoryLimit = t.config.MemoryLimit
			j.timeLimit = t.config.TimeLimit
			if limit, ok := t.config.TimeLimits[j.name]; ok {
				if limit < 100000 {
					j.timeLimit = limit
				}
			}
			if limit, ok := t.config.MemoryLimits[j.name]; ok {
				if limit < 1024*1024 {
					j.memoryLimit = limit
				}
			}
			Log(DEBUG, "Test: %s, time: %d, memory %d", j.name, j.timeLimit, j.memoryLimit)
		}
	}

}

func (t *Task) setTestsLimits() {
	t.setLimits(&t.initialTests)
	t.setLimits(&t.finalTests)
}

//LoadTask loads task data
func LoadTask(fs *fs.Fs, conf *config.Config, name *string) *Task {
	t := &Task{
		Name:               *name,
		fs:                 fs,
		defaultTimeLimit:   conf.DefaultTimeLimit,
		defaultMemoryLimit: conf.DefaultMemoryLimit,
		initialTests:       make(map[string]testGroup, 1),
		finalTests:         make(map[string]testGroup, 10),
	}
	t.parseConfig(conf)
	Log(INFO, "Loading task %s: %s", t.fs.Path, t.config.Title)
	t.addTests()
	t.setTestsLimits()
	return t
}
