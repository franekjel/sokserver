package submissions

import (
	. "github.com/franekjel/sokserver/logger"
	"gopkg.in/yaml.v2"
)

//Submission holds users submissions data
type Submission struct {
	User    string            `yaml:"user"`
	Task    string            `yaml:"task"`
	Round   string            `yaml:"round"`
	Code    string            `yaml:"code"`
	Results map[string]string `yaml:"results"` //status for each test like OK, Bad result, timeout etc
	Points  map[string]uint   `yaml:"points"`  //points for each testgroup
	Sum     uint              `yaml:"sum"`     //sum of points
}

//LoadSubmission load submission data from yaml string
func LoadSubmission(buff []byte) *Submission {
	var s Submission
	err := yaml.Unmarshal(buff, &s)
	if err != nil {
		Log(ERR, "Error parsing user submission: %s", buff)
		return nil
	}
	return &s
}
