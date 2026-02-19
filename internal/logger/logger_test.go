package logger_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/fzl-22/elgtm/internal/logger"
	"github.com/stretchr/testify/assert"
)

func TestLogger_Setup(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected slog.Level
	}{
		{"Success_DebugLevel", "debug", slog.LevelDebug},
		{"Success_WarnLevel", "warn", slog.LevelWarn},
		{"Success_ErrorLevel", "error", slog.LevelError},
		{"Success_InfoLevel", "info", slog.LevelInfo},
		{"Success_DefaultFallback", "other", slog.LevelInfo}, // default fallback
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger.Setup(tt.input)

			l := slog.Default()
			ctx := context.Background()

			// The expected level MUST be enabled
			assert.True(t, l.Enabled(ctx, tt.expected), "Expected level %v to be enabled", tt.expected)

			// The level below MUST be disabled (to ensure it's not set too low)
			if tt.expected > slog.LevelDebug {
				oneLevelLower := tt.expected - 1
				assert.False(t, l.Enabled(ctx, oneLevelLower), "Expected level %v to be disabled", oneLevelLower)
			}
		})
	}
}
