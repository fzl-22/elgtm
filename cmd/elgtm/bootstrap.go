package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/fzl-22/elgtm/internal/config"
	"github.com/fzl-22/elgtm/internal/llm"
	"github.com/fzl-22/elgtm/internal/reviewer"
	"github.com/fzl-22/elgtm/internal/scm"
)

func Initialize(ctx context.Context, cfg *config.Config) (*reviewer.Engine, error) {
	timeoutDuration := time.Duration(cfg.System.Timeout) * time.Second
	httpClient := http.Client{Timeout: timeoutDuration}

	var scmClient scm.SCMClient
	var scmErr error
	if cfg.SCM.Platform == "github" {
		scmClient = scm.NewGitHubClient(&httpClient, cfg.SCM)
	} else {
		scmErr = fmt.Errorf("unsupported platform: %s", cfg.SCM.Platform)
	}

	if scmErr != nil {
		slog.Error("SCM initialization failed", "error", scmErr)
		os.Exit(1)
	}

	var llmClient llm.LLMClient
	var llmErr error
	if cfg.LLM.Provider == "gemini" {
		llmClient, llmErr = llm.NewGeminiClient(ctx, cfg.LLM)
	} else {
		llmErr = fmt.Errorf("unsupported LLM provider: %s", cfg.LLM.Provider)
	}

	if llmErr != nil {
		slog.Error("LLM initialization failed", "error", llmErr)
		os.Exit(1)
	}

	return reviewer.NewEngine(*cfg, scmClient, llmClient), nil
}
