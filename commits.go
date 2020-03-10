package main

import (
	"fmt"
	"strings"

	"gopkg.in/src-d/go-git.v4/plumbing/storer"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func findCommonCommits(config Config) []*object.Commit {
	return nil
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

// commits, err := dockerCe.GetCommits(config.Tag)
// checkErr(err)
// for _, c := range commits {
// 	tree, err := c.Tree()
// 	checkErr(err)
// 	cli, err := tree.Tree("components/" + config.Component)
// 	if err == nil {
// 		fmt.Println(len(cli.Entries))
// 	}
// }

func CleanCommitMessage(ceMessage string) string {
	index := strings.Index(ceMessage, "\nUpstream-commit:")
	if index == -1 {
		return ceMessage
	}
	return ceMessage[:index]
}

func GetFirstParent(commit *object.Commit) (*object.Commit, error) {
	var parent *object.Commit
	err := commit.Parents().ForEach(func(c *object.Commit) error {
		if parent == nil {
			parent = c
		}
		return nil
	})
	return parent, err
}

func GetLastParent(commit *object.Commit) (*object.Commit, error) {
	var parent *object.Commit
	err := commit.Parents().ForEach(func(c *object.Commit) error {
		parent = c
		return nil
	})
	return parent, err
}

func FindCommitOnComponent(dockerCe, component *Remote, tagName, componentName string) (*object.Commit, *object.Commit, error) {
	// find commit related to selected tag
	tag, err := getTag(dockerCe, tagName)
	if err != nil {
		return nil, nil, err
	}

	// find latest merge commit comming from bot merging component
	componentMergeCommit, err := dockerCe.FindLatestCommonAncestor(tag.Commit, componentName)
	if err != nil {
		return nil, nil, err
	}
	// clean that message
	cleanedMessage := CleanCommitMessage(componentMergeCommit.Message)
	// find that message in the upstream repo
	componentCommit, err := component.FindCommitByMessage(cleanedMessage)
	if err != nil {
		return nil, nil, err
	}
	return componentMergeCommit, componentCommit, nil
}

func CherryPickOnBranch(dockerCe *Remote, branchName, component string) (*object.Commit, error) {
	// First find common ancestor between master branch and targeted branch
	master, err := dockerCe.GetHead("master")
	if err != nil {
		return nil, err
	}
	head, err := dockerCe.GetHead(branchName)
	if err != nil {
		return nil, err
	}

	cIter, err := dockerCe.repo.Log(&git.LogOptions{From: head.Hash})
	if err != nil {
		return nil, err
	}
	var commonAncestor *object.Commit
	if err := cIter.ForEach(func(c *object.Commit) error {
		fmt.Println("Checking ancestor with", c.Hash)
		if b, err := c.IsAncestor(master); err == nil && b {
			commonAncestor = c
			fmt.Printf("Ancestor %q found for branch %q\n", c.Hash, branchName)
			return storer.ErrStop
		}
		return nil
	}); err != nil {
		return nil, err
	}
	if commonAncestor == nil {
		return nil, fmt.Errorf("Could not find a common ancestor between %q and %q", "master", branchName)
	}
	// find all commits made only on DockerCE
	// commits, err := dockerCe.FindCommitsToCherryPick(head, commonAncestor)
	// if err != nil {
	// 	fmt.Println("error in FindCommitsToCherryPick", err)
	// 	return nil, nil, err
	// }
	return commonAncestor, nil
}

func getTag(remote *Remote, tagName string) (Tag, error) {
	ref, err := remote.FindReference(tagName)
	if err != nil {
		return Tag{}, err
	}
	tag, err := remote.GetTagFromRef(ref)
	if err != nil {
		return Tag{}, err
	}
	return tag, nil
}
