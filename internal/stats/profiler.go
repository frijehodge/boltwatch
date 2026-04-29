package stats

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// BucketProfile holds a detailed profile of a single bucket.
type BucketProfile struct {
	Name        string
	Keys        int
	Size        int64
	AvgKeySize  float64
	FillRatio   float64
	Label       string
	Score       float64
	CapturedAt  time.Time
}

// SnapshotProfile holds profiles for all buckets in a snapshot.
type SnapshotProfile struct {
	Buckets    []BucketProfile
	CapturedAt time.Time
	TotalKeys  int
	TotalSize  int64
}

// ProfileSnapshot builds a detailed profile from a snapshot.
// It combines label, score, and per-bucket statistics into a unified view.
func ProfileSnapshot(snap *Snapshot, weights ScoreWeights, cfg LabelConfig) *SnapshotProfile {
	if snap == nil || snap.IsEmpty() {
		return &SnapshotProfile{CapturedAt: time.Now()}
	}

	scores := ScoreSnapshot(snap, weights)
	labels := LabelSnapshot(snap, cfg)

	scoreMap := make(map[string]float64, len(scores))
	for _, s := range scores {
		scoreMap[s.Bucket] = s.Score
	}

	labelMap := make(map[string]string, len(labels))
	for _, l := range labels {
		labelMap[l.Bucket] = string(l.Label)
	}

	profiles := make([]BucketProfile, 0, len(snap.Buckets))
	for _, b := range snap.Buckets {
		avgKey := 0.0
		if b.Keys > 0 && b.Size > 0 {
			avgKey = float64(b.Size) / float64(b.Keys)
		}
		fill := 0.0
		if b.Size > 0 {
			fill = float64(b.Keys) / float64(b.Size) * 100.0
			if fill > 100.0 {
				fill = 100.0
			}
		}
		profiles = append(profiles, BucketProfile{
			Name:       b.Name,
			Keys:       b.Keys,
			Size:       b.Size,
			AvgKeySize: avgKey,
			FillRatio:  fill,
			Label:      labelMap[b.Name],
			Score:      scoreMap[b.Name],
			CapturedAt: snap.CapturedAt,
		})
	}

	sort.Slice(profiles, func(i, j int) bool {
		return profiles[i].Score > profiles[j].Score
	})

	return &SnapshotProfile{
		Buckets:    profiles,
		CapturedAt: snap.CapturedAt,
		TotalKeys:  snap.TotalKeys(),
		TotalSize:  snap.TotalSize(),
	}
}

// FormatProfile returns a human-readable string for a SnapshotProfile.
func FormatProfile(p *SnapshotProfile) string {
	if p == nil || len(p.Buckets) == 0 {
		return "No profile data available."
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Snapshot Profile — %s\n", p.CapturedAt.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("Total Keys: %d | Total Size: %s\n\n", p.TotalKeys, formatBytes(p.TotalSize)))
	sb.WriteString(fmt.Sprintf("%-20s %8s %10s %10s %8s %10s %8s\n",
		"Bucket", "Keys", "Size", "AvgKeyB", "Fill%", "Label", "Score"))
	sb.WriteString(strings.Repeat("-", 80) + "\n")
	for _, b := range p.Buckets {
		sb.WriteString(fmt.Sprintf("%-20s %8d %10s %10.1f %7.1f%% %10s %8.3f\n",
			b.Name, b.Keys, formatBytes(b.Size), b.AvgKeySize, b.FillRatio, b.Label, b.Score))
	}
	return sb.String()
}
