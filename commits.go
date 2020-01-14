package main

import (
	"strings"

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

func FindCommitOnComponent(dockerCe, component *Remote, tagName, componentName string) (*object.Commit, *object.Commit, []*object.Commit, error) {
	// find commit related to selected tag
	ref, err := dockerCe.FindReference(tagName)
	if err != nil {
		return nil, nil, nil, err
	}
	tag, err := dockerCe.GetTagFromRef(ref)
	if err != nil {
		return nil, nil, nil, err
	}
	// find latest mege commit comming from bot merging component
	componentMergeCommit, skipped, err := dockerCe.FindLatestCommonAncestor(tag.Commit, componentName)
	if err != nil {
		return nil, nil, nil, err
	}
	// clean that message
	cleanedMessage := CleanCommitMessage(componentMergeCommit.Message)
	// find that message in the upstream repo
	componentCommit, err := component.FindCommitByMessage(cleanedMessage)
	if err != nil {
		return nil, nil, nil, err
	}
	return componentMergeCommit, componentCommit, skipped, nil
}
