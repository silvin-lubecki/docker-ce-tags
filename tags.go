package main

import (
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

type Tag struct {
	Name   string
	Commit *object.Commit
}

func computeDiffTags(dockerCe, component *Remote) ([]Tag, error) {
	dockerCeTags, err := dockerCe.GetTags()
	if err != nil {
		return nil, err
	}
	componentTags, err := component.GetTags()
	if err != nil {
		return nil, err
	}

	return diffTags(dockerCeTags, componentTags), nil
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
