package tasks

import (
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

//LoadTask loads task data
func LoadTask(fs *fs.Fs, conf *config.Config, name *string) *Task {
	t := &Task{
		Name:               *name,
		fs:                 fs,
		defaultTimeLimit:   conf.DefaultTimeLimit,
		defaultMemoryLimit: conf.DefaultMemoryLimit,
	}
	//TODO

	return t
}
