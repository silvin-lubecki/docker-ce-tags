package main

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func checkErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: docker-ce-tags [diff-tags|branch|all-tags] conf.yml\n")
		os.Exit(1)
	}
	// Load config and remotes
	conf := loadConfig(os.Args[2])
	dockerCe, err := NewRemote(conf.UpstreamPath, conf.UpstreamRemote)
	checkErr(err)
	component, err := NewRemote(conf.ComponentPath, conf.ComponentRemote)
	checkErr(err)
	result, err := NewRemote("/Users/silvin/dev/go/src/github.com/silvin-lubecki/packaging-extract", "silvin-lubecki/packaging-extract")
	checkErr(err)

	switch os.Args[1] {
	case "diff-tags":
		tags, err := computeDiffTags(dockerCe, component)
		checkErr(err)
		for _, tag := range tags {
			fmt.Printf("| %s | ", tag.Name)
			branchName := fmt.Sprintf("%s-extract-%s", tag.Name[1:6], conf.Component)
			branch, err := result.GetHead(branchName)
			if err != nil {
				branchName := tag.Name[1:6]
				branch, err = result.GetHead(branchName)
				checkErr(err)
			}
			var commit *object.Commit
			for commit == nil {
				c := tag.Commit
				if strings.HasPrefix(c.Message, "Merge pull request") {
					c, err = GetLastParent(c)
					checkErr(err)
				} else if strings.HasPrefix(c.Message, "Merge component '"+conf.Component) {
					c, err = GetLastParent(c)
					checkErr(err)
					c.Message = CleanCommitMessage(c.Message)
				}
				m := c.Message
				if len(m) >= 25 {
					m = m[:25]
					m = strings.ReplaceAll(m, "\n", "")
				}
				commit, err = result.FindCommitByMessageOnBranch(c.Message, branch.Hash)
				if commit == nil {
					c, err := GetFirstParent(tag.Commit)
					checkErr(err)
					tag.Commit = c
				}
			}
			fmt.Println(tag.Commit.Hash, "|", commit.Hash, "|")
		}

	case "branch":
		extractHead, err := dockerCe.GetHead(conf.Branch + "-extract")
		checkErr(err)

		var ancestor *object.Commit
		if conf.Ancestor != "" {
			ancestor, err = dockerCe.GetCommit(conf.Ancestor, "master")
			checkErr(err)
		} else {
			ancestor, err = CherryPickOnBranch(dockerCe, conf.Branch, conf.Component)
			checkErr(err)
		}
		fmt.Println("Ancestor", ancestor.Hash)
		//fmt.Println("Branch", conf.Branch)
		//fmt.Println("Commits to cherry pick")
		//for _, c := range cherryPicked {
		//	fmt.Println(c)
		//}

		//fmt.Println("Common ancestor on docker-ce between master and branch", conf.Branch)
		//fmt.Println(ancestor)

		// find latest merge commit comming from bot merging component
		botMergeCommit, err := dockerCe.FindLatestCommonAncestor(ancestor, conf.Component)
		checkErr(err)
		fmt.Println("botMergeCommit", botMergeCommit.Hash)
		dockerCEMergeCommit, err := GetLastParent(botMergeCommit)
		checkErr(err)
		fmt.Println("dockerCEMergeCommit", dockerCEMergeCommit.Hash)

		// clean that message
		cleanedMessage := CleanCommitMessage(dockerCEMergeCommit.Message)
		// find that message in the upstream repo
		dockerProductCommit, err := component.FindCommitByMessage(cleanedMessage)
		checkErr(err)
		fmt.Println("dockerProductCommit", dockerProductCommit.Hash)
		// Now find the commit in the "git filter-branch" extracted branch on docker-ce
		//extractHead, err := dockerCe.GetHead(conf.Branch + "-extract")
		//	checkErr(err)
		fmt.Println("Extracted Head", extractHead.Hash)
		fmt.Println("Bot Merge commit", botMergeCommit.Hash)
		fmt.Println("Component Merge Commit", dockerCEMergeCommit)

		extractedCommit, err := dockerCe.FindCommitByMessageOnBranch(dockerCEMergeCommit.Message, extractHead.Hash)
		if err != nil {
			parent, err := GetLastParent(dockerCEMergeCommit)
			checkErr(err)
			extractedCommit, err = dockerCe.FindCommitByMessageOnBranch(parent.Message, extractHead.Hash)
		}

		fmt.Println("Extracted Commit", extractedCommit.Hash)

		//fmt.Println("******** Extracted Commit Ancestor")
		//fmt.Println(extractedCommit)
		fmt.Printf("git checkout -b %s-extract-%s %s\n", conf.Branch, conf.Component, dockerProductCommit.Hash)

		cherryPicked, err := dockerCe.FindCommitsToCherryPick(extractHead, extractedCommit)
		checkErr(err)

		f, err := os.Create("commits")
		checkErr(err)
		defer f.Close()
		for _, c := range cherryPicked {
			// found, err := dockerCe.FindCommitByMessageOnBranch(c.Message, extractHead.Hash)
			// if err != nil {
			// 	fmt.Println("*** Not found")
			// 	fmt.Println(c)
			// 	checkErr(err)
			// }
			// fmt.Println(c.Hash, found.Hash)
			fmt.Fprintln(f, c)
		}

	case "all-tags":
		tags, err := computeDiffTags(dockerCe, component)
		checkErr(err)
		for _, tag := range tags {
			fmt.Println("Checking", tag.Name)
			dockerCeCommit, componentCommit, err := FindCommitOnComponent(dockerCe, component, tag.Name, conf.Component)
			checkErr(err)
			if dockerCeCommit == nil || componentCommit == nil {
				checkErr(fmt.Errorf("%q failed to get commits", tag.Name))
			}
			fmt.Println("docker/docker-ce", dockerCeCommit.Hash, "docker/"+conf.Component, componentCommit.Hash)
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
		fmt.Fprintf(os.Stderr, "Usage: docker-ce-tags [diff-tags|branch|all-tags] conf.yml\n")
		os.Exit(1)
	}
}
