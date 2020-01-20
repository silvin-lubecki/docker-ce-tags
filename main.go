package main

import (
	"fmt"
	"os"
	"strings"
)

func checkErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: docker-ce-tags [diff-tags|branch|all-tags] config.yml\n")
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

	case "branch":
		cherryPicked, err := CherryPickOnBranch(dockerCe, config.Branch, config.Component)
		checkErr(err)
		fmt.Println("Branch", config.Branch)
		fmt.Println("Commits to cherry pick")
		for _, c := range cherryPicked {
			fmt.Println(c)
		}

	case "all-tags":
		tags, err := computeDiffTags(dockerCe, component)
		checkErr(err)
		for _, tag := range tags {
			fmt.Println("Checking", tag.Name)
			dockerCeCommit, componentCommit, err := FindCommitOnComponent(dockerCe, component, tag.Name, config.Component)
			checkErr(err)
			if dockerCeCommit == nil || componentCommit == nil {
				checkErr(fmt.Errorf("%q failed to get commits", tag.Name))
			}
			fmt.Println("docker/docker-ce", dockerCeCommit.Hash, "docker/"+config.Component, componentCommit.Hash)
		}

	case "sign-off":
		commits, err := dockerCe.GetCommits("master")
		checkErr(err)
		for _, c := range commits {
			if !strings.Contains(c.Message, "Signed-off-by:") && !strings.Contains(c.Message, "Merge pull request") {
				fmt.Println(c)
			}
		}

	default:
		fmt.Fprintf(os.Stderr, "Usage: docker-ce-tags [diff-tags|branch|all-tags] config.yml\n")
		os.Exit(1)
	}
}
