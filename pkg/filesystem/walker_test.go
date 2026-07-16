package filesystem

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ctxqa/ctxqa/pkg/config"
)

func TestClassifyFile(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"main.go", "code"},
		{"app.ts", "code"},
		{"component.tsx", "code"},
		{"script.js", "code"},
		{"config.yaml", "config"},
		{"README.md", "document"},
		{"data.csv", "data"},
		{"unknown.xyz", "other"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := classifyFile(tt.path)
			if result != tt.expected {
				t.Errorf("classifyFile(%s) = %s, want %s", tt.path, result, tt.expected)
			}
		})
	}
}

func TestScan(t *testing.T) {
	tmpDir := t.TempDir()

	testFiles := map[string]string{
		"main.go":       "package main\nfunc main() {}",
		"app.ts":        "console.log('hello');",
		"config.yaml":   "key: value",
		"README.md":     "# Test",
		"notes.txt":     "some notes",
		"skip.png":      "binary data",
		"skip.lock":     "lock file",
		"sub/helper.go": "package sub\nfunc Helper() {}",
	}

	for path, content := range testFiles {
		fullPath := filepath.Join(tmpDir, path)
		os.MkdirAll(filepath.Dir(fullPath), 0755)
		os.WriteFile(fullPath, []byte(content), 0644)
	}

	scanner := NewFileScanner(&config.Config{
		MaxFileSize: 1024 * 1024,
	})

	files, err := scanner.Scan(tmpDir)
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}

	if len(files) != 6 {
		t.Errorf("expected 6 files, got %d", len(files))
	}

	for _, f := range files {
		if f.Path == "skip.png" || f.Path == "skip.lock" {
			t.Errorf("should have skipped %s", f.Path)
		}
	}
}

func TestBatchSplit(t *testing.T) {
	batcher := NewBatcher(1000, 1024*1024)

	files := []FileContent{
		{Path: "a.go", Content: "package main\n", FileType: "code", Size: 12},
		{Path: "b.ts", Content: "console.log('hi');\n", FileType: "code", Size: 18},
		{Path: "c.py", Content: "print('hello')\n", FileType: "code", Size: 16},
	}

	batches := batcher.Split(files)
	if len(batches) != 1 {
		t.Errorf("expected 1 batch, got %d", len(batches))
	}

	if batches[0].Total != 3 {
		t.Errorf("expected total 3, got %d", batches[0].Total)
	}
}
