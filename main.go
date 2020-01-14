package main

import (
	"fmt"
	"os"
)

func checkErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: docker-ce-tags [diff-tags|commits|all-tags] config.yml\n")
		os.Exit(1)
	}
	// Load config and remotes
	config := loadConfig(os.Args[2])
	dockerCe, err := NewRemote(config.UpstreamPath, config.UpstreamRemote)
	checkErr(err)
	component, err := NewRemote(config.ComponentPath, config.ComponentRemote)
	checkErr(err)

	switch os.Args[1] {
	case "diff-tags":
		tags, err := computeDiffTags(dockerCe, component)
		checkErr(err)
		for _, tag := range tags {
			fmt.Println(tag.Name)
		}

	case "commits":
		dockerCeCommit, componentCommit, err := FindCommitOnComponent(dockerCe, component, config.Tag, config.Component)
		checkErr(err)
		fmt.Println("***** docker/docker-ce")
		fmt.Println(dockerCeCommit)
		fmt.Println("***** docker/cli")
		fmt.Println(componentCommit)

	case "all-tags":
		tags, err := computeDiffTags(dockerCe, component)
		checkErr(err)
		for _, tag := range tags {
			dockerCeCommit, componentCommit, err := FindCommitOnComponent(dockerCe, component, tag.Name, config.Component)
			checkErr(err)
			if dockerCeCommit == nil || componentCommit == nil {
				checkErr(fmt.Errorf("%q failed to get commits", tag.Name))
			}
		}

	default:
		fmt.Fprintf(os.Stderr, "Usage: docker-ce-tags [diff-tags|commits|all-tags] config.yml\n")
		os.Exit(1)
	}
}
