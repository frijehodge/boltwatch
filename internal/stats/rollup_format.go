package stats

import (
	"fmt"
	"strings"
)

// FormatRollup renders a RollupResult as a human-readable table.
func FormatRollup(result RollupResult) string {
	if len(result) == 0 {
		return "No rollup data available."
	}

	var sb strings.Builder

	// Pick any entry to read period metadata.
	var sample RollupEntry
	for _, e := range result {
		sample = e
		break
	}

	sb.WriteString(fmt.Sprintf("Rollup Period : %s\n", sample.Period))
	sb.WriteString(fmt.Sprintf("From          : %s\n", sample.PeriodStart.Format("2006-01-02 15:04:05")))
	sb.WriteString(fmt.Sprintf("To            : %s\n", sample.PeriodEnd.Format("2006-01-02 15:04:05")))
	sb.WriteString(strings.Repeat("-", 70) + "\n")
	sb.WriteString(fmt.Sprintf("%-24s %10s %10s %10s %10s\n",
		"Bucket", "AvgKeys", "MaxKeys", "AvgSize", "MaxSize"))
	sb.WriteString(strings.Repeat("-", 70) + "\n")

	for name, entry := range result {
		sb.WriteString(fmt.Sprintf("%-24s %10.1f %10d %10s %10s\n",
			truncateName(name, 24),
			entry.AvgKeys,
			entry.MaxKeys,
			formatBytes(int64(entry.AvgSize)),
			formatBytes(entry.MaxSize),
		))
	}

	sb.WriteString(strings.Repeat("-", 70) + "\n")
	sb.WriteString(fmt.Sprintf("Total buckets: %d\n", len(result)))
	return sb.String()
}

func truncateName(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}
