# Requirements Document: AI-Powered Folder Analyzer CLI

## Introduction

This project aims to build a terminal CLI tool named `ctxqa` that provides AI-powered analysis and organization suggestions for any folder. The tool is fully open source, supports user-provided API keys, one-click installation, and operates entirely in the terminal. Core value: when run in any folder, it traverses the entire directory tree, sends file contents to an LLM in batches, and produces a structured report covering folder/file structure, organization suggestions, progress analysis, issue detection, quality scoring, and improvement recommendations.

## Glossary

- **System**: `ctxqa` CLI tool
- **Folder analysis**: Automatic traversal of all files in a folder and its subfolders, sent in batches to an LLM for comprehensive review
- **Audit report**: Structured output containing progress analysis, file structure, organization suggestions, issues, quality scores, and improvement recommendations
- **API Key**: User-provided LLM service credential for calling the model API
- **One-click install**: Single command installation and initialization
- **Context**: Auxiliary project information including file structure, Git branch, command history

## Requirements

### Requirement 1: One-Click Install

**User Story**: As a developer, I want to install ctxqa with a single command so I can start using it quickly without complex setup.

#### Acceptance Criteria

1. WHEN user runs `curl -sSL https://ctxqa.dev/install.sh | sh`, System  SHALL download and install the CLI tool to a directory in the user PATH
2. WHEN user runs `curl -sSL https://ctxqa.dev/install.sh | sh`, System  SHALL auto-detect OS and architecture and select the matching binary
3. WHEN the install target directory is not writable, System  SHALL fall back to `$HOME/.local/bin`
4. WHEN installation completes, System  SHALL output a confirmation message showing available subcommands
5. IF the downloaded binary checksum does not match, System  SHALL abort installation and output an error message

### Requirement 2: User-Provided API Key

**User Story**: As a developer, I want to configure my own API key on first use so I can freely choose my LLM backend provider.

#### Acceptance Criteria

1. WHEN user runs ctxqa for the first time without a configured API key, System  SHALL enter an interactive configuration wizard prompting for API key, base URL, and model name
2. WHILE the configuration wizard is running, System  SHALL save the configuration to `~/.config/ctxqa/config.json`
3. WHEN user runs `ctxqa config set-api-key`, System  SHALL accept the new API key and update the config file
4. WHEN user runs `ctxqa config show`, System  SHALL display current configuration without revealing the API key value
5. IF the provided API key is invalid, System  SHALL output an error message and guide the user to reconfigure

### Requirement 3: Folder Traversal and Batch Analysis

**User Story**: As a developer, I want ctxqa to automatically traverse all files in any folder and its subfolders, sending them in batches to an LLM for review, so I get a comprehensive analysis of the folder contents.

#### Acceptance Criteria

1. WHEN user runs `ctxqa audit`, System  SHALL automatically traverse all files in the current working directory and its subdirectories
2. WHEN the file count exceeds a single API call limit, System  SHALL batch the files so each batch does not exceed 80% of the model context window
3. WHILE analysis is in progress, System  SHALL display real-time progress including current file count and percentage
4. WHEN encountering binary files, build artifacts, or files in excluded directories, System  SHALL automatically skip them and log the skip reason
5. IF reading a file fails, System  SHALL log the error and continue processing other files without interrupting the overall flow
6. WHEN all batches are complete, System  SHALL aggregate all batch results and generate a complete analysis report

### Requirement 4: Proactive Analysis Report

**User Story**: As a developer, I want ctxqa to output a structured report after analyzing a folder, covering file structure, organization suggestions, progress analysis, issues, quality scores, and improvement recommendations.

#### Acceptance Criteria

1. WHEN analysis completes, System  SHALL output a structured report containing the following sections:
   - Folder overview (total files, total size, file type breakdown, tech stack detection)
   - File structure and organization suggestions (proposed directory layout, file grouping recommendations)
   - Progress analysis (completed features, in-progress items, pending work, based on file content and commit history)
   - Issues found (classified by severity: critical, error, warning, info)
   - Quality scores (naming conventions, error handling, code duplication, security, with 0-10 ratings)
   - Improvement suggestions (prioritized: high, medium, low)
2. WHEN user runs `ctxqa audit --format json`, System  SHALL output the report as JSON to stdout
3. WHEN user runs `ctxqa audit --format markdown -o report.md`, System  SHALL save the report as a Markdown file
4. WHEN user runs `ctxqa audit --severity warning`, System  SHALL only output issues at warning severity and above
5. WHEN user runs `ctxqa audit --section progress`, System  SHALL only output the specified section
6. WHEN user runs `ctxqa audit --diff`, System  SHALL compare with the last analysis result and highlight changes

### Requirement 5: Incremental Analysis

**User Story**: As a developer, I want ctxqa to support incremental analysis so I can quickly see the impact of changes.

#### Acceptance Criteria

1. WHEN user runs `ctxqa audit --incremental`, System  SHALL only analyze files that have changed since the last analysis
2. WHEN running incremental analysis, System  SHALL identify changed files using Git diff if the folder is a Git repository
3. IF the current folder is not a Git repository, System  SHALL identify changed files using file modification timestamps
4. WHEN incremental analysis completes, System  SHALL output a change impact report including newly introduced issues, fixed issues, and remaining issues

### Requirement 6: Context Enhancement

**User Story**: As a developer, I want ctxqa to include Git branch, command history, and commit information in the analysis report for better project context.

#### Acceptance Criteria

1. WHEN user runs `ctxqa audit`, System  SHALL automatically collect the current Git branch name, last 20 shell commands, and last 10 commit messages
2. WHEN the working directory is not a Git repository, System  SHALL skip Git-related context collection and note this in the report
3. WHEN generating the analysis report, System  SHALL include the collected context as an introduction section
4. WHEN user runs `ctxqa audit --no-context`, System  SHALL skip all context collection and only analyze file contents

### Requirement 7: Open Source and Distribution

**User Story**: As a project maintainer, I want ctxqa to be fully open source with one-click installation so the community can contribute and easily install it.

#### Acceptance Criteria

1. WHEN user runs `ctxqa --version`, System  SHALL output the version number
2. WHEN user runs `ctxqa --license`, System  SHALL output the Apache 2.0 license text
3. WHEN user runs `ctxqa --help`, System  SHALL output complete help information including all subcommands and flags
4. System  SHALL provide a complete source code repository on GitHub with README, LICENSE, and CONTRIBUTING files

### Requirement 8: Configuration Management

**User Story**: As a developer, I want to manage analysis settings including exclusion rules, provider selection, and model configuration.

#### Acceptance Criteria

1. WHEN user runs `ctxqa config show`, System  SHALL display all current configuration settings
2. WHEN user runs `ctxqa config set-provider openai`, System  SHALL switch to the OpenAI provider
3. WHEN user runs `ctxqa config set-provider anthropic`, System  SHALL switch to the Anthropic provider
4. WHEN user runs `ctxqa config set-model <name>`, System  SHALL set the model name
5. WHEN user runs `ctxqa config set max-file-size <bytes>`, System  SHALL set the maximum single file size
6. WHEN user runs `ctxqa config set batch-size <count>`, System  SHALL set the number of files per batch
7. WHEN user runs `ctxqa config exclude-add <patterns...>`, System  SHALL add file exclusion patterns
8. WHEN user runs `ctxqa config exclude-remove <patterns...>`, System  SHALL remove the specified exclusion patterns

## API Provider Support

### OpenAI Compatible

- Base URL: user-configurable (default `https://api.openai.com/v1`)
- Interface format: OpenAI Chat Completions API
- Supported models: gpt-4o, gpt-4o-mini, gpt-4-turbo

### Anthropic Compatible

- Base URL: user-configurable (default `https://api.anthropic.com/v1`)
- Interface format: Anthropic Messages API
- Supported models: claude-sonnet-4-20250514, claude-opus-4-0-20250620

## Non-Functional Requirements

1. System  SHALL be developed in Go, producing a single binary with no runtime dependencies
2. System  SHALL have a binary size no larger than 20MB after installation
3. System  SHALL complete folder analysis within 30 minutes for folders containing up to 1000 files
4. System  SHALL support macOS, Linux, and Windows platforms
5. System  SHALL use JSON for configuration files
6. System  SHALL support interrupt (Ctrl+C) with progress persistence, resumable from the last checkpoint
7. System  SHALL keep memory usage under 500MB during analysis
