package stats

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// FormatMarkerBoard renders a MarkerBoard as a human-readable table.
func FormatMarkerBoard(mb *MarkerBoard) string {
	if mb == nil || mb.Len() == 0 {
		return "No markers recorded.\n"
	}

	markers := mb.All()
	sort.Slice(markers, func(i, j int) bool {
		return markers[i].Timestamp.Before(markers[j].Timestamp)
	})

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-20s %-12s %-22s %s\n", "Label", "Kind", "Timestamp", "Note"))
	sb.WriteString(strings.Repeat("-", 72) + "\n")
	for _, m := range markers {
		sb.WriteString(fmt.Sprintf("%-20s %-12s %-22s %s\n",
			truncateMarker(m.Label, 20),
			string(m.Kind),
			m.Timestamp.Format(time.RFC3339),
			m.Note,
		))
	}
	return sb.String()
}

// FormatMarkerDiff renders a comparison between a current snapshot and a marker's snapshot.
func FormatMarkerDiff(current *Snapshot, m *Marker) string {
	if m == nil {
		return "Marker not found.\n"
	}
	if current == nil {
		return "No current snapshot.\n"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Diff vs marker %q (%s):\n", m.Label, m.Timestamp.Format(time.RFC3339)))
	sb.WriteString(strings.Repeat("-", 48) + "\n")

	deltas := ComputeDeltas(m.Snapshot, current)
	if !HasChanges(deltas) {
		sb.WriteString("No changes since marker.\n")
		return sb.String()
	}
	sb.WriteString(FormatDeltas(deltas))
	return sb.String()
}

func truncateMarker(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}
