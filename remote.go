package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/storer"

	"gopkg.in/src-d/go-git.v4"
)

type Remote struct {
	repo       *git.Repository
	remote     *git.Remote
	remoteName string
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

	return &Remote{repository, upstream, remoteName}, nil
}

func (r *Remote) GetCommits(reference string) ([]*object.Commit, error) {
	ref, err := r.FindReference(reference)
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

func (r *Remote) GetCommit(hash, refName string) (*object.Commit, error) {
	ref, err := r.FindReference(refName)

	if err != nil {
		return nil, err
	}
	cIter, err := r.repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return nil, err
	}
	var commit *object.Commit
	if err := cIter.ForEach(func(c *object.Commit) error {
		if c.Hash.String() == hash {
			commit = c
			return storer.ErrStop
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return commit, nil
}

func (r *Remote) GetCommitsFromHash(h plumbing.Hash) ([]*object.Commit, error) {
	cIter, err := r.repo.Log(&git.LogOptions{From: h})
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

func (r *Remote) GetHead(branchName string) (*object.Commit, error) {
	// Local branch
	if branchRef, err := r.repo.Branch(branchName); err == nil && branchRef != nil {
		ref, err := r.repo.Reference(branchRef.Merge, true)
		if err != nil {
			return nil, err
		}
		return r.repo.CommitObject(ref.Hash())
	}

	// remote branch
	refs, err := r.repo.Storer.IterReferences()
	if err != nil {
		return nil, err
	}
	var branch *plumbing.Reference
	if err := refs.ForEach(func(ref *plumbing.Reference) error {
		if ref.Name().Short() == fmt.Sprintf("docker/%s", branchName) {
			branch = ref
			return storer.ErrStop
		}
		return nil
	}); err != nil {
		return nil, err
	}
	if branch == nil {
		return nil, fmt.Errorf("branch %q not found", branchName)
	}
	return r.repo.CommitObject(branch.Hash())
}

func (r *Remote) FindReference(name string) (*plumbing.Reference, error) {
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
		tag, err := r.GetTagFromRef(ref)
		if err != nil {
			return err
		}
		tags = append(tags, tag)
		return nil
	}); err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *Remote) GetTagFromRef(ref *plumbing.Reference) (Tag, error) {
	t, err := r.repo.TagObject(ref.Hash())
	var c *object.Commit
	// Not a tag? but commit object
	if err != nil {
		c, err = r.repo.CommitObject(ref.Hash())
		if err != nil {
			return Tag{}, fmt.Errorf("commit %q not found %s\n", ref.Name().Short(), err)
		}
	} else {
		c, err = t.Commit()
		if err != nil {
			return Tag{}, fmt.Errorf("tag commit %q not found %s", ref.Name().Short(), err)
		}
	}
	return Tag{ref.Name().Short(), c}, nil
}

func (r *Remote) FindCommitByMessage(message string) (*object.Commit, error) {
	cIter, err := r.repo.Log(&git.LogOptions{All: true, Order: git.LogOrderCommitterTime})
	if err != nil {
		return nil, err
	}
	var found *object.Commit
	if err := cIter.ForEach(func(c *object.Commit) error {
		if strings.HasPrefix(c.Message, message) {
			found = c
			return storer.ErrStop
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return found, nil
}

func (r *Remote) FindCommitByMessageStartByBranch(originalHash plumbing.Hash, dir, branch, message string) (*object.Commit, error) {
	commits, err := r.GitLogTopoOrder(dir, branch, "", true)
	if err != nil {
		return nil, err
	}

	for _, c := range commits {
		//fmt.Println(c.Hash)
		if strings.HasPrefix(c.Message, message) {
			return c, nil
		}
	}
	return nil, fmt.Errorf("message %q not found on branch %q", message, branch)
}

func (r *Remote) GitLogTopoOrder(dir, branch, fromCommit string, nomerge bool) ([]*object.Commit, error) {
	// Checkout the branch
	cmd := exec.Command("git", "checkout", branch)
	//cmd.Stderr = os.Stderr
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	// git log topo order
	cmd2 := exec.Command("git", "log", "--topo-order", "--format=format:%H")
	if nomerge {
		cmd2.Args = append(cmd2.Args, "--no-merges")
	}
	if fromCommit != "" {
		cmd2.Args = append(cmd2.Args, fmt.Sprintf("%s..HEAD", fromCommit))
	}
	//fmt.Println("running git", cmd2.Args)
	//cmd2.Stderr = os.Stderr
	cmd2.Dir = dir
	buff := bytes.NewBuffer(nil)
	cmd2.Stdout = buff
	if err := cmd2.Run(); err != nil {
		return nil, err
	}
	hashes := strings.Split(buff.String(), "\n")
	var commits []*object.Commit
	for _, h := range hashes {
		commitHash := plumbing.NewHash(h)
		commit, err := r.repo.CommitObject(commitHash)
		if err != nil {
			return nil, fmt.Errorf("commit %q: %s", h, err)
		}
		commits = append(commits, commit)
	}
	return commits, nil
}

func (r *Remote) FindCommitByMessageOnBranch(message string, from plumbing.Hash) (*object.Commit, error) {
	cIter, err := r.repo.Log(&git.LogOptions{From: from})
	if err != nil {
		return nil, err
	}
	var found *object.Commit
	if err := cIter.ForEach(func(c *object.Commit) error {
		if strings.HasPrefix(c.Message, message) {
			found = c
			return storer.ErrStop
		}
		return nil
	}); err != nil {
		return nil, err
	}
	if found == nil {
		return nil, fmt.Errorf("Couldn't find commit with message %q starting from %q", message, from)
	}
	return found, nil
}

func (r *Remote) FindLatestCommonAncestor(initialCommit *object.Commit, component string) (*object.Commit, error) {
	current := initialCommit
	for {
		// Merge commit from bot on selected component
		if strings.Contains(current.Message, fmt.Sprintf("Merge component '%s'", component)) {
			//found, err := GetLastParent(current)
			//if err != nil {
			//	return nil, err
			//}
			return current, nil
		} else {
			// Other commit (other component, or commit on docker-ce directly)
			next, err := GetFirstParent(current)
			if err != nil {
				return nil, err
			}
			current = next
		}
	}
}

func (r *Remote) FindComponentCommit(dir, branch string, initialCommit *object.Commit, component string) (*object.Commit, error) {
	commits, err := r.GitLogTopoOrder(dir, branch, "", false)
	if err != nil {
		return nil, err
	}
	started := false
	for _, c := range commits {
		if c.Hash == initialCommit.Hash {
			started = true
		}
		if !started {
			continue
		}
		if strings.HasPrefix(c.Message, "Merge pull request") || strings.HasPrefix(c.Message, fmt.Sprintf("Merge component '%s'", component)) {
			continue
		}
		if b, err := isCommitOnComponent(c, component); b && err == nil {
			return c, nil
		}
	}
	return nil, fmt.Errorf("Couldn't find commit %q", initialCommit.Hash)
}

func (r *Remote) FindCommitsToCherryPick(initialCommit, finalCommit *object.Commit) ([]string, error) {
	current := fmt.Sprintf("%s..HEAD", finalCommit.Hash)
	cmd := exec.Command("git", "log", "--topo-order", "--no-merges", "--format=format:%H", current)
	fmt.Println("running git", cmd.Args)
	cmd.Stderr = os.Stderr
	cmd.Dir = "/Users/silvin/dev/go/src/github.com/docker/docker-ce-extract/docker-ce"
	buff := bytes.NewBuffer(nil)
	cmd.Stdout = buff
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	return strings.Split(buff.String(), "\n"), nil
}

func getComponentCommitsFromMerge(mergeCommit *object.Commit) ([]*object.Commit, error) {
	var commits []*object.Commit
	current, err := GetLastParent(mergeCommit)
	if err != nil {
		return nil, err
	}
	for {
		if current.NumParents() > 1 || current.NumParents() == 0 {
			return commits, nil
		}
		commits = append(commits, current)
		current, err = GetFirstParent(current)
		if err != nil {
			return nil, err
		}
	}
}

func isCommitOnComponent(commit *object.Commit, component string) (bool, error) {
	// Check if this commit impacts component
	stats, err := commit.Stats()
	if err != nil {
		return false, err
	}
	for _, stat := range stats {
		if strings.HasPrefix(stat.Name, "components/"+component) {
			return true, nil
		}
	}
	return false, nil
	// componentTree, err := tree.Tree("components/" + component)
	// if err != nil {
	// 	return false, err
	// }
	// return len(componentTree.Entries) > 0, nil
}

func deduplicate(commits []*object.Commit) []*object.Commit {
	seen := map[string]bool{}
	deduplicated := []*object.Commit{}
	for _, c := range commits {
		if _, ok := seen[c.Hash.String()]; !ok {
			seen[c.Hash.String()] = true
			deduplicated = append(deduplicated, c)
		}
	}
	return deduplicated
}
