package stats

import (
	"fmt"
	"strings"
)

// FormatPinBoard returns a human-readable table of all pinned snapshots.
func FormatPinBoard(pb *PinBoard) string {
	if pb == nil || pb.Len() == 0 {
		return "No pinned snapshots."
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-20s  %-26s  %8s  %12s\n",
		"LABEL", "PINNED AT", "BUCKETS", "TOTAL KEYS"))
	sb.WriteString(strings.Repeat("-", 72) + "\n")

	for _, label := range pb.Labels() {
		p := pb.Get(label)
		if p == nil {
			continue
		}
		bucketCount := len(p.Snapshot.Buckets)
		totalKeys := TotalKeys(p.Snapshot.Buckets)
		sb.WriteString(fmt.Sprintf("%-20s  %-26s  %8d  %12d\n",
			p.Label,
			p.PinnedAt.Format("2006-01-02 15:04:05 MST"),
			bucketCount,
			totalKeys,
		))
	}

	return sb.String()
}

// FormatPinDiff compares a live snapshot against a pinned one and returns a diff summary.
func FormatPinDiff(pin *PinnedSnapshot, live *Snapshot) string {
	if pin == nil || live == nil {
		return "Cannot diff: nil input."
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Pin: %s (captured %s)\n",
		pin.Label, pin.PinnedAt.Format("2006-01-02 15:04:05")))
	sb.WriteString(strings.Repeat("-", 50) + "\n")

	for name, liveStat := range live.Buckets {
		if pinStat, ok := pin.Snapshot.Buckets[name]; ok {
			delta := liveStat.KeyCount - pinStat.KeyCount
			sign := "+"
			if delta < 0 {
				sign = ""
			}
			sb.WriteString(fmt.Sprintf("  %-24s keys: %d (%s%d since pin)\n",
				name, liveStat.KeyCount, sign, delta))
		} else {
			sb.WriteString(fmt.Sprintf("  %-24s keys: %d (new since pin)\n",
				name, liveStat.KeyCount))
		}
	}

	for name := range pin.Snapshot.Buckets {
		if _, ok := live.Buckets[name]; !ok {
			sb.WriteString(fmt.Sprintf("  %-24s (removed since pin)\n", name))
		}
	}

	return sb.String()
}
