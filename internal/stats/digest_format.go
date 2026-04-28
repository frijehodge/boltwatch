package stats

import (
	"fmt"
	"strings"
)

// FormatDigest returns a human-readable summary of a DigestResult.
func FormatDigest(d *DigestResult) string {
	if d == nil {
		return "(no digest)"
	}

	var sb strings.Builder
	sb.WriteString("Snapshot Digest\n")
	sb.WriteString(strings.Repeat("-", 36) + "\n")
	fmt.Fprintf(&sb, "  Hash:    %s\n", d.Hash)
	fmt.Fprintf(&sb, "  Buckets: %d\n", d.BucketCount)
	fmt.Fprintf(&sb, "  Keys:    %d\n", d.KeyCount)
	fmt.Fprintf(&sb, "  Size:    %s\n", formatBytes(d.ByteSize))
	return sb.String()
}

// FormatDigestDiff returns a short message describing whether two digests match.
func FormatDigestDiff(prev, curr *DigestResult) string {
	if DigestsMatch(prev, curr) {
		return "no changes detected"
	}
	if prev == nil {
		return "initial snapshot captured"
	}

	var parts []string

	keyDelta := curr.KeyCount - prev.KeyCount
	if keyDelta != 0 {
		sign := "+"
		if keyDelta < 0 {
			sign = ""
		}
		parts = append(parts, fmt.Sprintf("keys %s%d", sign, keyDelta))
	}

	bucketDelta := curr.BucketCount - prev.BucketCount
	if bucketDelta != 0 {
		sign := "+"
		if bucketDelta < 0 {
			sign = ""
		}
		parts = append(parts, fmt.Sprintf("buckets %s%d", sign, bucketDelta))
	}

	if len(parts) == 0 {
		return fmt.Sprintf("hash changed (%s → %s)", prev.Hash[:8], curr.Hash[:8])
	}
	return fmt.Sprintf("changed: %s", strings.Join(parts, ", "))
}
