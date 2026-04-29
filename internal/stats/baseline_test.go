package stats

import (
	"testing"

	bolt "go.etcd.io/bbolt"
)

func makeBaselineSnapshot(buckets map[string]int) *Snapshot {
	stats := make(map[string]bolt.BucketStats)
	for name, keys := range buckets {
		stats[name] = bolt.BucketStats{
			KeyN:      keys,
			LeafInuse: keys * 64,
		}
	}
	return NewSnapshot(stats)
}

func TestNewBaseline_Nil(t *testing.T) {
	b := NewBaseline(nil, "test")
	if b != nil {
		t.Error("expected nil baseline for nil snapshot")
	}
}

func TestNewBaseline_SetsLabel(t *testing.T) {
	snap := makeBaselineSnapshot(map[string]int{"users": 10})
	b := NewBaseline(snap, "initial")
	if b == nil {
		t.Fatal("expected non-nil baseline")
	}
	if b.Label != "initial" {
		t.Errorf("expected label 'initial', got %q", b.Label)
	}
}

func TestBaseline_Snapshot_ReturnsOriginal(t *testing.T) {
	snap := makeBaselineSnapshot(map[string]int{"orders": 5})
	b := NewBaseline(snap, "ref")
	if b.Snapshot() != snap {
		t.Error("Snapshot() did not return the original snapshot")
	}
}

func TestComputeDrift_NilInputs(t *testing.T) {
	snap := makeBaselineSnapshot(map[string]int{"a": 1})
	b := NewBaseline(snap, "base")

	if ComputeDrift(nil, snap) != nil {
		t.Error("expected nil result for nil baseline")
	}
	if ComputeDrift(b, nil) != nil {
		t.Error("expected nil result for nil current snapshot")
	}
}

func TestComputeDrift_NoDrift(t *testing.T) {
	snap := makeBaselineSnapshot(map[string]int{"users": 10})
	b := NewBaseline(snap, "base")
	current := makeBaselineSnapshot(map[string]int{"users": 10})

	results := ComputeDrift(b, current)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].KeyDrift != 0 {
		t.Errorf("expected zero key drift, got %d", results[0].KeyDrift)
	}
}

func TestComputeDrift_NewBucket(t *testing.T) {
	base := makeBaselineSnapshot(map[string]int{"users": 5})
	b := NewBaseline(base, "base")
	current := makeBaselineSnapshot(map[string]int{"users": 5, "orders": 3})

	results := ComputeDrift(b, current)
	var found bool
	for _, r := range results {
		if r.Bucket == "orders" && r.IsNew {
			found = true
		}
	}
	if !found {
		t.Error("expected 'orders' to be marked as new bucket")
	}
}

func TestComputeDrift_RemovedBucket(t *testing.T) {
	base := makeBaselineSnapshot(map[string]int{"users": 5, "logs": 20})
	b := NewBaseline(base, "base")
	current := makeBaselineSnapshot(map[string]int{"users": 5})

	results := ComputeDrift(b, current)
	var found bool
	for _, r := range results {
		if r.Bucket == "logs" && r.IsGone {
			found = true
		}
	}
	if !found {
		t.Error("expected 'logs' to be marked as gone")
	}
}

func TestComputeDrift_KeyGrowth(t *testing.T) {
	base := makeBaselineSnapshot(map[string]int{"users": 10})
	b := NewBaseline(base, "base")
	current := makeBaselineSnapshot(map[string]int{"users": 25})

	results := ComputeDrift(b, current)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].KeyDrift != 15 {
		t.Errorf("expected key drift of 15, got %d", results[0].KeyDrift)
	}
}
