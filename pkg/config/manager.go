package config

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

const (
	DefaultProvider   = "openai"
	DefaultBaseURL    = "https://api.openai.com/v1"
	DefaultModel      = "gpt-4o"
	DefaultMaxFileSize  = 1048576 // 1MB
	DefaultBatchSize    = 50
	DefaultMaxBatchChars = 100000
	DefaultTimeout      = 120
	DefaultRetryCount   = 3
	configDir           = ".config/ctxqa"
	configFileName      = "config.json"
)

type Config struct {
	Provider        string           `json:"provider"`
	APIKey          string           `json:"api_key"`
	BaseURL         string           `json:"base_url"`
	Model           string           `json:"model"`
	MaxFileSize     int              `json:"max_file_size"`
	BatchSize       int              `json:"batch_size"`
	MaxBatchChars   int              `json:"max_batch_chars"`
	TimeoutSeconds  int              `json:"timeout_seconds"`
	RetryCount      int              `json:"retry_count"`
	Context         ContextConfig    `json:"context"`
	Defaults        DefaultConfig    `json:"defaults"`
}

type ContextConfig struct {
	Exclude        []string `json:"exclude"`
	Include        []string `json:"include"`
	CollectHistory bool     `json:"collect_history"`
	CollectCommits bool     `json:"collect_commits"`
}

type DefaultConfig struct {
	Severity string `json:"severity"`
	Format   string `json:"format"`
}

var (
	instance *Config
	once     sync.Once
	mu       sync.RWMutex
)

func Get() *Config {
	once.Do(func() {
		instance = &Config{
			Provider:       DefaultProvider,
			BaseURL:        DefaultBaseURL,
			Model:          DefaultModel,
			MaxFileSize:    DefaultMaxFileSize,
			BatchSize:      DefaultBatchSize,
			MaxBatchChars:  DefaultMaxBatchChars,
			TimeoutSeconds: DefaultTimeout,
			RetryCount:     DefaultRetryCount,
			Context: ContextConfig{
				Exclude:        defaultExcludes(),
				CollectHistory: true,
				CollectCommits: true,
			},
			Defaults: DefaultConfig{
				Severity: "all",
				Format:   "text",
			},
		}
	})
	mu.RLock()
	defer mu.RUnlock()
	return instance
}

func Load(path string) (*Config, error) {
	cfg := &Config{
		Provider:       DefaultProvider,
		BaseURL:        DefaultBaseURL,
		Model:          DefaultModel,
		MaxFileSize:    DefaultMaxFileSize,
		BatchSize:      DefaultBatchSize,
		MaxBatchChars:  DefaultMaxBatchChars,
		TimeoutSeconds: DefaultTimeout,
		RetryCount:     DefaultRetryCount,
		Context: ContextConfig{
			Exclude:        defaultExcludes(),
			CollectHistory: true,
			CollectCommits: true,
		},
		Defaults: DefaultConfig{
			Severity: "all",
			Format:   "text",
		},
	}

	if path == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return cfg, nil
		}
		path = filepath.Join(home, configDir, configFileName)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, nil
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return cfg, nil
	}

	cfg.normalize()
	return cfg, nil
}

func (c *Config) Save(path string) error {
	if path == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		dir := filepath.Join(home, configDir)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
		path = filepath.Join(dir, configFileName)
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

func (c *Config) normalize() {
	if c.Provider == "" {
		c.Provider = DefaultProvider
	}
	if c.BaseURL == "" {
		c.BaseURL = DefaultBaseURL
	}
	if c.Model == "" {
		c.Model = DefaultModel
	}
	if c.MaxFileSize <= 0 {
		c.MaxFileSize = DefaultMaxFileSize
	}
	if c.BatchSize <= 0 {
		c.BatchSize = DefaultBatchSize
	}
	if c.MaxBatchChars <= 0 {
		c.MaxBatchChars = DefaultMaxBatchChars
	}
	if c.TimeoutSeconds <= 0 {
		c.TimeoutSeconds = DefaultTimeout
	}
	if c.RetryCount <= 0 {
		c.RetryCount = DefaultRetryCount
	}
}

func defaultExcludes() []string {
	return []string{
		".git/",
		"node_modules/",
		"vendor/",
		".terraform/",
		".aws/",
		".azure/",
		"dist/",
		"build/",
		".output/",
		"target/",
		"__pycache__/",
		".pytest_cache/",
		".mypy_cache/",
		"*.lock",
		"*.sum",
		"*.pb.go",
		"*.min.js",
		"*.min.css",
		"*.png",
		"*.jpg",
		"*.gif",
		"*.ico",
		"*.svg",
		"*.woff",
		"*.woff2",
		"*.ttf",
		filepath.Join(".config", "ctxqa", "**"),
	}
}

func DetectShell() string {
	shell := os.Getenv("SHELL")
	if shell != "" {
		return filepath.Base(shell)
	}
	if _, err := exec.LookPath("zsh"); err == nil {
		return "zsh"
	}
	return "bash"
}
