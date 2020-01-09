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
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Need a path to config file\n")
		os.Exit(1)
	}
	config := loadConfig(os.Args[1])
	printDiffTags(config)
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
