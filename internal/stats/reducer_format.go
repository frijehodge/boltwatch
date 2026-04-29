package stats

import (
	"fmt"
	"strings"
)

// FormatReduced returns a human-readable table of a reduced snapshot.
// It includes a header indicating how many buckets are shown vs total.
func FormatReduced(snap *Snapshot, originalCount int) string {
	if snap == nil {
		return "(no snapshot)"
	}

	var sb strings.Builder

	shown := len(snap.Buckets)
	if originalCount > shown {
		fmt.Fprintf(&sb, "Showing top %d of %d buckets (reduced)\n", shown, originalCount)
	} else {
		fmt.Fprintf(&sb, "Showing all %d buckets\n", shown)
	}

	if shown == 0 {
		sb.WriteString("  (empty)\n")
		return sb.String()
	}

	fmt.Fprintf(&sb, "%-30s %10s %12s\n", "Bucket", "Keys", "Size")
	sb.WriteString(strings.Repeat("-", 56) + "\n")

	for _, b := range snap.Buckets {
		fmt.Fprintf(&sb, "%-30s %10d %12s\n", b.Name, b.KeyCount, formatBytes(b.Size))
	}

	return sb.String()
}
