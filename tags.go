package main

import (
	"fmt"

	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

type Tag struct {
	Name   string
	Commit *object.Commit
}

func printDiffTags(dockerCe, component *Remote) error {
	dockerCeTags, err := dockerCe.GetTags()
	if err != nil {
		return err
	}
	componentTags, err := component.GetTags()
	if err != nil {
		return err
	}

	tagsToAdd := diffTags(dockerCeTags, componentTags)
	for _, tag := range tagsToAdd {
		fmt.Println(tag.Name)
	}
	return nil
}

func diffTags(upstreamTags, componentTags []Tag) []Tag {
	var tagsToAdd []Tag
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
