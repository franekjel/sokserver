package executor

//Command holds user command data (unmarshallized from yaml)
type Command struct {
	Login    string `yaml:"login"`    //user login
	Password string `yaml:"password"` //user password - deleted after check
	Command  string `yaml:"command"`  //given command - described in protocol.md
	Contest  string `yaml:"contest"`  //contest name for commands that need it
	Round    string `yaml:"round"`    //round name for commands that need it
	Task     string `yaml:"task"`     //task name for commands that need it
	Code     string `yaml:"code"`     //solution code if it is submission
}

//Execute given command. Return response to the client
func Execute(buff []byte) []byte {
	return nil
}
