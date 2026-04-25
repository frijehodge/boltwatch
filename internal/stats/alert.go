package stats

import "fmt"

// AlertLevel represents the severity of a bucket alert.
type AlertLevel int

const (
	AlertNone AlertLevel = iota
	AlertWarning
	AlertCritical
)

// Alert holds information about a threshold breach for a bucket.
type Alert struct {
	Bucket  string
	Level   AlertLevel
	Message string
}

// AlertConfig defines thresholds for triggering alerts.
type AlertConfig struct {
	WarnKeyCount     int64
	CriticalKeyCount int64
	WarnSizeBytes    int64
	CriticalSizeBytes int64
}

// DefaultAlertConfig returns a sensible default alert configuration.
func DefaultAlertConfig() AlertConfig {
	return AlertConfig{
		WarnKeyCount:      10_000,
		CriticalKeyCount:  100_000,
		WarnSizeBytes:     10 * 1024 * 1024,  // 10 MB
		CriticalSizeBytes: 100 * 1024 * 1024, // 100 MB
	}
}

// CheckAlerts evaluates a snapshot against the given config and returns any alerts.
func CheckAlerts(snap *Snapshot, cfg AlertConfig) []Alert {
	if snap == nil {
		return nil
	}
	var alerts []Alert
	for _, bs := range snap.Buckets {
		keys := int64(bs.KeyN)
		size := int64(bs.LeafInuse)

		var level AlertLevel
		var msg string

		switch {
		case keys >= cfg.CriticalKeyCount:
			level = AlertCritical
			msg = fmt.Sprintf("bucket %q has %d keys (critical threshold: %d)", bs.Name, keys, cfg.CriticalKeyCount)
		case keys >= cfg.WarnKeyCount:
			level = AlertWarning
			msg = fmt.Sprintf("bucket %q has %d keys (warn threshold: %d)", bs.Name, keys, cfg.WarnKeyCount)
		case size >= cfg.CriticalSizeBytes:
			level = AlertCritical
			msg = fmt.Sprintf("bucket %q uses %s (critical threshold: %s)", bs.Name, formatBytes(size), formatBytes(cfg.CriticalSizeBytes))
		case size >= cfg.WarnSizeBytes:
			level = AlertWarning
			msg = fmt.Sprintf("bucket %q uses %s (warn threshold: %s)", bs.Name, formatBytes(size), formatBytes(cfg.WarnSizeBytes))
		}

		if level != AlertNone {
			alerts = append(alerts, Alert{Bucket: bs.Name, Level: level, Message: msg})
		}
	}
	return alerts
}
