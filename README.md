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
│   ├── security.md
│   └── style.md
├── Dockerfile              # Multi-stage build (Distroless)
├── Jenkinsfile.example     # How to run it in Jenkins
└── Jenkinsfile             # Pipeline definition
```
