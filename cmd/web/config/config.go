package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type TwitterConfig struct {
	UserId string `yaml:"userId"`
	Bearer string `yaml:"bearer"`
}

type Config struct {
	Twitter *TwitterConfig `yaml:"twitter"`
}

var conf *Config

func NewConfigFile(filePath string) (*Config, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	conf = &Config{}
	if err = decoder.Decode(conf); err != nil {
		return nil, err
	}

	return conf, nil
}

func GetTwitterConfig() TwitterConfig {
	return *conf.Twitter
}
