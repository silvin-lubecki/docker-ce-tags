package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"

	"gopkg.in/src-d/go-git.v4"
)

func checkErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

type Config struct {
	UpstreamPath   string `yaml:"UpstreamPath"`
	ComponentPath  string `yaml:"ComponentPath"`
	ExportRepoPath string `yaml:"ExportRepoPath"`
	Tag            string `yaml:"Tag"`
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Need a path to config file\n")
		os.Exit(1)
	}
	data, err := ioutil.ReadFile(os.Args[1])
	checkErr(err)

	var config Config
	checkErr(yaml.Unmarshal(data, &config))
	printDiffTags(config)
}

func printDiffTags(config Config) {
	dockerCeTags := getTags(config.UpstreamPath)
	componentTags := getTags(config.ComponentPath)

	tagsToAdd := diffTags(dockerCeTags, componentTags)

	for _, tag := range tagsToAdd {
		fmt.Println(tag)
	}
}

func getTags(path string) []string {
	r, err := git.PlainOpen(path)
	checkErr(err)

	remotes, err := r.Remotes()
	checkErr(err)
	var upstream *git.Remote
	for _, remote := range remotes {
		for _, url := range remote.Config().URLs {
			if strings.Contains(url, "docker/") {
				upstream = remote
			}
		}
	}
	refs, err := upstream.List(&git.ListOptions{})
	checkErr(err)
	var tags []string
	for _, ref := range refs {
		if ref.Name().IsTag() {
			tags = append(tags, ref.Name().Short())
		}
	}
	sort.Strings(tags)
	return tags
}

func diffTags(upstreamTags, componentTags []string) []string {
	var tagsToAdd []string
	for _, up := range upstreamTags {
		found := false
		for _, comp := range componentTags {
			if up == comp {
				found = true
				break
			}
		}
		if !found {
			tagsToAdd = append(tagsToAdd, up)
		}
	}
	return tagsToAdd
}

// ref, err := r.Head()
// if err != nil {
// 	panic(err.Error())
// }
// cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
// if err != nil {
// 	panic(err.Error())
// }
// err = cIter.ForEach(func(c *object.Commit) error {
// 	//fmt.Println(c)
// 	return nil
// })
// if err != nil {
// 	panic(err.Error())
// }
