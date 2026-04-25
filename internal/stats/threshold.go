package stats

import "fmt"

// ThresholdLevel represents the severity of a threshold breach.
type ThresholdLevel string

const (
	ThresholdWarning  ThresholdLevel = "WARNING"
	ThresholdCritical ThresholdLevel = "CRITICAL"
)

// ThresholdConfig defines key and size limits for a bucket.
type ThresholdConfig struct {
	Bucket        string
	MaxKeys       int
	MaxSizeBytes  int64
	Level         ThresholdLevel
}

// ThresholdViolation describes a single threshold breach.
type ThresholdViolation struct {
	Bucket  string
	Field   string
	Level   ThresholdLevel
	Message string
}

// CheckThresholds evaluates a snapshot against a set of threshold configs
// and returns any violations found.
func CheckThresholds(snap *Snapshot, configs []ThresholdConfig) []ThresholdViolation {
	if snap == nil || snap.IsEmpty() {
		return nil
	}

	var violations []ThresholdViolation

	for _, cfg := range configs {
		bs, ok := snap.Buckets[cfg.Bucket]
		if !ok {
			continue
		}

		if cfg.MaxKeys > 0 && bs.KeyN > cfg.MaxKeys {
			violations = append(violations, ThresholdViolation{
				Bucket:  cfg.Bucket,
				Field:   "keys",
				Level:   cfg.Level,
				Message: fmt.Sprintf("[%s] bucket %q has %d keys (limit: %d)", cfg.Level, cfg.Bucket, bs.KeyN, cfg.MaxKeys),
			})
		}

		if cfg.MaxSizeBytes > 0 && bs.LeafInuse > cfg.MaxSizeBytes {
			violations = append(violations, ThresholdViolation{
				Bucket:  cfg.Bucket,
				Field:   "size",
				Level:   cfg.Level,
				Message: fmt.Sprintf("[%s] bucket %q uses %s (limit: %s)", cfg.Level, cfg.Bucket, formatBytes(bs.LeafInuse), formatBytes(cfg.MaxSizeBytes)),
			})
		}
	}

	return violations
}
