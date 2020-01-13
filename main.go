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
		fmt.Fprintf(os.Stderr, "Usage: docker-ce-tags [diff-tags|commits] config.yml\n")
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
		checkErr(printDiffTags(dockerCe, component))
	case "commits":
		commits, err := dockerCe.GetCommits(config.Tag)
		checkErr(err)
		for _, c := range commits {
			tree, err := c.Tree()
			checkErr(err)
			cli, err := tree.Tree("components/" + config.Component)
			if err == nil {
				fmt.Println(len(cli.Entries))
			}
		}
	default:
		fmt.Fprintf(os.Stderr, "Usage: docker-ce-tags [diff-tags|commits] config.yml\n")
		os.Exit(1)
	}
}
