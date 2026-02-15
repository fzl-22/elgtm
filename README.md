# ELGTM (Enhanced LGTM)

**ELGTM** is an automated code review tool that connects your Source Control Management (SCM) with Large Language Models (LLMs) to provide intelligent, context-aware feedback on Pull Requests.

It doesn't just say "LGTM" — it analyzes your diffs for bugs, maintainability issues, and security risks based on custom prompt templates you define in your repo.

## Features

- **AI-Powered Reviews**: Currently supports **Google Gemini**.
- **Fully Configurable Prompts**: Define your own review personas (e.g., "Security Auditor", "Nitpicker", "Senior Engineer") using simple Markdown templates in your repo.
- **CI/CD Native**: Runs effortlessly in GitHub Actions. Support for Jenkins is coming soon.
- **Zero-Dependency Binary**: Built as a static Go binary on Alpine Linux for speed and security.

---

## Quick Start

### GitHub Actions

The easiest way to use ELGTM is via our composite GitHub Action. Add this to your `.github/workflows/review.yml`:

```yaml
name: AI Code Review

on:
  pull_request:
    types: [opened, synchronize]

permissions:
  pull-requests: write # Required to post comments
  contents: read

jobs:
  review:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout Code
        uses: actions/checkout@v4

      - name: Run ELGTM
        uses: fzl-22/elgtm/action@v1
        with:
          llm_provider: "gemini"
          llm_model: "gemini-2.5-flash"
          llm_api_key: ${{ secrets.GEMINI_API_KEY }}
          # github_token is automatically picked up
```

## Configuration

| Variable            | Description                                             | Default                                 |
| ------------------- | ------------------------------------------------------- | --------------------------------------- |
| **LLM Settings**    |                                                         |                                         |
| LLM_PROVIDER        | LLM Provider (e.g., `gemini`)                           | Required                                |
| LLM_MODEL           | Model ID (e.g., `gemini-2.5-flash`)                     | Required                                |
| LLM_API_KEY         | Your AI provider's API Key                              | Required                                |
| LLM_TEMPERATURE     | Creativity (0.0 - 1.0)                                  | `0.2`                                   |
| LLM_MAX_TOKENS      | Max output tokens for the review                        | `4096`                                  |
| **SCM Settings**    |                                                         |
| SCM_PLATFORM        | Source control platform (github)                        | `github` in GitHub Actions              |
| SCM_TOKEN           | Access token (`PAT` or `GITHUB_TOKEN`)                  | `${{ github.token }}` in GitHub Actions |
| SCM_OWNER           | Repo owner                                              | Auto in GitHub Actions                  |
| SCM_REPO            | Repo name                                               | Auto in GitHub Actions                  |
| SCM_PR_NUMBER       | The PR number to review                                 | Auto in GitHub Actions                  |
| SCM_MAX_DIFF_SIZE   | Max characters of diff to process                       | `2097152`                               |
| **Review Settings** |                                                         |
| REVIEW_PROMPT_DIR   | Prompt directory (e.g. `.reviewer`)                     | `.reviewer`                             |
| REVIEW_PROMPT_TYPE  | Prompt filename at `REVIEW_PROMPT_DIR` (e.g. `general`) | `general`                               |

## Customizing Prompts

ELGTM allows you to define custom personas and review criteria by creating Markdown templates. This lets you switch between different "modes" (e.g., a "Security Auditor", a "Nitpicker", or a "Senior Architect") simply by changing a configuration variable.

### 1. Create the Prompt Directory

Create a directory in the root of your repository to store your templates. The default is `.reviewer`.

```bash
mkdir .reviewer
```

### 2. Create a Template File

Create a Markdown file inside that directory. The filename becomes the `prompt_type` you will use later.

Example: `.reviewer/security.md`:

````markdown
# Role

You are a Principal Security Engineer at a Fintech company.

# Task

Review the following Pull Request code changes for security vulnerabilities.
Be extremely strict.

# Focus Areas

1. **SQL Injection**: Ensure all database queries use parameterized statements.
2. **PII Leaks**: Check for hardcoded credentials, API keys, or logging of sensitive user data.
3. **Authorization**: Verify that the caller has permission to execute the action.

# Context

**Title**: {{ .Title }}
**Author**: {{ .Author }}

**Description**:
{{ .Body }}

**Code Changes**:

```text
{{ .RawDiff }}
```
````

### 3. Available Variables

You can use the following Go template variables in your Markdown files to inject PR context:

| Variable         | Description                                              |
| :--------------- | :------------------------------------------------------- |
| `{{ .Title }}`   | The title of the Pull Request                            |
| `{{ .Body }}`    | The description/body of the Pull Request                 |
| `{{ .Author }}`  | The username of the PR author                            |
| `{{ .RawDiff }}` | The raw git diff of the changes (truncated if too large) |
| `{{ .Number }}`  | The Pull Request number                                  |
| `{{ .URL }}`     | The URL of the Pull Request                              |

### 4. Activate the Prompt

To use your new prompt, set the `REVIEW_PROMPT_TYPE` environment variable (or `prompt_type` input in GitHub Actions) to the filename **without the extension**.

**GitHub Actions Example:**

```yaml
- uses: fzl-22/elgtm/action@v1
  with:
    # ... other inputs ...
    prompt_type: "security" # Uses .reviewer/security.md
```

## Contributing

Contributions are welcome! If you want to add support for GitLab, OpenAI, or Claude, feel free to open a PR.

1. Fork the repository.
2. Create a feature branch.
3. Commit your changes.
4. Open a Pull Request.

## License

This project is open-source and available under the **MIT License**.
