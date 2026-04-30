package stats

import (
	"fmt"
	"strings"
)

// FormatReplay returns a human-readable summary of a ReplayResult.
func FormatReplay(r *ReplayResult) string {
	if r == nil || len(r.Snapshots) == 0 {
		return "replay: no snapshots to display\n"
	}

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Replay: %d snapshot(s) from %s to %s\n",
		len(r.Snapshots),
		r.From.Format("15:04:05"),
		r.To.Format("15:04:05"),
	))
	sb.WriteString(strings.Repeat("-", 60) + "\n")

	for i, snap := range r.Snapshots {
		sb.WriteString(fmt.Sprintf("[%3d] %s  buckets=%-4d keys=%-8d size=%s\n",
			i+1,
			snap.Timestamp.Format("15:04:05.000"),
			len(snap.Buckets),
			TotalKeys(snap.Buckets),
			formatBytes(TotalSize(snap.Buckets)),
		))
	}

	sb.WriteString(strings.Repeat("-", 60) + "\n")
	return sb.String()
}

// FormatReplayDiff renders a per-step delta table across a ReplayResult.
func FormatReplayDiff(r *ReplayResult) string {
	if r == nil || len(r.Snapshots) < 2 {
		return "replay diff: need at least 2 snapshots\n"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Replay Diff: %d step(s)\n", len(r.Snapshots)-1))
	sb.WriteString(strings.Repeat("-", 60) + "\n")

	for i := 1; i < len(r.Snapshots); i++ {
		prev := r.Snapshots[i-1]
		curr := r.Snapshots[i]
		deltas := ComputeDeltas(prev, curr)
		if !HasChanges(deltas) {
			sb.WriteString(fmt.Sprintf("step %d → %s: no changes\n",
				i, curr.Timestamp.Format("15:04:05.000")))
			continue
		}
		sb.WriteString(fmt.Sprintf("step %d → %s:\n",
			i, curr.Timestamp.Format("15:04:05.000")))
		sb.WriteString(FormatDeltas(deltas))
	}

	return sb.String()
}
