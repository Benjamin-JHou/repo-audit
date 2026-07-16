package report

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type AuditReport struct {
	Metadata    ReportMetadata   `json:"metadata"`
	Overview    OverviewSection  `json:"overview"`
	Structure   StructureSection `json:"structure"`
	Progress    ProgressSection  `json:"progress"`
	Issues      IssuesSection    `json:"issues"`
	Quality     QualitySection   `json:"quality"`
	Suggestions []Suggestion     `json:"suggestions"`
}

type ReportMetadata struct {
	FolderPath        string            `json:"folder_path"`
	Branch            string            `json:"branch"`
	TotalFiles        int               `json:"total_files"`
	TotalSize         int64             `json:"total_size"`
	FileTypeBreakdown map[string]int    `json:"file_type_breakdown"`
	TopDirectories    []DirStat         `json:"top_directories"`
	TechStack         []string          `json:"tech_stack"`
	GeneratedAt       time.Time         `json:"generated_at"`
	Version           string            `json:"version"`
}

type DirStat struct {
	Path  string `json:"path"`
	Files int    `json:"files"`
	Size  int64  `json:"size"`
}

type OverviewSection struct {
	FileCountByExt map[string]int `json:"file_count_by_ext"`
	TopDirs        []DirStat      `json:"top_directories"`
}

type StructureSection struct {
	CurrentStructure    string   `json:"current_structure"`
	ProposedStructure   string   `json:"proposed_structure"`
	GroupingSuggestions []string `json:"grouping_suggestions"`
	NamingIssues        []string `json:"naming_issues"`
	DuplicateFiles      []string `json:"duplicate_files"`
}

type ProgressSection struct {
	Completed []Feature `json:"completed"`
	Running   []Feature `json:"in_progress"`
	Pending   []Feature `json:"pending"`
}

type Feature struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Files       []string `json:"files"`
	Status      string   `json:"status"`
}

type IssuesSection struct {
	Critical []Issue `json:"critical"`
	Error    []Issue `json:"error"`
	Warning  []Issue `json:"warning"`
	Info     []Issue `json:"info"`
}

type Issue struct {
	File       string `json:"file"`
	Line       int    `json:"line"`
	Severity   string `json:"severity"`
	Category   string `json:"category"`
	Message    string `json:"message"`
	Suggestion string `json:"suggestion,omitempty"`
}

type QualitySection struct {
	NamingScore        int `json:"naming_score"`
	ErrorHandlingScore int `json:"error_handling_score"`
	DuplicationScore   int `json:"duplication_score"`
	SecurityScore      int `json:"security_score"`
	OverallScore       int `json:"overall_score"`
}

type Suggestion struct {
	Priority string   `json:"priority"`
	Category string   `json:"category"`
	Title    string   `json:"title"`
	Detail   string   `json:"detail"`
	Files    []string `json:"files"`
}

func (r *AuditReport) FormatText(version string) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("\n%s\n", strings.Repeat("=", 60)))
	sb.WriteString("           FOLDER ANALYSIS REPORT\n")
	sb.WriteString(fmt.Sprintf("%s\n\n", strings.Repeat("=", 60)))

	sb.WriteString(fmt.Sprintf("Folder: %s\n", r.Metadata.FolderPath))
	sb.WriteString(fmt.Sprintf("Total files: %d\n", r.Metadata.TotalFiles))
	sb.WriteString(fmt.Sprintf("Total size: %s\n", formatBytes(r.Metadata.TotalSize)))
	sb.WriteString(fmt.Sprintf("Generated: %s\n", r.Metadata.GeneratedAt.Format("2006-01-02 15:04:05")))
	sb.WriteString(fmt.Sprintf("Version: %s\n\n", version))

	if len(r.Structure.CurrentStructure) > 0 {
		sb.WriteString(fmt.Sprintf("--- File Structure and Organization ---\n\n"))
		sb.WriteString(fmt.Sprintf("%s\n\n", r.Structure.CurrentStructure))
		if len(r.Structure.GroupingSuggestions) > 0 {
			sb.WriteString("Grouping suggestions:\n")
			for _, g := range r.Structure.GroupingSuggestions {
				sb.WriteString(fmt.Sprintf("  - %s\n", g))
			}
			sb.WriteString("\n")
		}
		if len(r.Structure.NamingIssues) > 0 {
			sb.WriteString("Naming issues:\n")
			for _, n := range r.Structure.NamingIssues {
				sb.WriteString(fmt.Sprintf("  - %s\n", n))
			}
			sb.WriteString("\n")
		}
		if len(r.Structure.DuplicateFiles) > 0 {
			sb.WriteString("Duplicate files:\n")
			for _, d := range r.Structure.DuplicateFiles {
				sb.WriteString(fmt.Sprintf("  - %s\n", d))
			}
			sb.WriteString("\n")
		}
	}

	sb.WriteString(fmt.Sprintf("--- Issues (%d critical, %d errors, %d warnings, %d infos) ---\n\n",
		len(r.Issues.Critical), len(r.Issues.Error),
		len(r.Issues.Warning), len(r.Issues.Info)))

	if len(r.Issues.Critical) > 0 {
		sb.WriteString("[CRITICAL]\n")
		for _, issue := range r.Issues.Critical {
			sb.WriteString(fmt.Sprintf("  %s:%d - %s\n", issue.File, issue.Line, issue.Message))
			if issue.Suggestion != "" {
				sb.WriteString(fmt.Sprintf("    Suggestion: %s\n", issue.Suggestion))
			}
		}
		sb.WriteString("\n")
	}

	if len(r.Issues.Error) > 0 {
		sb.WriteString("[ERROR]\n")
		for _, issue := range r.Issues.Error {
			sb.WriteString(fmt.Sprintf("  %s:%d - %s\n", issue.File, issue.Line, issue.Message))
		}
		sb.WriteString("\n")
	}

	if len(r.Issues.Warning) > 0 {
		sb.WriteString("[WARNING]\n")
		for _, issue := range r.Issues.Warning {
			sb.WriteString(fmt.Sprintf("  %s:%d - %s\n", issue.File, issue.Line, issue.Message))
		}
		sb.WriteString("\n")
	}

	if len(r.Issues.Info) > 0 {
		sb.WriteString("[INFO]\n")
		for _, issue := range r.Issues.Info {
			sb.WriteString(fmt.Sprintf("  %s:%d - %s\n", issue.File, issue.Line, issue.Message))
		}
		sb.WriteString("\n")
	}

	sb.WriteString(fmt.Sprintf("--- Quality Scores ---\n\n"))
	sb.WriteString(fmt.Sprintf("  Naming:           %2d/10\n", r.Quality.NamingScore))
	sb.WriteString(fmt.Sprintf("  Error Handling:   %2d/10\n", r.Quality.ErrorHandlingScore))
	sb.WriteString(fmt.Sprintf("  Duplication:      %2d/10\n", r.Quality.DuplicationScore))
	sb.WriteString(fmt.Sprintf("  Security:         %2d/10\n", r.Quality.SecurityScore))
	sb.WriteString(fmt.Sprintf("  Overall:          %2d/10\n\n", r.Quality.OverallScore))

	if len(r.Suggestions) > 0 {
		sb.WriteString("--- Improvement Suggestions ---\n\n")
		for _, s := range r.Suggestions {
			sb.WriteString(fmt.Sprintf("[%s] %s\n", strings.ToUpper(s.Priority), s.Title))
			sb.WriteString(fmt.Sprintf("  %s\n", s.Detail))
			if len(s.Files) > 0 {
				sb.WriteString(fmt.Sprintf("  Files: %s\n", strings.Join(s.Files, ", ")))
			}
			sb.WriteString("\n")
		}
	}

	sb.WriteString(fmt.Sprintf("%s\n", strings.Repeat("=", 60)))

	return sb.String()
}

func (r *AuditReport) FormatJSON() string {
	data, _ := json.MarshalIndent(r, "", "  ")
	return string(data)
}

func (r *AuditReport) FormatMarkdown() string {
	var sb strings.Builder

	sb.WriteString("# Folder Analysis Report\n\n")
	sb.WriteString(fmt.Sprintf("| Item | Value |\n|------|-------|\n"))
	sb.WriteString(fmt.Sprintf("| Folder | %s |\n", r.Metadata.FolderPath))
	sb.WriteString(fmt.Sprintf("| Files | %d |\n", r.Metadata.TotalFiles))
	sb.WriteString(fmt.Sprintf("| Total Size | %s |\n", formatBytes(r.Metadata.TotalSize)))
	sb.WriteString(fmt.Sprintf("| Generated | %s |\n\n", r.Metadata.GeneratedAt.Format("2006-01-02 15:04:05")))

	sb.WriteString("## Issues\n\n")
	sb.WriteString(fmt.Sprintf("| Severity | Count |\n|----------|-------|\n"))
	sb.WriteString(fmt.Sprintf("| Critical | %d |\n", len(r.Issues.Critical)))
	sb.WriteString(fmt.Sprintf("| Error | %d |\n", len(r.Issues.Error)))
	sb.WriteString(fmt.Sprintf("| Warning | %d |\n", len(r.Issues.Warning)))
	sb.WriteString(fmt.Sprintf("| Info | %d |\n\n", len(r.Issues.Info)))

	if len(r.Issues.Critical) > 0 {
		sb.WriteString("### Critical\n\n")
		for _, issue := range r.Issues.Critical {
			sb.WriteString(fmt.Sprintf("- **%s:%d** %s\n", issue.File, issue.Line, issue.Message))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("## Quality Scores\n\n")
	sb.WriteString(fmt.Sprintf("| Category | Score |\n|----------|-------|\n"))
	sb.WriteString(fmt.Sprintf("| Naming | %d/10 |\n", r.Quality.NamingScore))
	sb.WriteString(fmt.Sprintf("| Error Handling | %d/10 |\n", r.Quality.ErrorHandlingScore))
	sb.WriteString(fmt.Sprintf("| Duplication | %d/10 |\n", r.Quality.DuplicationScore))
	sb.WriteString(fmt.Sprintf("| Security | %d/10 |\n", r.Quality.SecurityScore))
	sb.WriteString(fmt.Sprintf("| Overall | %d/10 |\n\n", r.Quality.OverallScore))

	if len(r.Suggestions) > 0 {
		sb.WriteString("## Improvement Suggestions\n\n")
		for _, s := range r.Suggestions {
			sb.WriteString(fmt.Sprintf("### [%s] %s\n\n", strings.ToUpper(s.Priority), s.Title))
			sb.WriteString(fmt.Sprintf("%s\n\n", s.Detail))
			if len(s.Files) > 0 {
				sb.WriteString(fmt.Sprintf("**Files:** %s\n\n", strings.Join(s.Files, ", ")))
			}
		}
	}

	return sb.String()
}

func formatBytes(bytes int64) string {
	unit := []string{"B", "KB", "MB", "GB"}
	i := 0
	size := float64(bytes)
	for size >= 1024 && i < len(unit)-1 {
		size /= 1024
		i++
	}
	return fmt.Sprintf("%.1f %s", size, unit[i])
}
