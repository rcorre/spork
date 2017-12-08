package main

import (
	"io/ioutil"
	"os"

	"github.com/ghodss/yaml"
)

// Config contains configuration options for spork
type Config struct {
	SparkURL       string
	SparkDeviceURL string
	Keys           map[string]string
}

func LoadConfig(path string) (*Config, error) {
	conf := Config{
		Keys: map[string]string{
			"<c-c>":   "quit",
			"<c-j>":   "nextroom",
			"<c-k>":   "prevroom",
			"<c-u>":   "halfpageup",
			"<c-d>":   "halfpagedown",
			"<enter>": "send",
		},
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, err
	}

	err = yaml.Unmarshal(bytes, &conf)
	return &conf, err
}
