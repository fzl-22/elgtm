package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fzl-22/elgtm/internal/config"
	"github.com/fzl-22/elgtm/internal/logger"
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

	engine, err := Initialize(ctx, cfg)
	if err != nil {
		slog.Error("Initialization failed", "error", err)
		os.Exit(1)
	}

	if err := engine.Run(ctx); err != nil {
		slog.Error("Review failed", "error", err)
		os.Exit(1)
	}

	slog.Info("Review completed successfully")
}
