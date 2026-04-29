package stats

import (
	"testing"
	"time"
)

func makeMergeSnapshot(buckets map[string]int) *Snapshot {
	stats := make(map[string]BucketStats, len(buckets))
	for name, keys := range buckets {
		stats[name] = BucketStats{Name: name, KeyCount: keys, Size: int64(keys * 64)}
	}
	return NewSnapshot(stats)
}

func TestMergeSnapshots_BothNil(t *testing.T) {
	result := MergeSnapshots(nil, nil, MergeStrategyPreferA)
	if result != nil {
		t.Fatal("expected nil result for two nil inputs")
	}
}

func TestMergeSnapshots_NilA(t *testing.T) {
	b := makeMergeSnapshot(map[string]int{"alpha": 5})
	result := MergeSnapshots(nil, b, MergeStrategyPreferA)
	if result == nil || result.Snapshot != b {
		t.Fatal("expected snapshot B when A is nil")
	}
}

func TestMergeSnapshots_NilB(t *testing.T) {
	a := makeMergeSnapshot(map[string]int{"alpha": 5})
	result := MergeSnapshots(a, nil, MergeStrategyPreferA)
	if result == nil || result.Snapshot != a {
		t.Fatal("expected snapshot A when B is nil")
	}
}

func TestMergeSnapshots_NoConflicts(t *testing.T) {
	a := makeMergeSnapshot(map[string]int{"alpha": 3})
	b := makeMergeSnapshot(map[string]int{"beta": 7})
	result := MergeSnapshots(a, b, MergeStrategyPreferA)
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if len(result.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %v", result.Conflicts)
	}
	if len(result.Snapshot.Buckets) != 2 {
		t.Errorf("expected 2 buckets, got %d", len(result.Snapshot.Buckets))
	}
}

func TestMergeSnapshots_PreferA(t *testing.T) {
	a := makeMergeSnapshot(map[string]int{"shared": 10})
	b := makeMergeSnapshot(map[string]int{"shared": 99})
	result := MergeSnapshots(a, b, MergeStrategyPreferA)
	if result.Snapshot.Buckets["shared"].KeyCount != 10 {
		t.Errorf("expected 10 keys from A, got %d", result.Snapshot.Buckets["shared"].KeyCount)
	}
	if len(result.Conflicts) != 1 || result.Conflicts[0] != "shared" {
		t.Errorf("expected conflict on 'shared', got %v", result.Conflicts)
	}
}

func TestMergeSnapshots_PreferB(t *testing.T) {
	a := makeMergeSnapshot(map[string]int{"shared": 10})
	b := makeMergeSnapshot(map[string]int{"shared": 99})
	result := MergeSnapshots(a, b, MergeStrategyPreferB)
	if result.Snapshot.Buckets["shared"].KeyCount != 99 {
		t.Errorf("expected 99 keys from B, got %d", result.Snapshot.Buckets["shared"].KeyCount)
	}
}

func TestMergeSnapshots_SumKeys(t *testing.T) {
	a := makeMergeSnapshot(map[string]int{"shared": 10})
	b := makeMergeSnapshot(map[string]int{"shared": 5})
	result := MergeSnapshots(a, b, MergeStrategySumKeys)
	if result.Snapshot.Buckets["shared"].KeyCount != 15 {
		t.Errorf("expected 15 summed keys, got %d", result.Snapshot.Buckets["shared"].KeyCount)
	}
}

func TestMergeSnapshots_UsesLatestTimestamp(t *testing.T) {
	a := makeMergeSnapshot(map[string]int{"a": 1})
	b := makeMergeSnapshot(map[string]int{"b": 2})
	b.Timestamp = time.Now().Add(10 * time.Second)
	result := MergeSnapshots(a, b, MergeStrategyPreferA)
	if !result.Snapshot.Timestamp.Equal(b.Timestamp) {
		t.Errorf("expected timestamp from B (newer), got %v", result.Snapshot.Timestamp)
	}
}
