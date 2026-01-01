package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github-repo-agent/agent"
	"github-repo-agent/github"
)

func main() {
	repoURL := flag.String("repo", "", "GitHub repository URL")
	goal := flag.String("goal", "", "High-level objective")
	flag.Parse()

	if *repoURL == "" || *goal == "" {
		fmt.Println("Error: --repo and --goal are required")
		flag.Usage()
		os.Exit(1)
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("Error: GITHUB_TOKEN environment variable is required")
	}

	repo, err := github.NewRepo(*repoURL, token)
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}

	if err := repo.Clone(); err != nil {
		log.Fatalf("Failed to clone repository: %v", err)
	}

	a := agent.New(repo, *goal)
	if err := a.Execute(); err != nil {
		log.Fatalf("Agent execution failed: %v", err)
	}

	if err := repo.CommitAndPush(); err != nil {
		log.Fatalf("Failed to commit and push: %v", err)
	}

	fmt.Println("Autonomous update completed successfully")
}
