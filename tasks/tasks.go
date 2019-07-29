package tasks

import (
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

//Task - struct for holding task data and performing tests
type Task struct {
	Name               string
	fs                 *fs.Fs
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
	var re = make([]string, len(files))
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
	var re = make([]string, len(files))
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
	if s[1] == "ocen" || s[1] == "0" {
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
	var out map[string]bool
	for i := 0; i < len(outList); i++ {
		out[outList[i]] = true
	}
	for _, i := range inList {
		_, ok := out[i]
		if !ok {
			Log(WARN, "Missing output for test %s", i)
		} else {
			reg := regexp.MustCompile(`^(\d+)([a-z]*)$`)
			s := reg.FindStringSubmatch(i)
			if len(s) == 0 {
				Log(WARN, "Bad test name %s, skipping", i)
				continue
			}
			if isInitialTest(s) {
				t.insertInitalTest(&test{name: i}, &s[1])
			} else {
				t.insertFinalTest(&test{name: i}, &s[1])
			}
		}
	}

}

//LoadTask loads task data
func LoadTask(fs *fs.Fs, conf *config.Config, name *string) *Task {
	t := &Task{
		Name:               *name,
		fs:                 fs,
		defaultTimeLimit:   conf.DefaultTimeLimit,
		defaultMemoryLimit: conf.DefaultMemoryLimit,
	}
	Log(INFO, "Loading task %s", t.fs.Path)

	return t
}
