# ctxqa - AI-Powered Folder Analyzer

A terminal CLI tool that analyzes any folder using AI. Traverses the entire directory tree, sends file contents to an LLM in batches, and outputs a structured report covering file structure, organization suggestions, progress analysis, issue detection, quality scoring, and improvement recommendations.

## Features

- **Any Folder Analysis**: Works on any directory, not just code repositories
- **Intelligent File Classification**: Automatically classifies files into code, config, data, document, web, and other categories
- **Batch Processing**: Handles large folders by batching files for API calls
- **Multiple API Providers**: Supports OpenAI and Anthropic (Claude) compatible APIs
- **Structured Reports**: Generates reports with file structure analysis, organization suggestions, severity-classified issues, quality scores, and prioritized suggestions
- **Incremental Analysis**: Only analyzes changed files since last run
- **Resume Support**: Interrupted analyses can be resumed from checkpoints
- **Multiple Output Formats**: text, JSON, and Markdown output
- **One-Click Install**: Single command installation via shell script

## Installation

```bash
curl -sSL https://raw.githubusercontent.com/Benjamin-JHou/repo-audit/main/scripts/install.sh | sh
```

## Quick Start

```bash
# Configure your API key
ctxqa config init

# Analyze current folder
ctxqa analyze

# Analyze a specific folder
ctxqa analyze --dir /path/to/folder

# Dry run to see which files will be analyzed
ctxqa analyze --dry-run

# Incremental analysis (only changed files)
ctxqa analyze --incremental

# Export as Markdown
ctxqa analyze --format markdown -o report.md
```

## Configuration

```bash
# Show current configuration
ctxqa config show

# Switch provider
ctxqa config set-provider anthropic

# Set model
ctxqa config set-model claude-sonnet-4-20250514

# Set API key
ctxqa config set-api-key

# Add exclusion patterns
ctxqa config exclude-add "**/*.test.ts" "**/__snapshots__/**"

# Reset configuration
ctxqa config reset
```

## Command Reference

### `ctxqa analyze`

Perform a comprehensive analysis of any folder.

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--format` | `-f` | `text` | Output format: text, json, markdown |
| `-o`, `--output` | | stdout | Output file path |
| `--severity` | | `all` | Minimum severity: all, info, warning, error, critical |
| `--section` | | `all` | Specific section: overview, structure, progress, issues, suggestions |
| `--incremental` | | `false` | Only analyze changed files |
| `--resume` | | `false` | Resume from last interrupted analysis |
| `--no-context` | | `false` | Skip context collection |
| `--dir` | `-d` | `.` | Folder to analyze |
| `--concurrency` | `-j` | `1` | Parallel batch processors |
| `--dry-run` | | `false` | List files only |
| `--summary` | | `false` | Brief summary only |

### `ctxqa config`

Manage configuration settings.

| Command | Description |
|---------|-------------|
| `ctxqa config init` | Interactive configuration wizard |
| `ctxqa config show` | Display current configuration |
| `ctxqa config set-provider [openai\|anthropic]` | Set API provider |
| `ctxqa config set-model [name]` | Set model name |
| `ctxqa config set-api-key` | Set API key interactively |
| `ctxqa config exclude-add [patterns]` | Add exclusion patterns |
| `ctxqa config exclude-remove [patterns]` | Remove exclusion patterns |
| `ctxqa config reset` | Reset to defaults |

## Report Structure

The analysis report includes:

1. **Folder Overview**: Total files, total size, file type breakdown, top-level directories, tech stack detection
2. **File Structure and Organization Suggestions**: Current layout summary, proposed reorganization plan, file grouping recommendations, naming issues, duplicate file detection
3. **Progress Analysis**: Completed work, work in progress, pending tasks based on file content and commit history
4. **Issues Found**: Classified by severity (critical, error, warning, info)
5. **Quality Scores**: Named scoring for naming, error handling, duplication, security
6. **Improvement Suggestions**: Prioritized (high, medium, low) actionable recommendations

## Supported API Providers

### OpenAI Compatible

- Default URL: `https://api.openai.com/v1`
- Models: `gpt-4o`, `gpt-4o-mini`, `gpt-4-turbo`

### Anthropic Compatible

- Default URL: `https://api.anthropic.com/v1`
- Models: `claude-sonnet-4-20250514`, `claude-opus-4-0-20250620`

## How It Works

1. You run `ctxqa analyze` in any folder
2. ctxqa scans all files recursively, skipping binaries and excluded paths
3. Files are batched to fit within the LLM context window
4. Each batch is sent to your configured LLM API
5. Results are aggregated into a structured report
6. Report is rendered in the terminal or saved to a file

## Development

```bash
# Build
go build -o ctxqa ./cmd/ctxqa/

# Run
./ctxqa analyze --dry-run

# Test
go test ./...
```

## License

Apache 2.0
