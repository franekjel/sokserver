package rounds

import (
	"time"

	"github.com/franekjel/sokserver/fs"
	"github.com/franekjel/sokserver/tasks"
)

//Round has tasks, start time, end time and time where results will be show
type Round struct {
	fs          *fs.Fs
	tasks       map[string]*tasks.Task
	start       time.Time
	end         time.Time
	resultsShow time.Time
}
