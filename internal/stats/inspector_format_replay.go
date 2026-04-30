package stats

import (
	"fmt"
	"strings"
	"time"
)

// ReplayInspection combines a ReplayResult with per-snapshot inspection data.
type ReplayInspection struct {
	Result   *ReplayResult
	Profiles []*SnapshotProfile
}

// InspectReplay runs ProfileSnapshot on every snapshot in a ReplayResult.
func InspectReplay(r *ReplayResult) *ReplayInspection {
	if r == nil {
		return &ReplayInspection{}
	}
	profiles := make([]*SnapshotProfile, len(r.Snapshots))
	for i, snap := range r.Snapshots {
		profiles[i] = ProfileSnapshot(snap)
	}
	return &ReplayInspection{Result: r, Profiles: profiles}
}

// FormatReplayInspection renders a concise multi-step inspection report.
func FormatReplayInspection(ri *ReplayInspection) string {
	if ri == nil || ri.Result == nil || len(ri.Result.Snapshots) == 0 {
		return "replay inspection: no data\n"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Replay Inspection: %d snapshot(s)  %s → %s\n",
		len(ri.Result.Snapshots),
		ri.Result.From.Format(time.RFC3339),
		ri.Result.To.Format(time.RFC3339),
	))
	sb.WriteString(strings.Repeat("=", 72) + "\n")

	for i, p := range ri.Profiles {
		if p == nil {
			continue
		}
		sb.WriteString(fmt.Sprintf("  [%d] buckets=%d  keys=%d  size=%s  avg_keys=%.1f\n",
			i+1,
			p.BucketCount,
			p.TotalKeys,
			formatBytes(p.TotalSize),
			p.AvgKeysPerBucket,
		))
	}

	sb.WriteString(strings.Repeat("=", 72) + "\n")
	return sb.String()
}
