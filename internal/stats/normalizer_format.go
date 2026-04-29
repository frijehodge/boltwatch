package stats

import "fmt"

// FormatNormalized returns a human-readable summary of a normalized snapshot,
// highlighting any values that were clamped or zeroed during normalization.
func FormatNormalized(original, normalized *Snapshot) string {
	if original == nil || normalized == nil {
		return "(no data to format)\n"
	}

	var out string
	out += fmt.Sprintf("Normalized Snapshot — %s\n", normalized.Timestamp.Format("15:04:05"))
	out += fmt.Sprintf("%-20s %10s %10s %8s\n", "Bucket", "Keys", "Size", "Changed")
	out += fmt.Sprintf("%-20s %10s %10s %8s\n", "------", "----", "----", "-------")

	for name, norm := range normalized.Buckets {
		orig, exists := original.Buckets[name]
		changed := ""
		if exists {
			if orig.KeyCount != norm.KeyCount || orig.Size != norm.Size {
				changed = "yes"
			}
		}
		out += fmt.Sprintf("%-20s %10d %10s %8s\n",
			name,
			norm.KeyCount,
			formatBytes(norm.Size),
			changed,
		)
	}

	return out
}
