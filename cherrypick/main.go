package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

func main() {
	path := os.Args[1]
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	lines := bytes.Split(buf, []byte{'\n'})
	for i := len(lines) - 1; i >= 0; i-- {
		line := string(lines[i])
		fmt.Println("******** Cherry-picking", line)
		cmd := exec.Command("git", "cherry-pick", "--strategy=recursive", "-X", "theirs", line)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Dir = "/Users/silvin/dev/go/src/github.com/silvin-lubecki/cli-extract"
		if err := cmd.Run(); err != nil {
			panic(err)
		}
	}
}
