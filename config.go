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
	Branch                string `yaml:"Branch"`
	Component             string `yaml:"Component"`
	Ancestor              string `yaml:"Ancestor"`
}

func loadConfig(path string) Config {
	data, err := ioutil.ReadFile(path)
	checkErr(err)

	var config Config
	checkErr(yaml.Unmarshal(data, &config))
	return config
}
