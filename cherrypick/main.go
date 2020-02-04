package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
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
		if line == "" {
			continue
		}
		fmt.Println("******** Cherry-picking", line)
		if output, err := runCmd("git", "cherry-pick", "-m", "1", "--strategy=recursive", "-X", "theirs", line); err != nil {
			fmt.Println(output)
			if strings.Contains(output, "CONFLICT") {
				panic(err)
			}
			fmt.Println("Empty commit, aborting and re-cherrypicking it")
			if _, err := runCmd("git", "cherry-pick", "--skip"); err != nil {
				panic(err)
			}
		}
	}
}

func runCmd(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	stdout, stderr := bytes.NewBuffer(nil), bytes.NewBuffer(nil)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.Dir = "/Users/silvin/dev/go/src/github.com/silvin-lubecki/engine-extract"
	err := cmd.Run()
	return fmt.Sprintf("Out: %s\nErr %s\n", stdout.String(), stderr.String()), err
}
