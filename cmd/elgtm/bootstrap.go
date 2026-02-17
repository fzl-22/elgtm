package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/fzl-22/elgtm/internal/config"
	"github.com/fzl-22/elgtm/internal/llm"
	"github.com/fzl-22/elgtm/internal/reviewer"
	"github.com/fzl-22/elgtm/internal/scm"
)

func Initialize(ctx context.Context, cfg *config.Config) (*reviewer.Engine, error) {
	timeoutDuration := time.Duration(cfg.System.Timeout) * time.Second
	httpClient := http.Client{Timeout: timeoutDuration}

	var err error

	// Initialize SCM
	var scmClient scm.SCMClient

	switch cfg.SCM.Platform {
	case config.PlatformGitHub:
		scmClient, err = scm.NewGitHubClient(&httpClient, cfg.SCM)
	default:
		return nil, fmt.Errorf("unsupported platform: %s", cfg.SCM.Platform)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to initialize SCM client: %w", err)
	}

	// Initialize LLM Driver
	var llmDriver llm.Driver
	switch cfg.LLM.Provider {
	case config.ProviderGemini:
		llmDriver, err = llm.NewGeminiDriver(ctx, cfg.LLM.APIKey)
	default:
		return nil, fmt.Errorf("unsupported LLM provider: %s", cfg.LLM.Provider)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to initialize LLM driver: %w", err)
	}

	llmClient := llm.NewClient(llmDriver, cfg.LLM)

	return reviewer.NewEngine(*cfg, scmClient, llmClient), nil
}
