# CodeReviewr

## Expected Code Structure

```plaintext
codereviewr/
├── cmd/
│   └── codereviewr/        # Main entry point
│       └── main.go
├── internal/
│   ├── platform/           # Git Provider Facade (GitHub, for now)
│   │   ├── provider.go     # Interface definition
│   │   └── github.go
│   ├── ai/                 # AI Provider Facade (Gemini, for now)
│   │   ├── client.go       # Interface definition
│   │   ├── gemini.go
│   │   └── claude.go
│   ├── config/             # Env var and Prompt parsing
│   │   └── config.go
│   └── reviewer/           # Core logic (orchestrator)
│       └── engine.go
├── pkg/                    # Publicly shareable code (optional)
├── .reviewr/               # Default user prompts
│   ├── security.md
│   └── style.md
├── Dockerfile
├── Jenkinsfile.example     # How to run it in Jenkins
└── Jenkinsfile
```
