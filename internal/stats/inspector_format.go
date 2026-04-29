package stats

import (
	"fmt"
	"strings"
	"time"
)

// InspectionReport holds a full diagnostic report for a snapshot,
// combining profile, labels, scores, and alerts into a single view.
type InspectionReport struct {
	Snapshot    *Snapshot
	GeneratedAt time.Time
	Profile     *SnapshotProfile
	Labels      map[string]string
	Scores      map[string]float64
	Alerts      []string
}

// InspectSnapshot builds a comprehensive InspectionReport from a snapshot.
// Returns nil if the snapshot is nil or empty.
func InspectSnapshot(snap *Snapshot, alertCfg AlertConfig) *InspectionReport {
	if snap == nil || snap.IsEmpty() {
		return nil
	}

	profile := ProfileSnapshot(snap)

	labelCfg := DefaultLabelConfig()
	labeled := LabelSnapshot(snap, labelCfg)
	labels := make(map[string]string, len(labeled))
	for _, b := range labeled {
		labels[b.Name] = b.Label
	}

	weights := DefaultScoreWeights()
	scored := ScoreSnapshot(snap, weights)
	scores := make(map[string]float64, len(scored))
	for _, b := range scored {
		scores[b.Name] = b.Score
	}

	alerts := CheckAlerts(snap, alertCfg)

	return &InspectionReport{
		Snapshot:    snap,
		GeneratedAt: time.Now(),
		Profile:     profile,
		Labels:      labels,
		Scores:      scores,
		Alerts:      alerts,
	}
}

// FormatInspectionReport renders an InspectionReport as a human-readable string.
// Includes profile summary, per-bucket labels/scores, and any active alerts.
func FormatInspectionReport(report *InspectionReport) string {
	if report == nil {
		return "(no inspection data)"
	}

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("=== BoltWatch Inspection Report ===\n"))
	sb.WriteString(fmt.Sprintf("Generated : %s\n", report.GeneratedAt.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("Snapshot  : %s\n", report.Snapshot.Timestamp.Format(time.RFC3339)))
	sb.WriteString("\n")

	if report.Profile != nil {
		sb.WriteString("--- Profile ---\n")
		sb.WriteString(FormatProfile(report.Profile))
		sb.WriteString("\n")
	}

	if len(report.Labels) > 0 {
		sb.WriteString("--- Bucket Labels & Scores ---\n")
		sb.WriteString(fmt.Sprintf("%-30s %-12s %s\n", "Bucket", "Label", "Score"))
		sb.WriteString(strings.Repeat("-", 52) + "\n")
		for _, b := range report.Snapshot.Buckets {
			label := report.Labels[b.Name]
			score := report.Scores[b.Name]
			sb.WriteString(fmt.Sprintf("%-30s %-12s %.4f\n", b.Name, label, score))
		}
		sb.WriteString("\n")
	}

	if len(report.Alerts) > 0 {
		sb.WriteString("--- Alerts ---\n")
		for _, alert := range report.Alerts {
			sb.WriteString(fmt.Sprintf("  ⚠ %s\n", alert))
		}
	} else {
		sb.WriteString("--- Alerts ---\n")
		sb.WriteString("  ✓ No alerts triggered.\n")
	}

	return sb.String()
}
