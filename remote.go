package main

import (
	"fmt"
	"strings"

	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"

	"gopkg.in/src-d/go-git.v4"
)

type Remote struct {
	repo   *git.Repository
	remote *git.Remote
}

func NewRemote(path, remoteName string) (*Remote, error) {
	repository, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}

	remotes, err := repository.Remotes()
	if err != nil {
		return nil, err
	}
	var upstream *git.Remote
	for _, remote := range remotes {
		for _, url := range remote.Config().URLs {
			if strings.Contains(url, remoteName) {
				upstream = remote
			}
		}
	}
	if upstream == nil {
		return nil, fmt.Errorf("remote %q not found", remoteName)
	}

	return &Remote{repository, upstream}, nil
}

func (r *Remote) GetCommits(reference string) ([]*object.Commit, error) {
	ref, err := r.findReference(reference)
	if err != nil {
		return nil, err
	}
	cIter, err := r.repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return nil, err
	}
	var commits []*object.Commit
	if err := cIter.ForEach(func(c *object.Commit) error {
		commits = append(commits, c)
		return nil
	}); err != nil {
		return nil, err
	}
	return commits, nil
}

func (r *Remote) findReference(name string) (*plumbing.Reference, error) {
	refs, err := r.remote.List(&git.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, ref := range refs {
		if ref.Name().Short() == name {
			return ref, nil
		}
	}
	return nil, fmt.Errorf("Reference %q not found", name)
}

func (r *Remote) GetTags() ([]Tag, error) {
	var tags []Tag
	iter, err := r.repo.Tags()
	if err != nil {
		return nil, err
	}
	if err := iter.ForEach(func(ref *plumbing.Reference) error {
		t, err := r.repo.TagObject(ref.Hash())
		var c *object.Commit
		// Not a tag? but commit object
		if err != nil {
			c, err = r.repo.CommitObject(ref.Hash())
			if err != nil {
				fmt.Printf("commit %q not found %s\n", ref.Name().Short(), err)
				return nil
			}
		} else {
			c, err = t.Commit()
			if err != nil {
				return fmt.Errorf("tag commit %q not found %s", ref.Name().Short(), err)
			}
		}
		tags = append(tags, Tag{ref.Name().Short(), c})
		return nil
	}); err != nil {
		return nil, err
	}
	return tags, nil
}
