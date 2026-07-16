package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := Get()
	if cfg.Provider != DefaultProvider {
		t.Errorf("expected provider %s, got %s", DefaultProvider, cfg.Provider)
	}
	if cfg.BaseURL != DefaultBaseURL {
		t.Errorf("expected base URL %s, got %s", DefaultBaseURL, cfg.BaseURL)
	}
	if cfg.Model != DefaultModel {
		t.Errorf("expected model %s, got %s", DefaultModel, cfg.Model)
	}
}

func TestLoadEmptyConfig(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "nonexistent.json")

	cfg, err := Load(cfgPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Provider != DefaultProvider {
		t.Errorf("expected default provider, got %s", cfg.Provider)
	}
}

func TestSaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "config.json")

	cfg := &Config{
		Provider:       "anthropic",
		APIKey:         "test-key-123",
		BaseURL:        "https://test.api/v1",
		Model:          "claude-test",
		MaxFileSize:    2048,
		BatchSize:      25,
		MaxBatchChars:  50000,
		TimeoutSeconds: 60,
		RetryCount:     5,
	}

	if err := cfg.Save(cfgPath); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	// Verify file permissions
	info, err := os.Stat(cfgPath)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600 permissions, got %v", info.Mode().Perm())
	}

	loaded, err := Load(cfgPath)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}
	if loaded.Provider != "anthropic" {
		t.Errorf("expected anthropic, got %s", loaded.Provider)
	}
	if loaded.Model != "claude-test" {
		t.Errorf("expected claude-test, got %s", loaded.Model)
	}
}

func TestDefaultExcludes(t *testing.T) {
	excludes := defaultExcludes()
	if len(excludes) == 0 {
		t.Error("expected default excludes, got none")
	}

	expected := []string{".git/", "node_modules/", "vendor/"}
	for _, exp := range expected {
		found := false
		for _, ex := range excludes {
			if ex == exp {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected %s in default excludes", exp)
		}
	}
}
