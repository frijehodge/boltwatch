package stats

import (
	"testing"

	"go.etcd.io/bbolt"
)

func makeAggregateSnapshot(t *testing.T, data map[string]int) *Snapshot {
	t.Helper()
	snap := &Snapshot{Buckets: make(map[string]bbolt.BucketStats)}
	for name, keys := range data {
		snap.Buckets[name] = bbolt.BucketStats{
			KeyN:      keys,
			LeafInuse: int64(keys) * 128,
		}
	}
	return snap
}

func TestAggregate_NilSnapshot(t *testing.T) {
	result := Aggregate(nil, 0)
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if result.Buckets != 0 {
		t.Errorf("expected 0 buckets, got %d", result.Buckets)
	}
}

func TestAggregate_EmptySnapshot(t *testing.T) {
	snap := &Snapshot{Buckets: map[string]bbolt.BucketStats{}}
	result := Aggregate(snap, 0)
	if result.Buckets != 0 {
		t.Errorf("expected 0 buckets, got %d", result.Buckets)
	}
}

func TestAggregate_TopByKeys(t *testing.T) {
	snap := makeAggregateSnapshot(t, map[string]int{
		"alpha": 10,
		"beta":  50,
		"gamma": 30,
	})
	result := Aggregate(snap, 0)
	if len(result.TopByKeys) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(result.TopByKeys))
	}
	if result.TopByKeys[0].Bucket != "beta" {
		t.Errorf("expected beta first, got %s", result.TopByKeys[0].Bucket)
	}
}

func TestAggregate_TopN(t *testing.T) {
	snap := makeAggregateSnapshot(t, map[string]int{
		"a": 1, "b": 2, "c": 3, "d": 4,
	})
	result := Aggregate(snap, 2)
	if len(result.TopByKeys) != 2 {
		t.Errorf("expected 2 top entries, got %d", len(result.TopByKeys))
	}
	if len(result.TopBySize) != 2 {
		t.Errorf("expected 2 top-by-size entries, got %d", len(result.TopBySize))
	}
}

func TestAggregate_Totals(t *testing.T) {
	snap := makeAggregateSnapshot(t, map[string]int{
		"x": 5,
		"y": 15,
	})
	result := Aggregate(snap, 0)
	if result.TotalKeys != 20 {
		t.Errorf("expected total keys 20, got %d", result.TotalKeys)
	}
	if result.TotalSize != 20*128 {
		t.Errorf("expected total size %d, got %d", 20*128, result.TotalSize)
	}
	if result.Buckets != 2 {
		t.Errorf("expected 2 buckets, got %d", result.Buckets)
	}
}
