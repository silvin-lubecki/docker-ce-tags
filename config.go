package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	UpstreamPath   string `yaml:"UpstreamPath"`
	ComponentPath  string `yaml:"ComponentPath"`
	ExportRepoPath string `yaml:"ExportRepoPath"`
	Tag            string `yaml:"Tag"`
}

func loadConfig(path string) Config {
	data, err := ioutil.ReadFile(path)
	checkErr(err)

	var config Config
	checkErr(yaml.Unmarshal(data, &config))
	return config
}
