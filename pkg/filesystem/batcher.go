package filesystem

import (
	"fmt"
	"strings"
)

type Batch struct {
	Index      int           `json:"index"`
	Total      int           `json:"total"`
	Files      []FileContent `json:"files"`
	TotalChars int           `json:"total_chars"`
	Status     BatchStatus   `json:"status"`
}

type BatchStatus string

const (
	StatusPending   BatchStatus = "pending"
	StatusProcessing BatchStatus = "processing"
	StatusCompleted  BatchStatus = "completed"
	StatusFailed     BatchStatus = "failed"
)

type Batcher struct {
	maxBatchChars int
	maxFileSize   int
}

func NewBatcher(maxBatchChars, maxFileSize int) *Batcher {
	return &Batcher{
		maxBatchChars: maxBatchChars,
		maxFileSize:   maxFileSize,
	}
}

func (b *Batcher) Split(files []FileContent) []*Batch {
	if len(files) == 0 {
		return nil
	}

	var batches []*Batch
	var currentBatch *Batch
	var currentChars int

	for i, file := range files {
		fileChars := len([]rune(file.Content))

		if currentBatch == nil {
			currentBatch = &Batch{
				Index:  len(batches),
				Files:  []FileContent{},
				Status: StatusPending,
			}
		}

		if fileChars > b.maxBatchChars && len(currentBatch.Files) > 0 {
			batches = append(batches, currentBatch)
			currentBatch = &Batch{
				Index:  len(batches),
				Files:  []FileContent{},
				Status: StatusPending,
			}
			currentChars = 0
		}

		if currentChars+fileChars > b.maxBatchChars && len(currentBatch.Files) > 0 {
			batches = append(batches, currentBatch)
			currentBatch = &Batch{
				Index:  len(batches),
				Files:  []FileContent{},
				Status: StatusPending,
			}
			currentChars = 0
		}

		currentBatch.Files = append(currentBatch.Files, file)
		currentBatch.Total = len(files)
		currentChars += fileChars
		currentBatch.TotalChars = currentChars

		if i == len(files)-1 {
			batches = append(batches, currentBatch)
		}
	}

	for _, batch := range batches {
		batch.Total = len(files)
	}

	return batches
}

func (b *Batch) BuildPrompt(contextStr string) string {
	var sb strings.Builder

	sb.WriteString("You are a professional code and file organization auditor. Review the following files and provide a comprehensive analysis report.\n\n")

	if contextStr != "" {
		sb.WriteString("## Project Context\n")
		sb.WriteString(contextStr)
		sb.WriteString("\n\n")
	}

	sb.WriteString("## Files to Analyze\n\n")

	for _, file := range b.Files {
		sb.WriteString(fmt.Sprintf("### File: %s (type: %s, size: %d bytes)\n\n", file.Path, file.FileType, file.Size))
		sb.WriteString("```\n")
		sb.WriteString(file.Content)
		sb.WriteString("\n```\n\n")
	}

	sb.WriteString("Please provide a comprehensive analysis report covering:\n")
	sb.WriteString("1. Folder overview (total files, total size, file type breakdown, top directories)\n")
	sb.WriteString("2. File structure and organization suggestions (current layout summary, proposed reorganization, file grouping recommendations, naming issues, duplicate files)\n")
	sb.WriteString("3. Progress analysis (completed work, work in progress, pending tasks based on content and commit history)\n")
	sb.WriteString("4. Issues found (classified by severity: critical, error, warning, info)\n")
	sb.WriteString("5. Quality scores (naming conventions, error handling, code duplication, security) with scores 0-10\n")
	sb.WriteString("6. Improvement suggestions (prioritized: high, medium, low)\n")
	sb.WriteString("\n")
	sb.WriteString("Format your response as structured sections with clear headings.\n")

	return sb.String()
}
