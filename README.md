# github-repo-agent

An autonomous CLI tool that works on GitHub repositories using local LLM models.

## What It Does

`github-repo-agent` autonomously clones a GitHub repository, analyzes its structure, generates a plan based on a high-level goal, creates or modifies Go source code, and commits/pushes the changes back to the repository.

All LLM interactions use local models (via Ollama) - no paid APIs, no external services, and no secrets are stored or logged.

## Installation

```bash
go install github.com/rdvn-ylmz/github-repo-agent@latest
```

## Usage

```bash
export GITHUB_TOKEN=your_token_here
export LLM_BASE_URL=http://localhost:11434/api/chat
export LLM_MODEL=llama3.2

github-repo-agent \
  --repo https://github.com/owner/repo \
  --goal "Add a REST API endpoint for user authentication"
```

## Required Environment Variables

- `GITHUB_TOKEN` - GitHub personal access token for authentication
- `LLM_BASE_URL` - URL of local LLM API (default: http://localhost:11434/api/chat)
- `LLM_MODEL` - Model name to use (default: llama3.2)

## How It Works

1. **Clone**: Clones the target GitHub repository to a temporary location
2. **Analyze**: Examines repository structure and existing code
3. **Plan**: Uses local LLM to create a detailed implementation plan
4. **Build**: Generates complete, production-ready Go source files
5. **Commit**: Runs `git add .` and `git commit -m "autonomous update"`
6. **Push**: Pushes changes back to the remote repository

## Security

- GitHub token is ONLY read from the `GITHUB_TOKEN` environment variable
- Token is NEVER stored in files, logs, or any other location
- Token is NEVER printed to stdout or stderr
- All LLM queries use local models only (Ollama/LM Studio/etc.)

## Local LLM Setup

Install Ollama:
```bash
curl https://ollama.ai/install.sh | sh
```

Pull a model:
```bash
ollama pull llama3.2
```

The LLM API will be available at `http://localhost:11434/api/chat`

## Project Structure

```
github-repo-agent/
├── main.go          # CLI entry point
├── agent/           # Core agent logic
│   ├── planner.go   # Plan generation
│   ├── builder.go   # Code generation
│   ├── committer.go # Git operations
│   └── types.go     # Shared types
├── github/          # GitHub repository handling
│   └── repo.go
├── llm/             # Local LLM client
│   └── client.go
└── prompts/         # Prompt templates
    ├── planner.txt
    └── builder.txt
```

## License

MIT
