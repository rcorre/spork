package main

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ConfigTestSuite struct {
	suite.Suite
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}

func (suite *ConfigTestSuite) TestLoadConfig() {
	tmpfile, err := ioutil.TempFile("", "sporktest")
	suite.Nil(err)

	_, err = tmpfile.WriteString(strings.Replace(`
		sparkURL: "http://example.com"
		sparkDeviceURL: "http://example.com/devices"
	`, "\t", "    ", -1))
	suite.Nil(err)

	expected := &Config{
		SparkURL:       "http://example.com",
		SparkDeviceURL: "http://example.com/devices",
	}
	actual, err := LoadConfig(tmpfile.Name())
	suite.Nil(err)
	suite.Equal(expected, actual)
}
