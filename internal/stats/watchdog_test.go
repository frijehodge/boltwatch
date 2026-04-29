package stats

import (
	"testing"
	"time"
)

func makeWatchdogSnapshot(buckets map[string]BucketStats) *Snapshot {
	return &Snapshot{
		Timestamp: time.Now(),
		Buckets:   buckets,
	}
}

func TestEvaluateWatchdog_NilInputs(t *testing.T) {
	cfg := DefaultWatchdogConfig()
	if got := EvaluateWatchdog(nil, nil, cfg); got != nil {
		t.Errorf("expected nil, got %v", got)
	}
	snap := makeWatchdogSnapshot(map[string]BucketStats{"a": {KeyCount: 10}})
	if got := EvaluateWatchdog(nil, snap, cfg); got != nil {
		t.Error("expected nil when prev is nil")
	}
}

func TestEvaluateWatchdog_NoIssues(t *testing.T) {
	cfg := DefaultWatchdogConfig()
	prev := makeWatchdogSnapshot(map[string]BucketStats{"a": {KeyCount: 100, Size: 1024}})
	curr := makeWatchdogSnapshot(map[string]BucketStats{"a": {KeyCount: 105, Size: 1050}})
	events := EvaluateWatchdog(prev, curr, cfg)
	if len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestEvaluateWatchdog_KeyGrowth(t *testing.T) {
	cfg := DefaultWatchdogConfig()
	prev := makeWatchdogSnapshot(map[string]BucketStats{"a": {KeyCount: 100}})
	curr := makeWatchdogSnapshot(map[string]BucketStats{"a": {KeyCount: 200}})
	events := EvaluateWatchdog(prev, curr, cfg)
	if len(events) == 0 {
		t.Fatal("expected at least one event for key growth")
	}
	if events[0].Level != WatchdogWarn {
		t.Errorf("expected warn level, got %v", events[0].Level)
	}
}

func TestEvaluateWatchdog_CriticalKeys(t *testing.T) {
	cfg := DefaultWatchdogConfig()
	cfg.CriticalKeys = 50
	prev := makeWatchdogSnapshot(map[string]BucketStats{"b": {KeyCount: 10}})
	curr := makeWatchdogSnapshot(map[string]BucketStats{"b": {KeyCount: 60}})
	events := EvaluateWatchdog(prev, curr, cfg)
	found := false
	for _, e := range events {
		if e.Level == WatchdogCritical && e.Bucket == "b" {
			found = true
		}
	}
	if !found {
		t.Error("expected critical event for bucket b")
	}
}

func TestFormatWatchdogEvents_Empty(t *testing.T) {
	out := FormatWatchdogEvents(nil)
	if out == "" {
		t.Error("expected non-empty output")
	}
}

func TestFormatWatchdogEvents_WithEvents(t *testing.T) {
	events := []WatchdogEvent{
		{Bucket: "users", Level: WatchdogWarn, Reason: "key count grew 60.0%"},
		{Bucket: "logs", Level: WatchdogCritical, Reason: "size exceeds threshold"},
	}
	out := FormatWatchdogEvents(events)
	if out == "" {
		t.Error("expected non-empty output")
	}
	for _, e := range events {
		if !containsStr(out, e.Bucket) {
			t.Errorf("expected output to contain bucket %q", e.Bucket)
		}
	}
}

func TestWatchdogSummary(t *testing.T) {
	if s := WatchdogSummary(nil); s != "no issues" {
		t.Errorf("unexpected summary: %q", s)
	}
	events := []WatchdogEvent{
		{Level: WatchdogCritical},
		{Level: WatchdogWarn},
		{Level: WatchdogWarn},
	}
	s := WatchdogSummary(events)
	if !containsStr(s, "1 critical") {
		t.Errorf("expected '1 critical' in %q", s)
	}
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		}())
}
