package main

import (
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
