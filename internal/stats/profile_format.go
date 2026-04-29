package stats

import (
	"fmt"
	"strings"
)

// ProfileSummary holds aggregated highlights from a SnapshotProfile.
type ProfileSummary struct {
	TopBucket   string
	TopScore    float64
	MostKeys    string
	LargestSize string
	BucketCount int
}

// SummarizeProfile extracts a concise summary from a SnapshotProfile.
func SummarizeProfile(p *SnapshotProfile) ProfileSummary {
	if p == nil || len(p.Buckets) == 0 {
		return ProfileSummary{}
	}

	topScore := p.Buckets[0] // already sorted by score

	mostKeys := p.Buckets[0]
	largest := p.Buckets[0]
	for _, b := range p.Buckets[1:] {
		if b.Keys > mostKeys.Keys {
			mostKeys = b
		}
		if b.Size > largest.Size {
			largest = b
		}
	}

	return ProfileSummary{
		TopBucket:   topScore.Name,
		TopScore:    topScore.Score,
		MostKeys:    mostKeys.Name,
		LargestSize: largest.Name,
		BucketCount: len(p.Buckets),
	}
}

// FormatProfileSummary returns a compact summary string for a SnapshotProfile.
func FormatProfileSummary(p *SnapshotProfile) string {
	if p == nil || len(p.Buckets) == 0 {
		return "Profile summary: no data."
	}
	s := SummarizeProfile(p)
	var sb strings.Builder
	sb.WriteString("=== Profile Summary ===\n")
	sb.WriteString(fmt.Sprintf("  Buckets      : %d\n", s.BucketCount))
	sb.WriteString(fmt.Sprintf("  Top Scored   : %s (%.3f)\n", s.TopBucket, s.TopScore))
	sb.WriteString(fmt.Sprintf("  Most Keys    : %s\n", s.MostKeys))
	sb.WriteString(fmt.Sprintf("  Largest Size : %s\n", s.LargestSize))
	sb.WriteString(fmt.Sprintf("  Total Keys   : %d\n", p.TotalKeys))
	sb.WriteString(fmt.Sprintf("  Total Size   : %s\n", formatBytes(p.TotalSize)))
	return sb.String()
}
