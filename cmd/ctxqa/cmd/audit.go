package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/ctxqa/ctxqa/pkg/api"
	"github.com/ctxqa/ctxqa/pkg/config"
	"github.com/ctxqa/ctxqa/pkg/context"
	"github.com/ctxqa/ctxqa/pkg/filesystem"
	"github.com/ctxqa/ctxqa/pkg/progress"
	"github.com/ctxqa/ctxqa/pkg/renderer"
	"github.com/ctxqa/ctxqa/pkg/report"
	"github.com/spf13/cobra"
)

const version = "0.1.0"

var analyzeCmd = &cobra.Command{
	Use:   "analyze [flags]",
	Short: "Analyze the folder",
	Long:  "Perform a comprehensive analysis of any folder by reviewing all files and generating a structured report with organization suggestions, issue detection, quality scoring, and improvement recommendations.",
	RunE:  runAnalyze,
}

var (
	format        string
	output        string
	severity      string
	section       string
	incremental   bool
	resume        bool
	diffMode      bool
	noContext     bool
	dir           string
	concurrency   int
	dryRun        bool
	summary       bool
)

func init() {
	analyzeCmd.Flags().StringVarP(&format, "format", "f", "text", "Output format: text, json, markdown")
	analyzeCmd.Flags().StringVarP(&output, "output", "o", "", "Output file path")
	analyzeCmd.Flags().StringVar(&severity, "severity", "all", "Minimum severity: all, info, warning, error, critical")
	analyzeCmd.Flags().StringVar(&section, "section", "all", "Output specific section: overview, structure, progress, issues, suggestions")
	analyzeCmd.Flags().BoolVar(&incremental, "incremental", false, "Incremental analysis mode")
	analyzeCmd.Flags().BoolVar(&resume, "resume", false, "Resume from last interrupted analysis")
	analyzeCmd.Flags().BoolVar(&diffMode, "diff", false, "Compare with last analysis result")
	analyzeCmd.Flags().BoolVar(&noContext, "no-context", false, "Skip context collection")
	analyzeCmd.Flags().StringVarP(&dir, "dir", "d", ".", "Folder to analyze")
	analyzeCmd.Flags().IntVarP(&concurrency, "concurrency", "j", 1, "Parallel batch processors")
	analyzeCmd.Flags().BoolVar(&dryRun, "dry-run", false, "List files only, do not analyze")
	analyzeCmd.Flags().BoolVar(&summary, "summary", false, "Brief summary only")

	rootCmd.AddCommand(analyzeCmd)
}

func runAnalyze(cmd *cobra.Command, args []string) error {
	renderer.PrintWelcome(version)

	cfg, err := config.Load("")
	if err != nil {
		renderer.PrintError("Failed to load config: " + err.Error())
		os.Exit(1)
	}

	if cfg.APIKey == "" {
		fmt.Println("API Key not configured. Run 'ctxqa config init' to set up.")
		os.Exit(1)
	}

	if dir == "" {
		dir = "."
	}

	walker := filesystem.NewFileScanner(cfg)
	files, err := walker.Scan(dir)
	if err != nil {
		renderer.PrintError("Failed to walk directory: " + err.Error())
		os.Exit(1)
	}

	if len(files) == 0 {
		fmt.Println("No code files found in the specified directory.")
		return nil
	}

	if dryRun {
		fmt.Printf("Found %d files:\n", len(files))
		for _, f := range files {
			fmt.Printf("  %s (%s, %d bytes)\n", f.Path, f.FileType, f.Size)
		}
		return nil
	}

	batcher := filesystem.NewBatcher(cfg.MaxBatchChars, cfg.MaxFileSize)
	batches := batcher.Split(files)

	tracker := progress.New(len(batches), len(files))

	var allReports []*report.AuditReport
	var ctxStr string

	if !noContext {
		ctx, err := context.Collect(cfg, dir, incremental)
		if err == nil {
			ctxStr = buildContextString(ctx)
		}
	}

	systemPrompt := `You are a professional code auditor. Analyze the provided code files and generate a structured audit report.`

	for i, batch := range batches {
		tracker.UpdateBatch(i + 1)

		userPrompt := batch.BuildPrompt(ctxStr)
		messages := []api.Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		}

		client := api.New(cfg.Provider, cfg.APIKey, cfg.BaseURL, cfg.Model, cfg.TimeoutSeconds, cfg.RetryCount)
		response, err := client.Chat(messages)
		if err != nil {
			renderer.PrintError("API call failed for batch " + fmt.Sprintf("%d", i+1) + ": " + err.Error())
			continue
		}

		gen := report.NewGenerator()
		batchReport, err := gen.ParseLLMResponse(response)
		if err == nil {
			allReports = append(allReports, batchReport)
		}

		tracker.IncrementFiles(len(batch.Files))
	}

	tracker.Done()

	finalReport := mergeReports(allReports, len(files), dir, cfg)
	formatted := formatReport(finalReport, format, version)

	if format == "text" {
		renderer.RenderReport(formatted)
	} else {
		if output == "-" {
			fmt.Print(formatted)
		} else if output != "" {
			if err := os.WriteFile(output, []byte(formatted), 0644); err != nil {
				renderer.PrintError("Failed to write output file: " + err.Error())
				os.Exit(1)
			}
			renderer.PrintSuccess("Report saved to " + output)
		} else {
			fmt.Print(formatted)
		}
	}

	return nil
}

func buildContextString(ctx *context.Context) string {
	var sb strings.Builder
	if ctx.Branch != "" && ctx.Branch != "unavailable" {
		sb.WriteString(fmt.Sprintf("Branch: %s\n", ctx.Branch))
	}
	if len(ctx.Commits) > 0 {
		sb.WriteString("\nRecent commits:\n")
		for _, c := range ctx.Commits {
			sb.WriteString(fmt.Sprintf("  - %s\n", c))
		}
	}
	if len(ctx.History) > 0 {
		sb.WriteString("\nRecent commands:\n")
		for _, h := range ctx.History {
			sb.WriteString(fmt.Sprintf("  $ %s\n", h))
		}
	}
	return sb.String()
}

func mergeReports(reports []*report.AuditReport, totalFiles int, dir string, cfg *config.Config) *report.AuditReport {
	if len(reports) == 0 {
		return &report.AuditReport{}
	}

	merged := &report.AuditReport{
		Metadata: report.ReportMetadata{
			FolderPath:        dir,
			TotalFiles:        totalFiles,
			GeneratedAt:       reports[0].Metadata.GeneratedAt,
			TechStack:         []string{},
			FileTypeBreakdown: map[string]int{},
		},
	}

	if len(reports) > 0 {
		merged.Metadata.Branch = reports[0].Metadata.Branch
		merged.Metadata.TechStack = reports[0].Metadata.TechStack
		merged.Metadata.FileTypeBreakdown = reports[0].Metadata.FileTypeBreakdown
		merged.Metadata.TotalSize = reports[0].Metadata.TotalSize
	}

	for _, r := range reports {
		merged.Issues.Critical = append(merged.Issues.Critical, r.Issues.Critical...)
		merged.Issues.Error = append(merged.Issues.Error, r.Issues.Error...)
		merged.Issues.Warning = append(merged.Issues.Warning, r.Issues.Warning...)
		merged.Issues.Info = append(merged.Issues.Info, r.Issues.Info...)
		merged.Suggestions = append(merged.Suggestions, r.Suggestions...)

		for ext, count := range r.Metadata.FileTypeBreakdown {
			merged.Metadata.FileTypeBreakdown[ext] += count
		}
	}

	if merged.Quality.OverallScore == 0 {
		scores := []int{
			merged.Quality.NamingScore,
			merged.Quality.ErrorHandlingScore,
			merged.Quality.DuplicationScore,
			merged.Quality.SecurityScore,
		}
		total := 0
		for _, s := range scores {
			total += s
		}
		if len(scores) > 0 {
			merged.Quality.OverallScore = total / len(scores)
		}
	}

	return merged
}

func formatReport(report *report.AuditReport, format string, version string) string {
	switch format {
	case "json":
		return report.FormatJSON()
	case "markdown":
		return report.FormatMarkdown()
	default:
		return report.FormatText(version)
	}
}
