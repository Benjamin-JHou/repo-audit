# ctxqa - AI-Powered Folder Analyzer

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Version](https://img.shields.io/badge/Go-1.22%2B-00ADD8)](https://golang.org/dl/)
[![Platform](https://img.shields.io/badge/platform-macOS%20%7C%20Linux%20%7C%20Windows-lightgrey)](https://github.com/Benjamin-JHou/repo-audit)

A terminal CLI tool that analyzes any folder using AI. Traverse the entire directory tree, send file contents to an LLM in batches, and get a structured report covering file structure, organization suggestions, progress analysis, issue detection, quality scoring, and improvement recommendations.

## Why ctxqa?

- **Any folder, not just code repos**: Works on documents, configs, data files, mixed content
- **Intelligent file classification**: Automatically categorizes files into code, config, data, document, web, and other types
- **Structured analysis**: Not just chat - get actionable reports with severity-classified issues and quality scores
- **Multiple LLM backends**: OpenAI GPT-4o, Claude Sonnet, or any compatible API
- **Privacy first**: Your files are sent to your API key, not to us
- **One command install**: No complex setup, works in seconds

## Features

- Full directory traversal with smart exclusions
- Batch processing for large folders
- File type classification and statistics
- File structure and organization suggestions
- Progress analysis based on file content and commit history
- Severity-classified issue detection (critical, error, warning, info)
- Quality scoring (naming, error handling, duplication, security)
- Prioritized improvement suggestions (high, medium, low)
- Multiple output formats: text, JSON, Markdown
- Incremental analysis for changed files only
- Resume support for interrupted analyses
- Context enrichment with Git branch and commit history

## Installation

```bash
curl -sSL https://raw.githubusercontent.com/Benjamin-JHou/repo-audit/main/scripts/install.sh | sh
```

**Supported platforms**: macOS (Intel/Apple Silicon), Linux (amd64/arm64), Windows (WSL)

## Quick Start

```bash
# Configure your API key (first time only)
ctxqa config init

# Analyze current folder
ctxqa analyze

# Analyze a specific folder
ctxqa analyze --dir /path/to/folder

# Preview files without analysis
ctxqa analyze --dry-run

# Incremental analysis (only changed files)
ctxqa analyze --incremental

# Export report as Markdown
ctxqa analyze --format markdown -o report.md

# Show only warnings and above
ctxqa analyze --severity warning
```

## Usage

### Analyze a Folder

```bash
ctxqa analyze [flags]
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--format` | `-f` | `text` | Output format: text, json, markdown |
| `-o`, `--output` | | stdout | Output file path |
| `--severity` | | `all` | Minimum severity: all, info, warning, error, critical |
| `--section` | | `all` | Specific section: overview, structure, progress, issues, suggestions |
| `--incremental` | | `false` | Only analyze changed files since last run |
| `--resume` | | `false` | Resume from last interrupted analysis |
| `--no-context` | | `false` | Skip Git/context collection |
| `--dir` | `-d` | `.` | Folder to analyze |
| `--concurrency` | `-j` | `1` | Parallel batch processors |
| `--dry-run` | | `false` | List files only, do not analyze |
| `--summary` | | `false` | Brief summary only |

### Configuration

```bash
# Interactive setup
ctxqa config init

# Show current configuration
ctxqa config show

# Switch API provider
ctxqa config set-provider anthropic

# Set model
ctxqa config set-model claude-sonnet-4-20250514

# Set API key
ctxqa config set-api-key

# Add exclusion patterns
ctxqa config exclude-add "**/*.test.ts" "**/__snapshots__/**"

# Remove exclusion pattern
ctxqa config exclude-remove "**/*.test.ts"

# Reset to defaults
ctxqa config reset
```

## Report Structure

The analysis report includes:

### 1. Folder Overview
- Total files and total size
- File type breakdown (code, config, data, document, web, other)
- Top-level directory statistics
- Detected tech stack
- Git branch and recent commits

### 2. File Structure and Organization Suggestions
- Current layout summary
- Proposed reorganization plan
- File grouping recommendations
- Naming convention issues
- Duplicate file detection

### 3. Progress Analysis
- Completed work (based on file content and commit history)
- Work in progress
- Pending tasks

### 4. Issues Found
Classified by severity:
- **Critical**: Blocking issues that need immediate attention
- **Error**: Serious problems affecting functionality
- **Warning**: Potential issues that should be reviewed
- **Info**: Suggestions and minor observations

### 5. Quality Scores
Each scored 0-10:
- Naming conventions
- Error handling
- Code duplication
- Security
- Overall score

### 6. Improvement Suggestions
Prioritized recommendations:
- **High**: Address immediately
- **Medium**: Plan for next iteration
- **Low**: Nice to have improvements

## Supported API Providers

### OpenAI Compatible

```bash
ctxqa config set-provider openai
ctxqa config set-model gpt-4o
```

Default URL: `https://api.openai.com/v1`

Supported models: `gpt-4o`, `gpt-4o-mini`, `gpt-4-turbo`

### Anthropic Compatible

```bash
ctxqa config set-provider anthropic
ctxqa config set-model claude-sonnet-4-20250514
```

Default URL: `https://api.anthropic.com/v1`

Supported models: `claude-sonnet-4-20250514`, `claude-opus-4-0-20250620`

## How It Works

```
1. You run `ctxqa analyze` in any folder
2. ctxqa scans all files recursively, skipping binaries and excluded paths
3. Files are classified by type (code, config, data, document, web, other)
4. Files are batched to fit within the LLM context window
5. Each batch is sent to your configured LLM API
6. Results are aggregated into a structured report
7. Report is rendered in the terminal or saved to a file
```

## Configuration File

Location: `~/.config/ctxqa/config.json`

```json
{
  "provider": "openai",
  "api_key": "your-api-key",
  "base_url": "https://api.openai.com/v1",
  "model": "gpt-4o",
  "max_file_size": 1048576,
  "batch_size": 50,
  "max_batch_chars": 100000,
  "timeout_seconds": 120,
  "retry_count": 3,
  "context": {
    "exclude": [
      ".git/",
      "node_modules/",
      "dist/",
      "build/"
    ],
    "collect_history": true,
    "collect_commits": true
  },
  "defaults": {
    "severity": "all",
    "format": "text"
  }
}
```

**Security note**: The config file is created with `0600` permissions (owner read/write only). Your API key is never exposed in logs or terminal output.

## Built-in Exclusions

By default, ctxqa skips these directories and files:

**Directories**: `.git/`, `node_modules/`, `vendor/`, `.terraform/`, `.aws/`, `.azure/`, `dist/`, `build/`, `.output/`, `target/`, `__pycache__/`, `.pytest_cache/`, `.mypy_cache/`

**Files**: `*.lock`, `*.sum`, `*.pb.go`, `*.min.js`, `*.min.css`, `*.png`, `*.jpg`, `*.gif`, `*.ico`, `*.svg`, `*.woff`, `*.woff2`, `*.ttf`, `*.zip`, `*.tar`, `*.gz`, `*.rar`, `*.7z`, `*.exe`, `*.dll`, `*.so`, `*.dylib`

Add custom exclusions:
```bash
ctxqa config exclude-add "**/*.test.ts" "**/__snapshots__/**"
```

## Examples

### Analyze a project folder

```bash
$ cd ~/projects/my-app
$ ctxqa analyze

============================================================
           FOLDER ANALYSIS REPORT
============================================================

Folder: /home/user/projects/my-app
Total files: 247
Total size: 4.2 MB
Generated: 2026-07-15 14:30:00
Version: 0.1.0

--- File Structure and Organization ---

Current structure:
  src/
    components/ (45 files)
    pages/ (12 files)
    utils/ (8 files)
  tests/ (34 files)
  docs/ (15 files)
  config/ (6 files)

Grouping suggestions:
  - Consider merging duplicate utility functions in src/utils/
  - Move test files into tests/unit/ and tests/integration/
  - Consolidate configuration files into config/

--- Issues (2 critical, 5 errors, 12 warnings, 8 infos) ---

[CRITICAL]
  src/auth/login.go:42 - Password hashing uses weak algorithm
    Suggestion: Use bcrypt or argon2 for password hashing

[ERROR]
  src/api/client.go:89 - Missing timeout on HTTP request
  src/db/connection.go:34 - Hardcoded database credentials

--- Quality Scores ---

  Naming:            8/10
  Error Handling:    6/10
  Duplication:       7/10
  Security:          5/10
  Overall:           7/10

--- Improvement Suggestions ---

[HIGH] Implement centralized error handling
  Add a global error handler middleware to standardize error responses
  Files: src/api/*.go

[MEDIUM] Extract configuration to environment variables
  Remove hardcoded values and use os.Getenv() with defaults
  Files: src/db/connection.go, src/api/client.go
```

### Export as JSON

```bash
ctxqa analyze --format json -o report.json
```

### Incremental analysis

```bash
# First run analyzes all files
ctxqa analyze

# Subsequent runs only analyze changed files
ctxqa analyze --incremental
```

### Dry run to preview

```bash
ctxqa analyze --dry-run

Found 247 files:
  main.go (code, 234 bytes)
  app.ts (code, 156 bytes)
  config.yaml (config, 89 bytes)
  ...
```

## Development

```bash
# Clone the repository
git clone https://github.com/Benjamin-JHou/repo-audit.git
cd repo-audit

# Build
go build -o ctxqa ./cmd/ctxqa/

# Run
./ctxqa analyze --dry-run

# Run tests
go test ./...

# Build for all platforms
make build
```

## Project Structure

```
repo-audit/
├── cmd/ctxqa/              # CLI commands
│   ├── main.go
│   └── cmd/
│       ├── root.go         # Root command
│       ├── analyze.go      # Analyze command
│       └── config.go       # Config command
├── pkg/                    # Core packages
│   ├── api/                # API clients (OpenAI, Anthropic)
│   ├── config/             # Configuration management
│   ├── context/            # Context enrichment
│   ├── filesystem/         # Folder scanning and batching
│   ├── progress/           # Progress tracking
│   ├── report/             # Report generation
│   └── renderer/           # Terminal rendering
├── scripts/
│   └── install.sh          # One-click install script
├── .github/workflows/      # CI/CD
├── Makefile
├── go.mod
├── README.md
├── LICENSE
└── CONTRIBUTING.md
```

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

1. Fork the repository
2. Create a feature branch: `git checkout -b feat/my-feature`
3. Make your changes
4. Run tests: `go test ./...`
5. Commit: `git commit -m "feat: add my feature"`
6. Push: `git push origin feat/my-feature`
7. Open a Pull Request

## Roadmap

- [ ] Streaming output for real-time analysis
- [ ] Parallel batch processing with configurable concurrency
- [ ] Plugin system for custom analyzers
- [ ] Diff mode to track changes between runs
- [ ] Export to HTML with syntax highlighting
- [ ] Integration with popular CI/CD platforms
- [ ] Support for custom LLM providers via plugins

## License

Apache 2.0 - see [LICENSE](LICENSE) for details.

## Acknowledgments

- Built with [Go](https://golang.org/)
- CLI framework: [Cobra](https://github.com/spf13/cobra)
- Terminal styling: [charmbracelet](https://github.com/charmbracelet)
- Git operations: [go-git](https://github.com/go-git/go-git)
