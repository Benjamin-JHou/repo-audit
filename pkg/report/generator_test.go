package report

import (
	"testing"
	"time"
)

func TestFormatText(t *testing.T) {
	report := &AuditReport{
		Metadata: ReportMetadata{
			FolderPath:  "/tmp/test",
			TotalFiles:  10,
			GeneratedAt: time.Now(),
		},
		Quality: QualitySection{
			NamingScore:       8,
			ErrorHandlingScore: 6,
			DuplicationScore:  7,
			SecurityScore:     9,
			OverallScore:      8,
		},
	}

	output := report.FormatText("0.1.0")
	if output == "" {
		t.Error("expected non-empty output")
	}
	if len(output) == 0 {
		t.Error("output should not be empty")
	}
}

func TestFormatJSON(t *testing.T) {
	report := &AuditReport{
		Metadata: ReportMetadata{
			FolderPath:  "/tmp/test",
			TotalFiles:  5,
			GeneratedAt: time.Now(),
		},
	}

	output := report.FormatJSON()
	if output == "" {
		t.Error("expected non-empty JSON output")
	}
}

func TestFormatMarkdown(t *testing.T) {
	report := AuditReport{
		Metadata: ReportMetadata{
			FolderPath:  "/tmp/test",
			TotalFiles:  5,
			GeneratedAt: time.Now(),
		},
	}

	output := report.FormatMarkdown()
	if output == "" {
		t.Error("expected non-empty markdown output")
	}
}

func TestParseIssues(t *testing.T) {
	report := &AuditReport{}
	parseIssues(report, `## Issues

### Critical
- main.go:42 - Missing error handling

### Warning
- utils.go:15 - Unused variable
`)

	if len(report.Issues.Critical) != 1 {
		t.Errorf("expected 1 critical issue, got %d", len(report.Issues.Critical))
	}
	if len(report.Issues.Warning) != 1 {
		t.Errorf("expected 1 warning, got %d", len(report.Issues.Warning))
	}
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		bytes    int64
		expected string
	}{
		{500, "500.0 B"},
		{1024, "1.0 KB"},
		{1048576, "1.0 MB"},
		{1073741824, "1.0 GB"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := formatBytes(tt.bytes)
			if result != tt.expected {
				t.Errorf("formatBytes(%d) = %s, want %s", tt.bytes, result, tt.expected)
			}
		})
	}
}
