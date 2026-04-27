package stats

import (
	"fmt"
	"sort"
	"strings"
)

// FormatLabels returns a formatted table of labeled buckets.
func FormatLabels(labeled []LabeledBucket) string {
	if len(labeled) == 0 {
		return "No buckets to display.\n"
	}

	sort.Slice(labeled, func(i, j int) bool {
		return labeled[i].Name < labeled[j].Name
	})

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-30s %-10s %10s %12s\n", "Bucket", "Label", "Keys", "Size"))
	sb.WriteString(strings.Repeat("-", 66) + "\n")

	for _, lb := range labeled {
		sb.WriteString(fmt.Sprintf("%-30s %-10s %10d %12s\n",
			lb.Name,
			string(lb.Label),
			lb.Stats.KeyCount,
			formatBytes(lb.Stats.Size),
		))
	}

	return sb.String()
}

// LabelSummary returns a count of buckets per label.
func LabelSummary(labeled []LabeledBucket) map[BucketLabel]int {
	summary := make(map[BucketLabel]int)
	for _, lb := range labeled {
		summary[lb.Label]++
	}
	return summary
}

// FormatLabelSummary renders the label summary as a short string.
func FormatLabelSummary(labeled []LabeledBucket) string {
	if len(labeled) == 0 {
		return "No buckets.\n"
	}
	summary := LabelSummary(labeled)
	return fmt.Sprintf("large=%d  medium=%d  small=%d  empty=%d\n",
		summary[LabelLarge],
		summary[LabelMedium],
		summary[LabelSmall],
		summary[LabelEmpty],
	)
}
