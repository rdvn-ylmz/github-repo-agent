package github

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Repo struct {
	URL       string
	Token     string
	LocalPath string
	Owner     string
	Name      string
}

func NewRepo(url, token string) (*Repo, error) {
	parts := strings.Split(url, "/")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid repository URL")
	}

	owner := parts[len(parts)-2]
	name := strings.TrimSuffix(parts[len(parts)-1], ".git")

	localPath := filepath.Join(os.TempDir(), name)

	return &Repo{
		URL:       url,
		Token:     token,
		LocalPath: localPath,
		Owner:     owner,
		Name:      name,
	}, nil
}

func (r *Repo) Clone() error {
	// Remove existing directory if it exists
	if _, err := os.Stat(r.LocalPath); err == nil {
		if err := os.RemoveAll(r.LocalPath); err != nil {
			return fmt.Errorf("failed to remove existing directory: %w", err)
		}
	}

	// Clone with auth
	authURL := r.insertAuth()
	cmd := exec.Command("git", "clone", authURL, r.LocalPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git clone failed: %w, output: %s", err, string(output))
	}

	return nil
}

func (r *Repo) Command(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	cmd.Dir = r.LocalPath
	return cmd
}

func (r *Repo) CommitAndPush() error {
	// git add .
	cmd := r.Command("git", "add", ".")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git add failed: %w, output: %s", err, string(output))
	}

	// git commit -m "autonomous update"
	cmd = r.Command("git", "commit", "-m", "autonomous update")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git commit failed: %w, output: %s", err, string(output))
	}

	// git push
	cmd = exec.Command("git", "-C", r.LocalPath, "push")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git push failed: %w, output: %s", err, string(output))
	}

	return nil
}

func (r *Repo) insertAuth() string {
	if r.Token == "" {
		return r.URL
	}

	// Insert token into URL for authentication
	// https://github.com/user/repo.git -> https://TOKEN@github.com/user/repo.git
	return strings.Replace(r.URL, "https://", fmt.Sprintf("https://%s@", r.Token), 1)
}
