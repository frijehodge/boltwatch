package stats

import (
	"fmt"
	"strings"
)

// FormatDeltas returns a human-readable string summarising a slice of Deltas.
// Each changed bucket is listed with its key and size differences.
func FormatDeltas(deltas []Delta) string {
	if len(deltas) == 0 {
		return "no changes detected"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-24s %10s %12s\n", "BUCKET", "KEYS DIFF", "SIZE DIFF"))
	sb.WriteString(strings.Repeat("-", 50) + "\n")

	for _, d := range deltas {
		keySign := "+"
		if d.KeysDiff < 0 {
			keySign = ""
		}
		sizeSign := "+"
		if d.SizeDiff < 0 {
			sizeSign = ""
		}
		sb.WriteString(fmt.Sprintf(
			"%-24s %s%9d %s%11s\n",
			d.Bucket,
			keySign, d.KeysDiff,
			sizeSign, formatBytes(d.SizeDiff),
		))
	}

	return sb.String()
}

// formatBytes formats an int64 byte count (possibly negative) as a
// human-readable string, e.g. "1.2 KB" or "-512 B".
func formatDeltaBytes(b int64) string {
	if b < 0 {
		return "-" + formatBytes(-b)
	}
	return formatBytes(b)
}
