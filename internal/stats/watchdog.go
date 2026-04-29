package stats

import "fmt"

// WatchdogLevel represents the severity of a watchdog event.
type WatchdogLevel int

const (
	WatchdogInfo WatchdogLevel = iota
	WatchdogWarn
	WatchdogCritical
)

// WatchdogEvent describes a condition detected during snapshot monitoring.
type WatchdogEvent struct {
	Bucket string
	Level  WatchdogLevel
	Reason string
}

// WatchdogConfig controls thresholds for watchdog evaluation.
type WatchdogConfig struct {
	MaxKeyGrowthPct  float64 // percent growth in keys between snapshots
	MaxSizeGrowthPct float64 // percent growth in size between snapshots
	CriticalKeys     int64
	CriticalSize     int64
}

// DefaultWatchdogConfig returns sensible defaults.
func DefaultWatchdogConfig() WatchdogConfig {
	return WatchdogConfig{
		MaxKeyGrowthPct:  50.0,
		MaxSizeGrowthPct: 50.0,
		CriticalKeys:     100_000,
		CriticalSize:     500 * 1024 * 1024,
	}
}

// EvaluateWatchdog compares two snapshots and emits watchdog events.
func EvaluateWatchdog(prev, curr *Snapshot, cfg WatchdogConfig) []WatchdogEvent {
	if prev == nil || curr == nil {
		return nil
	}
	var events []WatchdogEvent
	for name, cs := range curr.Buckets {
		ps, ok := prev.Buckets[name]
		if ok && ps.KeyCount > 0 {
			growth := float64(cs.KeyCount-ps.KeyCount) / float64(ps.KeyCount) * 100
			if growth >= cfg.MaxKeyGrowthPct {
				events = append(events, WatchdogEvent{
					Bucket: name,
					Level:  WatchdogWarn,
					Reason: fmt.Sprintf("key count grew %.1f%% (threshold %.1f%%)", growth, cfg.MaxKeyGrowthPct),
				})
			}
			sizeGrowth := float64(cs.Size-ps.Size) / float64(ps.Size) * 100
			if ps.Size > 0 && sizeGrowth >= cfg.MaxSizeGrowthPct {
				events = append(events, WatchdogEvent{
					Bucket: name,
					Level:  WatchdogWarn,
					Reason: fmt.Sprintf("size grew %.1f%% (threshold %.1f%%)", sizeGrowth, cfg.MaxSizeGrowthPct),
				})
			}
		}
		if cs.KeyCount >= cfg.CriticalKeys {
			events = append(events, WatchdogEvent{
				Bucket: name,
				Level:  WatchdogCritical,
				Reason: fmt.Sprintf("key count %d exceeds critical threshold %d", cs.KeyCount, cfg.CriticalKeys),
			})
		}
		if cs.Size >= cfg.CriticalSize {
			events = append(events, WatchdogEvent{
				Bucket: name,
				Level:  WatchdogCritical,
				Reason: fmt.Sprintf("size %s exceeds critical threshold %s", formatBytes(cs.Size), formatBytes(cfg.CriticalSize)),
			})
		}
	}
	return events
}
