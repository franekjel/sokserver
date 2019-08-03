package submissions

//Submission holds users submissions data
type Submission struct {
	User       string
	Task       string
	Round      string
	Code       *string
	Results    map[string]string //status for each test like OK, Bad result, timeout etc
	TestPoints map[string]uint   //points for each testgroup
	Points     uint              //sum of points
}

//LoadSubmission load submission data from yaml string
func LoadSubmission(buff []byte) *Submission {
	//TODO
	return nil
}
