package stats

import (
	"testing"
	"time"
)

func makePruneSnapshot(ts time.Time, buckets map[string]int) *Snapshot {
	snap := &Snapshot{
		Timestamp: ts,
		Buckets:   make(map[string]BucketStats),
	}
	for name, keys := range buckets {
		snap.Buckets[name] = BucketStats{KeyCount: keys, Size: int64(keys * 100)}
	}
	return snap
}

func TestPruneSnapshot_Nil(t *testing.T) {
	result := PruneSnapshot(nil, DefaultPruneOptions())
	if result != nil {
		t.Error("expected nil for nil input")
	}
}

func TestPruneSnapshot_NoFilter(t *testing.T) {
	snap := makePruneSnapshot(time.Now(), map[string]int{"a": 10, "b": 5})
	opts := DefaultPruneOptions()
	opts.MinKeys = 0
	result := PruneSnapshot(snap, opts)
	if len(result.Buckets) != 2 {
		t.Errorf("expected 2 buckets, got %d", len(result.Buckets))
	}
}

func TestPruneSnapshot_MinKeys(t *testing.T) {
	snap := makePruneSnapshot(time.Now(), map[string]int{"big": 50, "small": 2, "empty": 0})
	opts := DefaultPruneOptions()
	opts.MinKeys = 10
	result := PruneSnapshot(snap, opts)
	if len(result.Buckets) != 1 {
		t.Errorf("expected 1 bucket, got %d", len(result.Buckets))
	}
	if _, ok := result.Buckets["big"]; !ok {
		t.Error("expected 'big' bucket to survive pruning")
	}
}

func TestPruneSnapshot_PreservesTimestamp(t *testing.T) {
	ts := time.Now().Add(-1 * time.Hour)
	snap := makePruneSnapshot(ts, map[string]int{"a": 5})
	result := PruneSnapshot(snap, DefaultPruneOptions())
	if !result.Timestamp.Equal(ts) {
		t.Error("timestamp should be preserved after pruning")
	}
}

func TestPruneHistory_Nil(t *testing.T) {
	result := PruneHistory(nil, DefaultPruneOptions())
	if result != nil {
		t.Error("expected nil for nil input")
	}
}

func TestPruneHistory_MaxAge(t *testing.T) {
	now := time.Now()
	snaps := []*Snapshot{
		makePruneSnapshot(now.Add(-48*time.Hour), map[string]int{"old": 1}),
		makePruneSnapshot(now.Add(-1*time.Hour), map[string]int{"recent": 1}),
	}
	opts := DefaultPruneOptions()
	opts.MaxAge = 24 * time.Hour
	opts.MaxCount = 0
	result := PruneHistory(snaps, opts)
	if len(result) != 1 {
		t.Errorf("expected 1 snapshot after age pruning, got %d", len(result))
	}
}

func TestPruneHistory_MaxCount(t *testing.T) {
	now := time.Now()
	var snaps []*Snapshot
	for i := 0; i < 10; i++ {
		snaps = append(snaps, makePruneSnapshot(now.Add(time.Duration(-i)*time.Minute), map[string]int{"b": i + 1}))
	}
	opts := PruneOptions{MaxAge: 0, MaxCount: 3}
	result := PruneHistory(snaps, opts)
	if len(result) != 3 {
		t.Errorf("expected 3 snapshots after count pruning, got %d", len(result))
	}
}
