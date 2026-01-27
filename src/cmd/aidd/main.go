package main

import (
	"errors"
	"log"

	"github.com/tomoyuki65/go-aidd/internal/config"
	"github.com/tomoyuki65/go-aidd/internal/provider/container"
	"github.com/tomoyuki65/go-aidd/internal/provider/github"
)

// Generate "task.md" from task information
func generateTaskMd(provider, repository, label string) error {
	// Switch processing by provider
	switch provider {
	case "GitHub":
		return github.GenerateTaskMd(repository, label)
	case "container":
		return container.ExecContainer()
	default:
		return errors.New("unsupported provider is set")
	}
}

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Generate "task.md" from task info
	if err := generateTaskMd(cfg.Issue.Provider, cfg.GitHub.Repository, cfg.Issue.Label); err != nil {
		log.Fatalf("failed to generate file 'task.md': %v", err)
	}
}
