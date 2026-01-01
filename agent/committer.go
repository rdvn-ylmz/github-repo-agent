package agent

import (
	"fmt"
	"log"
	"os/exec"

	"github-repo-agent/github"
)

type Committer struct {
	repo *github.Repo
}

func NewCommitter(repo *github.Repo) *Committer {
	return &Committer{repo: repo}
}

func (c *Committer) Commit(message string) error {
	// git add .
	cmd := c.repo.Command("git", "add", ".")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git add failed: %w, output: %s", err, string(output))
	}

	// git commit -m message
	cmd = c.repo.Command("git", "commit", "-m", message)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git commit failed: %w, output: %s", err, string(output))
	}

	log.Printf("Committed changes: %s", message)
	return nil
}

func (c *Committer) Push() error {
	// git push
	cmd := exec.Command("git", "-C", c.repo.LocalPath, "push")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git push failed: %w, output: %s", err, string(output))
	}

	log.Println("Pushed changes to remote repository")
	return nil
}

func (c *Committer) CommitAndPush(message string) error {
	if err := c.Commit(message); err != nil {
		return err
	}
	return c.Push()
}
