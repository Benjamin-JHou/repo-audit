package context

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ctxqa/ctxqa/pkg/config"
	"github.com/go-git/go-git/v5"
)

type Context struct {
	Branch      string    `json:"branch"`
	History     []string  `json:"history"`
	Commits     []string  `json:"commits"`
	ChangedFiles []string `json:"changed_files"`
	TechStack   []string  `json:"tech_stack"`
	CollectedAt time.Time `json:"collected_at"`
}

func Collect(cfg *config.Config, dir string, incremental bool) (*Context, error) {
	ctx := &Context{
		CollectedAt: time.Now(),
	}

	ctx.Branch = collectBranch(dir)
	ctx.History = collectHistory()
	ctx.Commits = collectCommits(dir, cfg.Context.CollectCommits)
	ctx.ChangedFiles = collectChangedFiles(dir, incremental)
	ctx.TechStack = detectTechStack(dir)

	return ctx, nil
}

func collectBranch(dir string) string {
	repo, err := git.PlainOpen(dir)
	if err != nil {
		return "unavailable"
	}

	ref, err := repo.Head()
	if err != nil {
		return "unavailable"
	}

	name := ref.Name().Short()
	if name == "HEAD" {
		return "detached"
	}
	return name
}

func collectHistory() []string {
	shell := config.DetectShell()
	var historyPath string

	switch shell {
	case "zsh":
		home, _ := os.UserHomeDir()
		historyPath = filepath.Join(home, ".zsh_history")
	default:
		home, _ := os.UserHomeDir()
		historyPath = filepath.Join(home, ".bash_history")
	}

	var lines []string
	file, err := os.Open(historyPath)
	if err != nil {
		return lines
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() && count < 20 {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line != "" {
			lines = append(lines, line)
			count++
		}
	}
	return lines
}

func collectCommits(dir string, enabled bool) []string {
	if !enabled {
		return nil
	}

	repo, err := git.PlainOpen(dir)
	if err != nil {
		return nil
	}

	commitIter, err := repo.Log(&git.LogOptions{})
	if err != nil {
		return nil
	}
	defer commitIter.Close()

	var commits []string
	count := 0
	for {
		c, err := commitIter.Next()
		if err != nil {
			break
		}
		if count >= 10 {
			break
		}
		msg := strings.TrimSpace(c.Message)
		msg = strings.ReplaceAll(msg, "\n", " ")
		commits = append(commits, c.Hash.String()[:8]+" "+msg)
		count++
	}
	return commits
}

func collectChangedFiles(dir string, incremental bool) []string {
	if !incremental {
		return nil
	}

	repo, err := git.PlainOpen(dir)
	if err != nil {
		return nil
	}

	w, err := repo.Worktree()
	if err != nil {
		return nil
	}

	status, err := w.Status()
	if err != nil {
		return nil
	}

	var changed []string
	for file, s := range status {
		if s.Worktree != git.Unmodified {
			changed = append(changed, file)
		}
	}
	return changed
}

func detectTechStack(dir string) []string {
	files := []string{
		"package.json", "go.mod", "Cargo.toml", "pom.xml",
		"build.gradle", "Gemfile", "composer.json",
		"requirements.txt", "setup.py", "pyproject.toml",
		"dotnet.csproj", "project.clj", "mix.exs",
	}

	var techStack []string
	for _, f := range files {
		path := filepath.Join(dir, f)
		if _, err := os.Stat(path); err == nil {
			techStack = append(techStack, f)
		}
	}

	return techStack
}
