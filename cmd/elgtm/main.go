package main

import (
	"log/slog"
	"os"

	"github.com/fzl-22/elgtm/internal/config"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		slog.Error("Config load failed", "error", err)
		os.Exit(1)
	}

	var level slog.Level
	switch cfg.System.LogLevel {
	case "debug ":
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

	slog.Info("Configuration loaded successfully", "scm platform", cfg.SCM.Platform, "llm provider", cfg.LLM.Provider)
}
