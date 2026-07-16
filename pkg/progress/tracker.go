package progress

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

type Tracker struct {
	totalBatches  int
	totalFiles    int
	currentBatch  int
	completedFiles int
	startTime     time.Time
	style         lipgloss.Style
}

func New(totalBatches, totalFiles int) *Tracker {
	return &Tracker{
		totalBatches: totalBatches,
		totalFiles:   totalFiles,
		startTime:    time.Now(),
		style: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00BFBF")),
	}
}

func (t *Tracker) UpdateBatch(batch int) {
	t.currentBatch = batch
	t.printProgress()
}

func (t *Tracker) IncrementFiles(count int) {
	t.completedFiles += count
	t.printProgress()
}

func (t *Tracker) printProgress() {
	if t.totalBatches == 0 {
		return
	}

	percent := float64(t.currentBatch) / float64(t.totalBatches) * 100
	barWidth := 30
	filled := int(percent / 100.0 * float64(barWidth))
	if filled > barWidth {
		filled = barWidth
	}

	bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)
	elapsed := time.Since(t.startTime).Round(time.Second)

	fmt.Fprintf(os.Stderr, "\r%s [%s] %5.1f%% (%d/%d files) %s   ",
		t.style.Render("ctxqa"),
		bar,
		percent,
		t.completedFiles,
		t.totalFiles,
		elapsed,
	)
}

func (t *Tracker) Done() {
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintf(os.Stderr, "\r"+strings.Repeat(" ", 120)+"\r")
}
