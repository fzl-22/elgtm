package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fzl-22/elgtm/internal/config"
	"github.com/fzl-22/elgtm/internal/logger"
	"github.com/fzl-22/elgtm/internal/reviewer"
	"github.com/fzl-22/elgtm/internal/scm"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		slog.Error("Config load failed", "error", err)
		os.Exit(1)
	}

	logger.Setup(cfg.System.LogLevel)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	timeoutDuration := time.Duration(cfg.System.Timeout) * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeoutDuration)
	defer cancel()

	slog.Info("Starting ELGTM",
		"scm_platform", cfg.SCM.Platform,
		"llm_provider", cfg.LLM.Provider,
		"system_log_level", cfg.System.LogLevel,
		"system_timeout", timeoutDuration.String(),
	)

	httpClient := http.Client{Timeout: timeoutDuration}

	var scmClient scm.SCM
	if cfg.SCM.Platform == "github" {
		scmClient = scm.NewGitHubClient(&httpClient, cfg.SCM.Token)
	} else {
		slog.Error("Unsupported SCM platform", "error", err)
	}

	engine := reviewer.NewEngine()
	if err := engine.Run(ctx, *cfg, scmClient); err != nil {
		slog.Error("Review failed", "error", err)
		os.Exit(1)
	}

	slog.Info("Review completed successfully")
}
