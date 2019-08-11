package contests

import (
	"gopkg.in/yaml.v2"

	"github.com/franekjel/sokserver/fs"
	log "github.com/franekjel/sokserver/logger"
	"github.com/franekjel/sokserver/rounds"
	"github.com/franekjel/sokserver/tasks"
)

//Contest holds rounds in contest and groups allowed to participate
type Contest struct {
	fs      *fs.Fs
	Rounds  map[string]*rounds.Round `yaml:"-"`
	Ranking map[string]uint          `yaml:"-"`
	Name    string                   `yaml:"name"`
	Key     string                   `yaml:"key"`
}

func (c *Contest) loadRanking() {
	c.Ranking = make(map[string]uint)
	for _, round := range c.Rounds {
		for i, score := range round.Ranking.Points {
			_, ok := c.Ranking[round.Ranking.Names[i]]
			if ok {
				c.Ranking[round.Ranking.Names[i]] += score[0]
			} else {
				c.Ranking[round.Ranking.Names[i]] = score[0]
			}
		}
	}
}

func (c *Contest) loadRounds(tasks map[string]*tasks.Task) {
	c.Rounds = make(map[string]*rounds.Round)
	for _, round := range c.fs.ListDirs("") {
		c.Rounds[round] = rounds.LoadRound(fs.Init(c.fs.Path, round), tasks)
	}
}

func (c *Contest) loadConfig() {
	if !c.fs.FileExist("contest.yml") {
		log.Fatal("Contest settings missing! %s", c.fs.Path)
	}
	buff := c.fs.ReadFile("contest.yml")
	err := yaml.Unmarshal(buff, c)
	if err != nil {
		log.Error("%s", err.Error())
	}
}

//LoadContest loads contest in given folder
func LoadContest(fs *fs.Fs, tasks map[string]*tasks.Task) *Contest {
	log.Info("Loading contest %s", fs.Path)
	contest := new(Contest)
	contest.fs = fs
	contest.loadConfig()
	contest.loadRounds(tasks)
	contest.loadRanking()
	return contest
}
