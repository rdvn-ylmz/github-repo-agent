package agent

import (
	"log"

	"github-repo-agent/github"
	"github-repo-agent/llm"
)

type Agent struct {
	repo    *github.Repo
	goal    string
	planner *Planner
	builder *Builder
	client  *llm.Client
}

func New(repo *github.Repo, goal string) *Agent {
	client := llm.NewClient()
	return &Agent{
		repo:    repo,
		goal:    goal,
		client:  client,
		planner: NewPlanner(client, repo, goal),
		builder: NewBuilder(client, repo, goal),
	}
}

func (a *Agent) Execute() error {
	log.Printf("Creating execution plan for goal: %s", a.goal)

	plan, err := a.planner.CreatePlan()
	if err != nil {
		return err
	}

	log.Printf("Plan created with %d steps and %d files", len(plan.Steps), len(plan.Files))

	for i, step := range plan.Steps {
		log.Printf("Step %d/%d: %s", i+1, len(plan.Steps), step.Description)
	}

	result, err := a.builder.BuildFromPlan(plan)
	if err != nil {
		return err
	}

	log.Printf("Built %d files", len(result.Files))

	return nil
}
