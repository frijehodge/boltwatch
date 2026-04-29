package stats

import (
	"fmt"
	"strings"
)

// levelLabel returns a short string label for a WatchdogLevel.
func levelLabel(l WatchdogLevel) string {
	switch l {
	case WatchdogCritical:
		return "CRIT"
	case WatchdogWarn:
		return "WARN"
	default:
		return "INFO"
	}
}

// levelIcon returns an emoji icon for display.
func levelIcon(l WatchdogLevel) string {
	switch l {
	case WatchdogCritical:
		return "🔴"
	case WatchdogWarn:
		return "🟡"
	default:
		return "🟢"
	}
}

// FormatWatchdogEvents renders watchdog events as a human-readable report.
func FormatWatchdogEvents(events []WatchdogEvent) string {
	if len(events) == 0 {
		return "✅ Watchdog: no issues detected."
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("⚠️  Watchdog: %d event(s) detected\n", len(events)))
	sb.WriteString(strings.Repeat("-", 60) + "\n")
	for _, e := range events {
		sb.WriteString(fmt.Sprintf("%s [%s] %-20s %s\n",
			levelIcon(e.Level),
			levelLabel(e.Level),
			e.Bucket,
			e.Reason,
		))
	}
	return sb.String()
}

// WatchdogSummary returns a one-line summary of watchdog events.
func WatchdogSummary(events []WatchdogEvent) string {
	if len(events) == 0 {
		return "no issues"
	}
	warn, crit := 0, 0
	for _, e := range events {
		switch e.Level {
		case WatchdogWarn:
			warn++
		case WatchdogCritical:
			crit++
		}
	}
	return fmt.Sprintf("%d critical, %d warning(s)", crit, warn)
}
