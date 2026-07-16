package report

import (
	"strings"
	"time"
)

type Generator struct{}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) ParseLLMResponse(response string) (*AuditReport, error) {
	report := &AuditReport{
		Metadata: ReportMetadata{
			GeneratedAt: time.Now(),
		},
		Overview: OverviewSection{
			FileCountByExt: make(map[string]int),
		},
	}

	report = parseSections(report, response)
	return report, nil
}

func parseSections(report *AuditReport, response string) *AuditReport {
	sections := strings.Split(response, "##")
	for _, section := range sections {
		section = strings.TrimSpace(section)
		if strings.HasPrefix(strings.ToLower(section), "repository overview") ||
			strings.HasPrefix(strings.ToLower(section), "repo overview") {
			parseOverview(report, section)
		} else if strings.HasPrefix(strings.ToLower(section), "issues") {
			parseIssues(report, section)
		} else if strings.HasPrefix(strings.ToLower(section), "quality scores") ||
			strings.HasPrefix(strings.ToLower(section), "code quality") {
			parseQuality(report, section)
		} else if strings.HasPrefix(strings.ToLower(section), "improvement") ||
			strings.HasPrefix(strings.ToLower(section), "suggestions") {
			parseSuggestions(report, section)
		} else if strings.HasPrefix(strings.ToLower(section), "completed work") ||
			strings.HasPrefix(strings.ToLower(section), "progress") {
			parseProgress(report, section)
		}
	}

	if report.Quality.OverallScore == 0 {
		report.Quality.OverallScore = calculateOverall(report)
	}

	return report
}

func parseOverview(report *AuditReport, section string) {
	lower := strings.ToLower(section)
	lines := strings.Split(section, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(lower, "project") || strings.Contains(lower, "name") || strings.Contains(lower, "folder") {
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				report.Metadata.FolderPath = strings.TrimSpace(parts[len(parts)-1])
			}
		}
		if strings.Contains(lower, "branch") {
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				report.Metadata.Branch = strings.TrimSpace(parts[len(parts)-1])
			}
		}
		if strings.Contains(lower, "file") && strings.Contains(lower, "count") {
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				val := strings.TrimSpace(parts[len(parts)-1])
				parseInt(val, &report.Metadata.TotalFiles)
			}
		}
	}
}

func parseIssues(report *AuditReport, section string) {
	lines := strings.Split(section, "\n")
	var currentSeverity string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		lower := strings.ToLower(line)

		if strings.Contains(lower, "critical") {
			currentSeverity = "critical"
			continue
		}
		if strings.Contains(lower, "error") && !strings.Contains(lower, "handling") {
			currentSeverity = "error"
			continue
		}
		if strings.Contains(lower, "warning") {
			currentSeverity = "warning"
			continue
		}
		if strings.Contains(lower, "[info]") || strings.Contains(lower, "suggestion") {
			currentSeverity = "info"
			continue
		}

		if strings.HasPrefix(line, "-") || strings.HasPrefix(line, "*") {
			line = strings.TrimPrefix(line, "-")
			line = strings.TrimPrefix(line, "*")
			line = strings.TrimSpace(line)

			issue := Issue{
				Severity: currentSeverity,
			}

			parts := strings.SplitN(line, " - ", 2)
			if len(parts) >= 1 {
				loc := strings.TrimSpace(parts[0])
				fileParts := strings.Split(loc, ":")
				if len(fileParts) >= 1 {
					issue.File = strings.Trim(fileParts[0], "`")
				}
				if len(fileParts) >= 2 {
					parseInt(fileParts[1], &issue.Line)
				}
			}
			if len(parts) >= 2 {
				issue.Message = strings.TrimSpace(parts[1])
			}

			switch currentSeverity {
			case "critical":
				report.Issues.Critical = append(report.Issues.Critical, issue)
			case "error":
				report.Issues.Error = append(report.Issues.Error, issue)
			case "warning":
				report.Issues.Warning = append(report.Issues.Warning, issue)
			default:
				report.Issues.Info = append(report.Issues.Info, issue)
			}
		}
	}
}

func parseQuality(report *AuditReport, section string) {
	lines := strings.Split(section, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(strings.ToLower(line), "overall") || strings.Contains(strings.ToLower(line), "total") {
			parts := strings.Split(line, "/")
			if len(parts) >= 1 {
				val := strings.TrimSpace(parts[0])
				parseInt(val, &report.Quality.OverallScore)
			}
		}
		if strings.Contains(strings.ToLower(line), "naming") {
			parts := strings.Split(line, "/")
			if len(parts) >= 1 {
				val := strings.TrimSpace(parts[0])
				parseInt(val, &report.Quality.NamingScore)
			}
		}
		if strings.Contains(strings.ToLower(line), "error handling") {
			parts := strings.Split(line, "/")
			if len(parts) >= 1 {
				val := strings.TrimSpace(parts[0])
				parseInt(val, &report.Quality.ErrorHandlingScore)
			}
		}
		if strings.Contains(strings.ToLower(line), "duplication") {
			parts := strings.Split(line, "/")
			if len(parts) >= 1 {
				val := strings.TrimSpace(parts[0])
				parseInt(val, &report.Quality.DuplicationScore)
			}
		}
		if strings.Contains(strings.ToLower(line), "security") {
			parts := strings.Split(line, "/")
			if len(parts) >= 1 {
				val := strings.TrimSpace(parts[0])
				parseInt(val, &report.Quality.SecurityScore)
			}
		}
	}
}

func parseSuggestions(report *AuditReport, section string) {
	lines := strings.Split(section, "\n")
	var currentSuggestion *Suggestion

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "###") || strings.HasPrefix(line, "##") {
			if currentSuggestion != nil {
				report.Suggestions = append(report.Suggestions, *currentSuggestion)
			}
			title := strings.TrimLeft(line, "#")
			title = strings.TrimSpace(title)
			currentSuggestion = &Suggestion{Title: title}
			continue
		}

		if currentSuggestion == nil {
			continue
		}

		lower := strings.ToLower(line)
		if strings.Contains(lower, "high") || strings.Contains(lower, "critical") {
			currentSuggestion.Priority = "high"
		} else if strings.Contains(lower, "medium") || strings.Contains(lower, "moderate") {
			currentSuggestion.Priority = "medium"
		} else if strings.Contains(lower, "low") || strings.Contains(lower, "minor") {
			currentSuggestion.Priority = "low"
		}

		if strings.Contains(lower, "file") && strings.Contains(lower, ":") {
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				files := strings.Split(parts[len(parts)-1], ",")
				for _, f := range files {
					f = strings.TrimSpace(strings.Trim(f, "`"))
					if f != "" {
						currentSuggestion.Files = append(currentSuggestion.Files, f)
					}
				}
			}
		}

		if line != "" && !strings.HasPrefix(line, "-") && !strings.Contains(lower, "file") {
			currentSuggestion.Detail = line
		}
	}

	if currentSuggestion != nil {
		report.Suggestions = append(report.Suggestions, *currentSuggestion)
	}
}

func parseProgress(report *AuditReport, section string) {
	lines := strings.Split(section, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "-") || strings.HasPrefix(line, "*") {
			line = strings.TrimPrefix(line, "-")
			line = strings.TrimPrefix(line, "*")
			line = strings.TrimSpace(line)
			if line != "" {
				report.Progress.Completed = append(report.Progress.Completed, Feature{
					Name: line,
					Status: "completed",
				})
			}
		}
	}
}

func calculateOverall(report *AuditReport) int {
	scores := []int{
		report.Quality.NamingScore,
		report.Quality.ErrorHandlingScore,
		report.Quality.DuplicationScore,
		report.Quality.SecurityScore,
	}
	total := 0
	for _, s := range scores {
		total += s
	}
	return total / len(scores)
}

func parseInt(s string, target *int) {
	for _, c := range s {
		if c >= '0' && c <= '9' {
			*target = *target*10 + int(c-'0')
		}
	}
}
