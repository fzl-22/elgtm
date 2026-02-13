package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fzl-22/elgtm/internal/config"
	"github.com/fzl-22/elgtm/internal/reviewer"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		slog.Error("Config load failed", "error", err)
		os.Exit(1)
	}

	setupLogger(cfg.System.LogLevel)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	timeoutDuration := time.Duration(cfg.System.Timeout) * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeoutDuration)
	defer cancel()

	slog.Info("Starting ELGTM",
		"platform", cfg.SCM.Platform,
		"provider", cfg.LLM.Provider,
		"timeout", timeoutDuration.String(),
	)

	engine := reviewer.NewEngine()
	if err := engine.Run(ctx, *cfg); err != nil {
		slog.Error("Review failed", "error", err)
		os.Exit(1)
	}

	slog.Info("Review completed successfully")
}

func setupLogger(levelStr string) {
	var level slog.Level
	switch levelStr {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{Level: level}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	slog.SetDefault(logger)
}
