package stats

import (
	"fmt"
	"strings"
)

// FormatDiffSummary returns a compact one-line summary of a DiffResult.
func FormatDiffSummary(d *DiffResult) string {
	if d == nil || d.IsEmpty() {
		return "buckets: no changes"
	}
	parts := []string{}
	if n := len(d.Added); n > 0 {
		parts = append(parts, fmt.Sprintf("%d added", n))
	}
	if n := len(d.Removed); n > 0 {
		parts = append(parts, fmt.Sprintf("%d removed", n))
	}
	if n := len(d.Changed); n > 0 {
		parts = append(parts, fmt.Sprintf("%d changed", n))
	}
	return "buckets: " + strings.Join(parts, ", ")
}

// FormatDiffTable returns a formatted table-style string for a DiffResult,
// listing each bucket with its change status.
func FormatDiffTable(d *DiffResult) string {
	if d == nil || d.IsEmpty() {
		return "No bucket differences detected.\n"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-30s %s\n", "BUCKET", "STATUS"))
	sb.WriteString(strings.Repeat("-", 42) + "\n")

	for _, name := range d.Added {
		sb.WriteString(fmt.Sprintf("%-30s %s\n", name, "ADDED"))
	}
	for _, name := range d.Removed {
		sb.WriteString(fmt.Sprintf("%-30s %s\n", name, "REMOVED"))
	}
	for _, name := range d.Changed {
		sb.WriteString(fmt.Sprintf("%-30s %s\n", name, "CHANGED"))
	}
	for _, name := range d.Unchanged {
		sb.WriteString(fmt.Sprintf("%-30s %s\n", name, "UNCHANGED"))
	}
	return sb.String()
}
