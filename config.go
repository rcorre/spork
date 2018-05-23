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
	MessageFormat  messageFormat
}

type messageFormat struct {
	emph    string
	strong  string
	code    string
	quote   string
	link    string
	mention string
}

func defaultConfig() *Config {
	return &Config{
		SparkURL:       "https://api.ciscospark.com/v1/",
		SparkDeviceURL: "https://wdm-a.wbx2.com/wdm/api/v1/devices",
		Keys: map[string]string{
			"<c-c>":   "quit",
			"<c-j>":   "nextroom",
			"<c-k>":   "prevroom",
			"<c-u>":   "halfpageup",
			"<c-d>":   "halfpagedown",
			"<enter>": "send",
		},
		MessageFormat: messageFormat{
			emph:    "white+h",
			strong:  "white+b",
			code:    "white:gray",
			quote:   "gray",
			link:    "blue",
			mention: "blue+b",
		},
	}
}

func LoadConfig(path string) (*Config, error) {
	conf := defaultConfig()
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return conf, nil
		}
		return nil, err
	}

	err = yaml.Unmarshal(bytes, &conf)
	return conf, err
}
