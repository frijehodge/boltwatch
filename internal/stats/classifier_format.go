package stats

import (
	"fmt"
	"sort"
	"strings"
)

// FormatClassifications returns a human-readable table of bucket growth classifications.
func FormatClassifications(classifications map[string]BucketClassification) string {
	if len(classifications) == 0 {
		return "No bucket classifications available."
	}

	names := make([]string, 0, len(classifications))
	for name := range classifications {
		names = append(names, name)
	}
	sort.Strings(names)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-30s %-12s %10s %12s %14s\n",
		"Bucket", "Class", "KeyDelta", "SizeDelta", "GrowthRate/s"))
	sb.WriteString(strings.Repeat("-", 82) + "\n")

	for _, name := range names {
		c := classifications[name]
		classStr := classIcon(c.Class) + " " + string(c.Class)
		sb.WriteString(fmt.Sprintf("%-30s %-12s %+10d %12s %+14.2f\n",
			truncate(name, 30),
			classStr,
			c.KeyDelta,
			formatDeltaBytes(c.SizeDelta),
			c.GrowthRate,
		))
	}

	return sb.String()
}

func classIcon(c GrowthClass) string {
	switch c {
	case GrowthClassGrowing:
		return "↑"
	case GrowthClassShrinking:
		return "↓"
	case GrowthClassNew:
		return "+"
	case GrowthClassRemoved:
		return "✗"
	default:
		return "~"
	}
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}
