package agent

type (
	Config struct {
		RepoURL string
		Goal    string
		Token   string
	}

	Plan struct {
		Steps      []Step
		Files      map[string]string
		ModuleInit bool
	}

	Step struct {
		Description string
		Type        string
		Params      map[string]string
	}

	BuildResult struct {
		Files     map[string]string
		Commands  []string
		CommitMsg string
	}
)
