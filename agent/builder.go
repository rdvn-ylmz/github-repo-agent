package agent

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github-repo-agent/github"
	"github-repo-agent/llm"
)

type Builder struct {
	client *llm.Client
	repo   *github.Repo
	goal   string
}

func NewBuilder(client *llm.Client, repo *github.Repo, goal string) *Builder {
	return &Builder{
		client: client,
		repo:   repo,
		goal:   goal,
	}
}

func (b *Builder) BuildFromPlan(plan *Plan) (*BuildResult, error) {
	result := &BuildResult{
		Files: make(map[string]string),
	}

	// Initialize Go module if needed
	if plan.ModuleInit {
		if err := b.initGoModule(); err != nil {
			return nil, fmt.Errorf("failed to initialize Go module: %w", err)
		}
		result.Files["go.mod"] = "module initialized"
	}

	// Create all files
	for path, content := range plan.Files {
		fullPath := filepath.Join(b.repo.LocalPath, path)
		if err := b.createFile(fullPath, content); err != nil {
			return nil, fmt.Errorf("failed to create file %s: %w", path, err)
		}
		result.Files[path] = content
	}

	// Generate commit message
	result.CommitMsg = b.generateCommitMsg(plan)

	return result, nil
}

func (b *Builder) initGoModule() error {
	moduleName := b.extractModuleName()
	cmd := b.repo.Command("go", "mod", "init", moduleName)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func (b *Builder) createFile(path, content string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return err
	}

	return nil
}

func (b *Builder) extractModuleName() string {
	parts := strings.Split(b.repo.URL, "/")
	if len(parts) >= 2 {
		return parts[len(parts)-2] + "/" + parts[len(parts)-1]
	}
	return "github-repo-agent"
}

func (b *Builder) generateCommitMsg(plan *Plan) string {
	if len(plan.Steps) > 0 {
		return "autonomous update: " + plan.Steps[0].Description
	}
	return "autonomous update"
}
