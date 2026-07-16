package filesystem

import (
	"bufio"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/ctxqa/ctxqa/pkg/config"
)

type FileContent struct {
	Path     string `json:"path"`
	Content  string `json:"content"`
	FileType string `json:"file_type"`
	Size     int    `json:"size"`
}

type FileScanner struct {
	excludePatterns []string
	includePatterns []string
	maxFileSize     int
}

func NewFileScanner(cfg *config.Config) *FileScanner {
	return &FileScanner{
		excludePatterns: cfg.Context.Exclude,
		includePatterns: cfg.Context.Include,
		maxFileSize:     cfg.MaxFileSize,
	}
}

func (fs_ *FileScanner) Scan(root string) ([]FileContent, error) {
	var files []FileContent

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(root, path)
		if err != nil {
			rel = path
		}

		if d.IsDir() {
			return fs_.shouldSkipDir(rel, d)
		}

		if fs_.shouldSkipFile(rel, d) {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return nil
		}

		if info.Size() > int64(fs_.maxFileSize) {
			return nil
		}

		content, err := fs_.readFile(path)
		if err != nil {
			return nil
		}

		files = append(files, FileContent{
			Path:     rel,
			Content:  content,
			FileType: classifyFile(path),
			Size:     int(info.Size()),
		})

		return nil
	})

	return files, err
}

func (fs_ *FileScanner) shouldSkipDir(rel string, d fs.DirEntry) error {
	name := d.Name()

	skipDirs := map[string]bool{
		".git": true, "node_modules": true, "vendor": true,
		".terraform": true, ".aws": true, ".azure": true,
		"dist": true, "build": true, ".output": true, "target": true,
		"__pycache__": true, ".pytest_cache": true, ".mypy_cache": true,
	}

	if skipDirs[name] {
		return filepath.SkipDir
	}

	for _, pattern := range fs_.excludePatterns {
		if strings.HasSuffix(pattern, "/") {
			dirName := strings.TrimSuffix(pattern, "/")
			if strings.Contains(rel, dirName) {
				return filepath.SkipDir
			}
		}
	}

	return nil
}

func (fs_ *FileScanner) shouldSkipFile(rel string, d fs.DirEntry) bool {
	name := d.Name()

	skipExtensions := map[string]bool{
		".lock": true, ".sum": true, ".pb.go": true,
		".min.js": true, ".min.css": true,
		".png": true, ".jpg": true, ".jpeg": true, ".gif": true,
		".ico": true, ".svg": true, ".woff": true,
		".woff2": true, ".ttf": true, ".zip": true, ".tar": true,
		".gz": true, ".rar": true, ".7z": true, ".exe": true,
		".dll": true, ".so": true, ".dylib": true,
	}

	if skipExtensions[strings.ToLower(filepath.Ext(name))] {
		return true
	}

	for _, pattern := range fs_.excludePatterns {
		if !strings.HasSuffix(pattern, "/") {
			matched, err := filepath.Match(pattern, name)
			if err == nil && matched {
				return true
			}
			matched, err = filepath.Match(pattern, rel)
			if err == nil && matched {
				return true
			}
		}
	}

	return false
}

func (fs_ *FileScanner) readFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var sb strings.Builder
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	for scanner.Scan() {
		sb.WriteString(scanner.Text())
		sb.WriteString("\n")
	}

	return sb.String(), scanner.Err()
}

func classifyFile(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	classMap := map[string]string{
		".go": "code", ".ts": "code", ".tsx": "code", ".js": "code",
		".jsx": "code", ".py": "code", ".rb": "code", ".rs": "code",
		".java": "code", ".c": "code", ".cpp": "code", ".h": "code",
		".hpp": "code", ".cs": "code", ".swift": "code", ".kt": "code",
		".scala": "code", ".php": "code", ".sh": "code", ".bash": "code",
		".zsh": "code", ".yaml": "config", ".yml": "config",
		".toml": "config", ".json": "config", ".xml": "config",
		".ini": "config", ".conf": "config", ".env": "config",
		".csv": "data", ".parquet": "data", ".avro": "data",
		".db": "data", ".sqlite": "data",
		".md": "document", ".txt": "document", ".pdf": "document",
		".doc": "document", ".docx": "document",
		".html": "web", ".css": "web", ".scss": "web",
		".vue": "web", ".svelte": "web",
		".proto": "other", ".graphql": "other", ".sql": "other",
	}

	if t, ok := classMap[ext]; ok {
		return t
	}
	return "other"
}

