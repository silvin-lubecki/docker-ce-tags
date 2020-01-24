package main

import (
	"fmt"
	"os"
	"strings"
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

	switch os.Args[1] {
	case "diff-tags":
		tags, err := computeDiffTags(dockerCe, component)
		checkErr(err)
		for _, tag := range tags {
			fmt.Println(tag.Name)
		}

	case "branch":
		cherryPicked, ancestor, err := CherryPickOnBranch(dockerCe, conf.Branch, conf.Component)
		checkErr(err)
		fmt.Println("Branch", conf.Branch)
		fmt.Println("Commits to cherry pick")
		for _, c := range cherryPicked {
			fmt.Println(c)
		}

		fmt.Println("Common ancestor on docker-ce between master and branch", conf.Branch)
		fmt.Println(ancestor)

		// find latest merge commit comming from bot merging component
		botMergeCommit, err := dockerCe.FindLatestCommonAncestor(ancestor, conf.Component)
		checkErr(err)
		componentMergeCommit, err := GetLastParent(botMergeCommit)
		checkErr(err)
		// clean that message
		cleanedMessage := CleanCommitMessage(componentMergeCommit.Message)
		// find that message in the upstream repo
		componentCommit, err := component.FindCommitByMessage(cleanedMessage)
		checkErr(err)

		// Now find the commit in the "git filter-branch" extracted branch on docker-ce
		extractHead, err := dockerCe.GetHead(conf.Branch + "-extract")
		checkErr(err)
		fmt.Println("Extracted Head", extractHead.Hash)

		extractedCommit, err := dockerCe.FindCommitByMessageOnBranch(botMergeCommit.Message, extractHead.Hash)
		checkErr(err)
		fmt.Println("******** Extracted Commit Ancestor")
		fmt.Println(extractedCommit)

		fmt.Println("Commit to branch for", conf.Branch)
		fmt.Println(componentCommit)

		fmt.Println("****************")

		cherryPicked, err = dockerCe.FindCommitsToCherryPick(extractHead, extractedCommit)
		checkErr(err)

		for _, c := range cherryPicked {
			fmt.Println(c.Hash)
			//extractedCommit, err := dockerCe.FindCommitByMessageOnBranch(c.Message, extractHead.Hash)
			//checkErr(err)
			//fmt.Println(extractedCommit.Hash)
		}

		// // Create a branch on the component fork (delete it if already exists)
		// fmt.Println("Creating the branch on the extracted fork")

		// w, err := component.repo.Worktree()
		// checkErr(err)
		// checkErr(w.Checkout(&git.CheckoutOptions{
		// 	Hash:   componentCommit.Hash,
		// 	Branch: plumbing.NewBranchReferenceName(conf.Branch),
		// 	Create: true,
		// }))

		// head := componentCommit.Hash

		// fmt.Println("Cherry picking the commits")
		// for i := len(cherryPicked) - 1; i >= 0; i-- {
		// 	commitHead, err := component.repo.CommitObject(head)
		// 	checkErr(err)

		// 	c := cherryPicked[i]
		// 	fmt.Println("Cherrypicking", c.Hash, "HEAD commit", head)
		// 	t, err := c.Tree()
		// 	checkErr(err)
		// 	fromTree, err := t.Tree("components/" + conf.Component)
		// 	checkErr(err)
		// 	toTree, err := commitHead.Tree()
		// 	checkErr(err)

		// 	// if c.Hash.String() == "7d395933ee04cc9b981590d2f65cf0aa80fa9694" {
		// 	// 	f, err := fromTree.File("docs/reference/commandline/run.md")
		// 	// 	checkErr(err)
		// 	// 	content, err := f.Contents()
		// 	// 	checkErr(err)
		// 	// 	fmt.Println("docker ce content")
		// 	// 	fmt.Println(content)

		// 	// 	f, err = toTree.File("docs/reference/commandline/run.md")
		// 	// 	checkErr(err)
		// 	// 	content, err = f.Contents()
		// 	// 	checkErr(err)
		// 	// 	fmt.Println("cli content")
		// 	// 	fmt.Println(content)
		// 	// }

		// 	// patch, err := toTree.Patch(fromTree)
		// 	// checkErr(err)
		// 	diffs, err := toTree.Diff(fromTree)
		// 	checkErr(err)
		// 	for _, diff := range diffs {
		// 		fmt.Println(diff.String())
		// 		_, toFile, err := diff.Files()
		// 		checkErr(err)

		// 		action, err := diff.Action()
		// 		checkErr(err)
		// 		toPath := filepath.Join(conf.ComponentPath, diff.To.Name)
		// 		fromPath := filepath.Join(conf.ComponentPath, diff.From.Name)
		// 		// Applying action
		// 		switch action {
		// 		case merkletrie.Insert:
		// 			fmt.Println("Inserting")
		// 			{
		// 				rTo, err := toFile.Reader()
		// 				checkErr(err)
		// 				bTo, err := ioutil.ReadAll(rTo)
		// 				checkErr(err)
		// 				checkErr(os.MkdirAll(filepath.Dir(toPath), 0755))
		// 				f, err := os.Create(toPath)
		// 				checkErr(err)
		// 				defer f.Close()
		// 				_, err = f.Write(bTo)
		// 				checkErr(err)
		// 			}
		// 			_, err := w.Add(diff.To.Name)
		// 			checkErr(err)

		// 		case merkletrie.Delete:
		// 			fmt.Println("Deleting", fromPath)
		// 			checkErr(os.Remove(fromPath))
		// 			_, err := w.Add(diff.From.Name)
		// 			checkErr(err)

		// 		case merkletrie.Modify:
		// 			fmt.Println("Modifying", toPath)
		// 			if fromPath != toPath {
		// 				fmt.Println("Moving", fromPath)
		// 				checkErr(os.Rename(fromPath, toPath))
		// 				_, err := w.Add(diff.From.Name)
		// 				checkErr(err)
		// 			}
		// 			{
		// 				rTo, err := toFile.Reader()
		// 				checkErr(err)
		// 				bTo, err := ioutil.ReadAll(rTo)
		// 				checkErr(err)

		// 				f, err := os.Create(toPath)
		// 				checkErr(err)
		// 				defer f.Close()
		// 				_, err = f.Write(bTo)
		// 				checkErr(err)
		// 			}
		// 			_, err := w.Add(diff.To.Name)
		// 			if err != nil {
		// 				fmt.Println("worktree add failed", toPath)
		// 			}
		// 			checkErr(err)
		// 		}
		// 	}
		// 	// Now commit
		// 	fmt.Println("Status")
		// 	status, err := w.Status()
		// 	checkErr(err)
		// 	fmt.Println(status)

		// 	head, err = w.Commit(c.Message, &git.CommitOptions{
		// 		All:       true,
		// 		Author:    &c.Author,
		// 		Committer: &c.Committer,
		// 	})
		// 	checkErr(err)
		// 	fmt.Println("Commited", head)

		// 	fmt.Println("Status after")
		// 	status, err = w.Status()
		// 	checkErr(err)
		// 	fmt.Println(status)

		// 	commit, err := component.repo.CommitObject(head)
		// 	checkErr(err)

		// 	// Check commit
		// 	stats, err := commit.Stats()
		// 	checkErr(err)
		// 	fmt.Println("Component commit stats")
		// 	fmt.Println(stats)
		// 	cherryPickedStats, err := c.Stats()
		// 	checkErr(err)
		// 	fmt.Println("DockerCe commit stats")
		// 	fmt.Println(cherryPickedStats)

		// 	for _, s := range stats {
		// 		var found bool
		// 		for _, cs := range cherryPickedStats {
		// 			if strings.Contains(cs.Name, s.Name) && cs.Addition == s.Addition && cs.Deletion == s.Deletion {
		// 				found = true
		// 			}
		// 		}
		// 		if !found {
		// 			checkErr(fmt.Errorf("Failed to find %q in %s", s.String(), cherryPickedStats))
		// 		}
		// 	}
		// }

		fmt.Println("Done")

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
