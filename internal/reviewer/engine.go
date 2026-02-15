package reviewer

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/fzl-22/elgtm/internal/config"
	"github.com/fzl-22/elgtm/internal/llm"
	"github.com/fzl-22/elgtm/internal/scm"
)

type Engine struct {
	cfg       config.Config
	scmClient scm.SCMClient
	llmClient llm.LLMClient
}

func NewEngine(cfg config.Config, scmClient scm.SCMClient, llmClient llm.LLMClient) *Engine {
	return &Engine{
		cfg:       cfg,
		scmClient: scmClient,
		llmClient: llmClient,
	}
}

func (e *Engine) Run(ctx context.Context) error {
	promptPath, err := e.resolvePromptPath(e.cfg.Review.PromptDir, e.cfg.Review.PromptType)
	if err != nil {
		return fmt.Errorf("prompt resolution failed: %w", err)
	}

	promptContent, err := os.ReadFile(promptPath)
	if err != nil {
		return fmt.Errorf("failed to load prompt file [%s]: %w", promptPath, err)
	}

	if e.cfg.System.LogLevel == "debug" {
		slog.Debug("Loaded Prompt", "path", promptPath, "content_length", len(promptContent))
	}

	pr, err := e.scmClient.GetPullRequest(ctx, e.cfg.SCM.Owner, e.cfg.SCM.Repo, e.cfg.SCM.PRNumber)
	if err != nil {
		return fmt.Errorf("failed to get pull request: %w", err)
	}

	slog.Info("PR Fetched", "pr_number", pr.Number, "title", pr.Title, "author", pr.Author, "diff_size", len(pr.RawDiff))

	issueComment, err := e.llmClient.GenerateIssueComment(ctx, *pr)

	slog.Info("Posting comment", "repo", e.cfg.SCM.Repo, "pr", e.cfg.SCM.PRNumber)

	err = e.scmClient.PostIssueComment(ctx, e.cfg.SCM.Owner, e.cfg.SCM.Repo, e.cfg.SCM.PRNumber, &scm.IssueComment{
		Body: issueComment.Body,
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
