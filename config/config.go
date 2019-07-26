package config

import (
	"gopkg.in/yaml.v2"
)

//Config stores global config variables
type Config struct {
	Port               uint16 `yaml:"port"`
	Workers            uint16 `yaml:"workers"`
	DefaultMemoryLimit uint32 `yaml:"default_memory_limit"`
	DefaultTimeLimit   uint16 `yaml:"default_time_limit"`
	Info               string `yaml:"info"`
}

//MakeConfig loads variables from yaml buff to conf or sets defaults
func MakeConfig(buff *[]byte) *Config {
	conf := Config{
		Port:               19151,
		Workers:            4,
		DefaultMemoryLimit: 32 * 1024,
		DefaultTimeLimit:   1000,
		Info:               "Sok server",
	}
	yaml.Unmarshal(*buff, &conf)
	return &conf
}

//GetConfig return current config in yaml format
func (conf *Config) GetConfig() []byte {
	buff, _ := yaml.Marshal(conf)
	return buff
}
