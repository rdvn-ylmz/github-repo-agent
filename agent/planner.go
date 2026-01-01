package agent

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github-repo-agent/github"
	"github-repo-agent/llm"
)

type Planner struct {
	client *llm.Client
	repo   *github.Repo
	goal   string
}

func NewPlanner(client *llm.Client, repo *github.Repo, goal string) *Planner {
	return &Planner{
		client: client,
		repo:   repo,
		goal:   goal,
	}
}

func (p *Planner) CreatePlan() (*Plan, error) {
	// Load prompt template
	promptContent, err := os.ReadFile("prompts/planner.txt")
	if err != nil {
		return nil, fmt.Errorf("failed to load prompt: %w", err)
	}

	// Read repository structure
	repoInfo, err := p.analyzeRepo()
	if err != nil {
		return nil, fmt.Errorf("failed to analyze repository: %w", err)
	}

	// Build prompt
	prompt := strings.ReplaceAll(string(promptContent), "{{GOAL}}", p.goal)
	prompt = strings.ReplaceAll(prompt, "{{REPO_INFO}}", repoInfo)

	// Query LLM
	response, err := p.client.Query(prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to get plan from LLM: %w", err)
	}

	// Parse response into plan
	plan := p.parsePlan(response)
	return plan, nil
}

func (p *Planner) analyzeRepo() (string, error) {
	var info strings.Builder

	info.WriteString(fmt.Sprintf("Repository URL: %s\n", p.repo.URL))
	info.WriteString(fmt.Sprintf("Local Path: %s\n", p.repo.LocalPath))

	// List directories
	dirs, err := os.ReadDir(p.repo.LocalPath)
	if err != nil {
		return "", err
	}

	info.WriteString("\nDirectories:\n")
	for _, d := range dirs {
		if d.IsDir() {
			info.WriteString(fmt.Sprintf("  - %s/\n", d.Name()))
		} else {
			info.WriteString(fmt.Sprintf("  - %s\n", d.Name()))
		}
	}

	// Check for go.mod
	goModPath := filepath.Join(p.repo.LocalPath, "go.mod")
	if _, err := os.Stat(goModPath); err == nil {
		info.WriteString("\nHas Go module: Yes\n")
	} else {
		info.WriteString("\nHas Go module: No\n")
	}

	return info.String(), nil
}

func (p *Planner) parsePlan(response string) *Plan {
	plan := &Plan{
		Files: make(map[string]string),
		Steps: make([]Step, 0),
	}

	lines := strings.Split(response, "\n")
	var currentFile string
	var currentContent strings.Builder

	for _, line := range lines {
		if strings.HasPrefix(line, "STEP:") {
			if currentFile != "" {
				plan.Files[currentFile] = strings.TrimSpace(currentContent.String())
				currentContent.Reset()
			}

			stepDesc := strings.TrimPrefix(line, "STEP:")
			plan.Steps = append(plan.Steps, Step{
				Description: strings.TrimSpace(stepDesc),
				Type:        "create",
			})
		} else if strings.HasPrefix(line, "FILE:") {
			if currentFile != "" {
				plan.Files[currentFile] = strings.TrimSpace(currentContent.String())
				currentContent.Reset()
			}
			currentFile = strings.TrimPrefix(line, "FILE:")
			currentFile = strings.TrimSpace(currentFile)
		} else if strings.HasPrefix(line, "MODULE_INIT:") {
			value := strings.TrimPrefix(line, "MODULE_INIT:")
			if strings.TrimSpace(value) == "true" {
				plan.ModuleInit = true
			}
		} else {
			currentContent.WriteString(line + "\n")
		}
	}

	if currentFile != "" {
		plan.Files[currentFile] = strings.TrimSpace(currentContent.String())
	}

	return plan
}
