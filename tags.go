package main

import (
	"fmt"
	"sort"
	"strings"

	"gopkg.in/src-d/go-git.v4"
)

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
