package stats

import (
	"strings"
	"testing"
	"time"
)

// TestWatchdog_FullPipeline exercises EvaluateWatchdog → FormatWatchdogEvents together.
func TestWatchdog_FullPipeline_NoIssues(t *testing.T) {
	cfg := DefaultWatchdogConfig()
	prev := &Snapshot{
		Timestamp: time.Now().Add(-5 * time.Second),
		Buckets:   map[string]BucketStats{"orders": {KeyCount: 200, Size: 2048}},
	}
	curr := &Snapshot{
		Timestamp: time.Now(),
		Buckets:   map[string]BucketStats{"orders": {KeyCount: 210, Size: 2100}},
	}
	events := EvaluateWatchdog(prev, curr, cfg)
	out := FormatWatchdogEvents(events)
	if !strings.Contains(out, "no issues") {
		t.Errorf("expected no-issues message, got: %s", out)
	}
}

func TestWatchdog_FullPipeline_WithIssues(t *testing.T) {
	cfg := DefaultWatchdogConfig()
	cfg.CriticalKeys = 10
	prev := &Snapshot{
		Timestamp: time.Now().Add(-5 * time.Second),
		Buckets:   map[string]BucketStats{"sessions": {KeyCount: 5, Size: 512}},
	}
	curr := &Snapshot{
		Timestamp: time.Now(),
		Buckets:   map[string]BucketStats{"sessions": {KeyCount: 50, Size: 5000}},
	}
	events := EvaluateWatchdog(prev, curr, cfg)
	if len(events) == 0 {
		t.Fatal("expected events but got none")
	}
	out := FormatWatchdogEvents(events)
	if !strings.Contains(out, "sessions") {
		t.Errorf("expected 'sessions' in output, got: %s", out)
	}
	summary := WatchdogSummary(events)
	if summary == "no issues" {
		t.Error("expected non-trivial summary")
	}
}

func TestWatchdog_DefaultConfig_Thresholds(t *testing.T) {
	cfg := DefaultWatchdogConfig()
	if cfg.MaxKeyGrowthPct <= 0 {
		t.Error("MaxKeyGrowthPct should be positive")
	}
	if cfg.MaxSizeGrowthPct <= 0 {
		t.Error("MaxSizeGrowthPct should be positive")
	}
	if cfg.CriticalKeys <= 0 {
		t.Error("CriticalKeys should be positive")
	}
	if cfg.CriticalSize <= 0 {
		t.Error("CriticalSize should be positive")
	}
}
