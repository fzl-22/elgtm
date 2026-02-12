# ELGTM - Enhanced LGTM

## Expected Code Structure

```plaintext
elgtm/
├── cmd/
│   └── elgtm/
│       └── main.go         # Entry point: Wires Config -> SCM -> LLM -> Engine
├── internal/
│   ├── scm/                # Package: scm (Source Control Management)
│   │   ├── scm.go          # Defines 'type Client interface'
│   │   ├── github.go       # Implements GitHub logic
│   │   └── gitlab.go       # (Future) Implements GitLab logic
│   ├── llm/                # Package: llm (Large Language Model)
│   │   ├── llm.go          # Defines 'type Client interface'
│   │   ├── gemini.go       # Implements Gemini logic
│   │   └── claude.go       # (Future) Implements Claude logic
│   ├── config/             # Package: config
│   │   └── config.go       # Parses env vars (SCM_*, LLM_*, etc)
│   └── reviewer/           # Package: reviewer (Core Logic)
│       └── engine.go       # Orchestrator: Reads .reviewer/ -> Fetches SCM -> Asks LLM
├── .reviewer/              # User Config: Custom Prompts
│   ├── general.md
│   └── security.md
├── .github/
│   └── workflows/
│       └── publish.yml     # Automated Docker Hub publishing
├── Dockerfile              # Multi-stage build (Scratch)
└── Jenkinsfile             # Pipeline definition (Alternative CI)
```
