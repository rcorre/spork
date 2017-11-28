package main

import (
	"io/ioutil"

	"github.com/ghodss/yaml"
)

// Config contains configuration options for spork
type Config struct {
	SparkURL       string
	SparkDeviceURL string
}

func LoadConfig(path string) (*Config, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var conf Config
	err = yaml.Unmarshal(bytes, &conf)
	return &conf, err
}
