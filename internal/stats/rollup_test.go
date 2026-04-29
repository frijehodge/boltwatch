package stats

import (
	"strings"
	"testing"
	"time"
)

func makeRollupSnapshot(t *testing.T, ts time.Time, buckets map[string]BucketStats) *Snapshot {
	t.Helper()
	snap := NewSnapshot()
	snap.Timestamp = ts
	snap.Buckets = buckets
	return snap
}

func TestRollupSnapshots_Empty(t *testing.T) {
	result := RollupSnapshots(nil, RollupMinute)
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d entries", len(result))
	}
}

func TestRollupSnapshots_SingleSnapshot(t *testing.T) {
	now := time.Now()
	snap := makeRollupSnapshot(t, now, map[string]BucketStats{
		"users": {KeyCount: 10, Size: 1024},
	})
	result := RollupSnapshots([]*Snapshot{snap}, RollupHour)
	entry, ok := result["users"]
	if !ok {
		t.Fatal("expected entry for 'users'")
	}
	if entry.SampleCount != 1 {
		t.Errorf("expected SampleCount=1, got %d", entry.SampleCount)
	}
	if entry.AvgKeys != 10 {
		t.Errorf("expected AvgKeys=10, got %.1f", entry.AvgKeys)
	}
	if entry.MaxKeys != 10 {
		t.Errorf("expected MaxKeys=10, got %d", entry.MaxKeys)
	}
	if entry.Period != RollupHour {
		t.Errorf("expected period=hour, got %s", entry.Period)
	}
}

func TestRollupSnapshots_MultipleSnapshots(t *testing.T) {
	now := time.Now()
	snaps := []*Snapshot{
		makeRollupSnapshot(t, now, map[string]BucketStats{
			"orders": {KeyCount: 20, Size: 2048},
		}),
		makeRollupSnapshot(t, now.Add(time.Minute), map[string]BucketStats{
			"orders": {KeyCount: 40, Size: 4096},
		}),
	}
	result := RollupSnapshots(snaps, RollupMinute)
	entry, ok := result["orders"]
	if !ok {
		t.Fatal("expected entry for 'orders'")
	}
	if entry.SampleCount != 2 {
		t.Errorf("expected SampleCount=2, got %d", entry.SampleCount)
	}
	if entry.AvgKeys != 30 {
		t.Errorf("expected AvgKeys=30, got %.1f", entry.AvgKeys)
	}
	if entry.MaxKeys != 40 {
		t.Errorf("expected MaxKeys=40, got %d", entry.MaxKeys)
	}
	if entry.MaxSize != 4096 {
		t.Errorf("expected MaxSize=4096, got %d", entry.MaxSize)
	}
}

func TestRollupSnapshots_NilSnapshotSkipped(t *testing.T) {
	now := time.Now()
	snaps := []*Snapshot{
		nil,
		makeRollupSnapshot(t, now, map[string]BucketStats{
			"logs": {KeyCount: 5, Size: 512},
		}),
	}
	result := RollupSnapshots(snaps, RollupDay)
	if _, ok := result["logs"]; !ok {
		t.Error("expected entry for 'logs'")
	}
}

func TestFormatRollup_Empty(t *testing.T) {
	out := FormatRollup(RollupResult{})
	if !strings.Contains(out, "No rollup data") {
		t.Errorf("expected no-data message, got: %s", out)
	}
}

func TestFormatRollup_ContainsHeaders(t *testing.T) {
	now := time.Now()
	snap := makeRollupSnapshot(t, now, map[string]BucketStats{
		"events": {KeyCount: 100, Size: 8192},
	})
	result := RollupSnapshots([]*Snapshot{snap}, RollupHour)
	out := FormatRollup(result)
	for _, want := range []string{"Bucket", "AvgKeys", "MaxKeys", "AvgSize", "MaxSize", "events"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q\n%s", want, out)
		}
	}
}
