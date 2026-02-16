package logger_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/fzl-22/elgtm/internal/logger"
	"github.com/stretchr/testify/assert"
)

func TestSetup(t *testing.T) {
	tests := []struct {
		input    string
		expected slog.Level
	}{
		{"debug", slog.LevelDebug},
		{"warn", slog.LevelWarn},
		{"error", slog.LevelError},
		{"info", slog.LevelInfo},
		{"other", slog.LevelInfo}, // default fallback
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
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
