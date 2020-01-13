package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	UpstreamPath          string `yaml:"UpstreamPath"`
	UpstreamRemote        string `yaml:"UpstreamRemote"`
	ComponentPath         string `yaml:"ComponentPath"`
	ComponentRemote       string `yaml:"ComponentRemote"`
	ComponentResultRemote string `yaml:"ComponentResultRemote"`
	Tag                   string `yaml:"Tag"`
	Component             string `yaml:"Component"`
}

func loadConfig(path string) Config {
	data, err := ioutil.ReadFile(path)
	checkErr(err)

	var config Config
	checkErr(yaml.Unmarshal(data, &config))
	return config
}
