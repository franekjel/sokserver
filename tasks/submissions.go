package tasks

import (
	log "github.com/franekjel/sokserver/logger"
	"gopkg.in/yaml.v2"
)

//Submission holds users submissions data
type Submission struct {
	Id            string            `yaml:"id"`
	User          string            `yaml:"user"`
	Task          string            `yaml:"task"`
	Round         string            `yaml:"round"`
	Contest       string            `yaml:"contest"`
	Code          string            `yaml:"code"`
	Results       map[string]string `yaml:"results,omitempty"`        //status for each test like OK, Bad result, timeout etc
	Time          map[string]uint   `yaml:"time,omitempty"`           //time of execution of each test
	MaxTime       map[string]uint   `yaml:"max_time,omitempty"`       //max time of execution - empty, used only during marshall to send to user TODO
	Points        map[string]uint   `yaml:"points,omitempty"`         //points for each testgroup
	Sum           uint              `yaml:"sum,omitempty"`            //sum of points
	InitialStatus string            `yaml:"initial_status,omitempty"` //status of initial tests
	FinalStatus   string            `yaml:"final_status,omitempty"`   //status of final tests
}

//LoadSubmission load submission data from yaml string
func LoadSubmission(buff []byte) *Submission {
	var s Submission
	err := yaml.Unmarshal(buff, &s)
	if err != nil {
		log.Error("Error parsing user submission: %s", buff)
		return nil
	}
	return &s
}
