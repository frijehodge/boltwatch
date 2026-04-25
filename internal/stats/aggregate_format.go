package stats

import (
	"fmt"
	"strings"
)

// FormatAggregate renders an AggregateResult as a human-readable summary table.
func FormatAggregate(result *AggregateResult, topN int) string {
	if result == nil || result.Buckets == 0 {
		return "No bucket data available.\n"
	}

	var sb strings.Builder

	fmt.Fprintf(&sb, "Buckets: %d | Total Keys: %d | Total Size: %s\n",
		result.Buckets,
		result.TotalKeys,
		formatBytes(result.TotalSize),
	)

	limit := topN
	if limit <= 0 || limit > len(result.TopByKeys) {
		limit = len(result.TopByKeys)
	}

	sb.WriteString("\nTop Buckets by Keys:\n")
	sb.WriteString(fmt.Sprintf("  %-30s %10s\n", "Bucket", "Keys"))
	sb.WriteString("  " + strings.Repeat("-", 42) + "\n")
	for i := 0; i < limit; i++ {
		r := result.TopByKeys[i]
		fmt.Fprintf(&sb, "  %-30s %10d\n", r.Bucket, r.Keys)
	}

	limitSize := topN
	if limitSize <= 0 || limitSize > len(result.TopBySize) {
		limitSize = len(result.TopBySize)
	}

	sb.WriteString("\nTop Buckets by Size:\n")
	sb.WriteString(fmt.Sprintf("  %-30s %10s\n", "Bucket", "Size"))
	sb.WriteString("  " + strings.Repeat("-", 42) + "\n")
	for i := 0; i < limitSize; i++ {
		r := result.TopBySize[i]
		fmt.Fprintf(&sb, "  %-30s %10s\n", r.Bucket, formatBytes(r.Size))
	}

	return sb.String()
}
