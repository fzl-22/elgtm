package reviewer

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fzl-22/elgtm/internal/config"
	"github.com/fzl-22/elgtm/internal/scm"
)

type Engine struct{}

func NewEngine() *Engine {
	return &Engine{}
}

func (e *Engine) Run(ctx context.Context, cfg config.Config, scmClient scm.SCMClient) error {
	promptPath, err := e.resolvePromptPath(cfg.Review.PromptDir, cfg.Review.PromptType)
	if err != nil {
		return fmt.Errorf("prompt resolution failed: %w", err)
	}

	promptContent, err := os.ReadFile(promptPath)
	if err != nil {
		return fmt.Errorf("failed to load prompt file [%s]: %w", promptPath, err)
	}

	if cfg.System.LogLevel == "debug" {
		fmt.Printf("\n--- [DEBUG] LOADED PROMPT (%s) ---\n%s\n------------------------------------\n", promptPath, string(promptContent))
	}

	pr, err := scmClient.GetPullRequest(ctx, cfg.SCM.Owner, cfg.SCM.Repo, cfg.SCM.PRNumber)
	if err != nil {
		return fmt.Errorf("failed to get pull request: %w", err)
	}

	fmt.Println(pr)

	commentBody := "Hi, It is a test comment!"
	err = scmClient.PostIssueComment(ctx, cfg.SCM.Owner, cfg.SCM.Repo, cfg.SCM.PRNumber, &scm.IssueComment{
		Body: &commentBody,
	})

	return err
}

func (e *Engine) resolvePromptPath(userDir, promptType string) (string, error) {
	filename := fmt.Sprintf("%s.md", promptType)

	// PRIORITY 1: User's local configuration
	userPath := filepath.Join(userDir, filename)
	if _, err := os.Stat(userPath); err == nil {
		return userPath, nil
	}

	// PRIORITY 2: System Defaults (/etc/elgtm/defaults/general.md)
	systemDir := os.Getenv("PROMPT_DEFAULTS")
	if systemDir != "" {
		systemPath := filepath.Join(systemDir, filename)
		if _, err := os.Stat(systemPath); err == nil {
			return systemPath, nil
		}
	}

	// ERROR: Not found anywhere
	return "", fmt.Errorf("prompt '%s' not found in local [%s] or system [%s]", filename, userPath, systemDir)
}
