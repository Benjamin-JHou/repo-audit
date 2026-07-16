package renderer

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	headerStyle    = lipgloss.NewStyle().Bold(true).MarginTop(1).MarginBottom(1)
	subHeaderStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00BFBF")).MarginTop(1)
	criticalStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF5555"))
	errorStyle     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF8700"))
	warningStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFD700"))
	infoStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#8787AF"))
	codeStyle      = lipgloss.NewStyle().Background(lipgloss.Color("#1a1b26")).Padding(1, 2).MarginBottom(1)
	separatorStyle = lipgloss.NewStyle().MarginTop(1).MarginBottom(1)
)

func RenderReport(reportText string) {
	lines := strings.Split(reportText, "\n")
	for _, line := range lines {
		renderLine(line)
	}
	fmt.Println()
}

func renderLine(line string) {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		fmt.Println()
		return
	}

	switch {
	case strings.HasPrefix(trimmed, "# "):
		fmt.Println(headerStyle.Render(trimmed[2:]))
	case strings.HasPrefix(trimmed, "## "):
		fmt.Println(subHeaderStyle.Render(trimmed[3:]))
	case strings.HasPrefix(trimmed, "### "):
		fmt.Println(subHeaderStyle.Render(trimmed[4:]))
	case strings.HasPrefix(trimmed, "---"):
		fmt.Println(separatorStyle.Render(strings.Repeat("-", 60)))
	case strings.HasPrefix(trimmed, "[CRITICAL]"):
		fmt.Println(criticalStyle.Render(trimmed))
	case strings.HasPrefix(trimmed, "[ERROR]"):
		fmt.Println(errorStyle.Render(trimmed))
	case strings.HasPrefix(trimmed, "[WARNING]"):
		fmt.Println(warningStyle.Render(trimmed))
	case strings.HasPrefix(trimmed, "[INFO]"):
		fmt.Println(infoStyle.Render(trimmed))
	case strings.HasPrefix(trimmed, "  ") || strings.HasPrefix(trimmed, "- "):
		fmt.Println(infoStyle.Render(trimmed))
	default:
		fmt.Println(trimmed)
	}
}

func PrintWelcome(version string) {
	welcome := fmt.Sprintf(`
%s
  ctxqa v%s - AI-Powered Folder Analyzer
  Analyze any folder with AI-powered insights.

  Usage:
    ctxqa analyze              Analyze current folder
    ctxqa analyze --dir <path> Analyze specific folder
    ctxqa analyze --dry-run    List files without analysis
    ctxqa config init          Initialize configuration

  See "ctxqa --help" for more information.
%s`,
		lipgloss.NewStyle().Bold(true).Render("ctxqa"),
		version,
		strings.Repeat("-", 60),
	)
	fmt.Println(welcome)
}

func PrintError(msg string) {
	fmt.Fprintln(os.Stderr, lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5555")).Render("Error: "+msg))
}

func PrintWarning(msg string) {
	fmt.Fprintln(os.Stderr, lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700")).Render("Warning: "+msg))
}

func PrintSuccess(msg string) {
	fmt.Fprintln(os.Stderr, lipgloss.NewStyle().Foreground(lipgloss.Color("#50FA7B")).Render("Success: "+msg))
}
