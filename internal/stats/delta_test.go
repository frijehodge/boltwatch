package stats

import (
	"testing"

	bolt "go.etcd.io/bbolt"
)

func makeSnapshot(buckets map[string]bolt.BucketStats) *Snapshot {
	s := &Snapshot{
		Buckets: make(map[string]bolt.BucketStats),
	}
	for k, v := range buckets {
		s.Buckets[k] = v
	}
	return s
}

func TestComputeDeltas_NilInputs(t *testing.T) {
	if got := ComputeDeltas(nil, nil); got != nil {
		t.Errorf("expected nil, got %v", got)
	}
	if got := ComputeDeltas(nil, &Snapshot{}); got != nil {
		t.Errorf("expected nil for nil prev, got %v", got)
	}
}

func TestComputeDeltas_NoChanges(t *testing.T) {
	stats := bolt.BucketStats{KeyN: 5, LeafInuse: 1024}
	prev := makeSnapshot(map[string]bolt.BucketStats{"alpha": stats})
	curr := makeSnapshot(map[string]bolt.BucketStats{"alpha": stats})

	deltas := ComputeDeltas(prev, curr)
	if len(deltas) != 0 {
		t.Errorf("expected 0 deltas, got %d", len(deltas))
	}
}

func TestComputeDeltas_KeysChanged(t *testing.T) {
	prev := makeSnapshot(map[string]bolt.BucketStats{"alpha": {KeyN: 3, LeafInuse: 512}})
	curr := makeSnapshot(map[string]bolt.BucketStats{"alpha": {KeyN: 7, LeafInuse: 512}})

	deltas := ComputeDeltas(prev, curr)
	if len(deltas) != 1 {
		t.Fatalf("expected 1 delta, got %d", len(deltas))
	}
	if deltas[0].KeysDiff != 4 {
		t.Errorf("expected KeysDiff=4, got %d", deltas[0].KeysDiff)
	}
}

func TestComputeDeltas_NewBucket(t *testing.T) {
	prev := makeSnapshot(map[string]bolt.BucketStats{})
	curr := makeSnapshot(map[string]bolt.BucketStats{"beta": {KeyN: 2, LeafInuse: 256}})

	deltas := ComputeDeltas(prev, curr)
	if len(deltas) != 1 {
		t.Fatalf("expected 1 delta, got %d", len(deltas))
	}
	if deltas[0].Bucket != "beta" {
		t.Errorf("expected bucket 'beta', got %q", deltas[0].Bucket)
	}
	if deltas[0].KeysDiff != 2 {
		t.Errorf("expected KeysDiff=2, got %d", deltas[0].KeysDiff)
	}
}

func TestComputeDeltas_RemovedBucket(t *testing.T) {
	prev := makeSnapshot(map[string]bolt.BucketStats{"gone": {KeyN: 10, LeafInuse: 2048}})
	curr := makeSnapshot(map[string]bolt.BucketStats{})

	deltas := ComputeDeltas(prev, curr)
	if len(deltas) != 1 {
		t.Fatalf("expected 1 delta, got %d", len(deltas))
	}
	if deltas[0].KeysDiff != -10 {
		t.Errorf("expected KeysDiff=-10, got %d", deltas[0].KeysDiff)
	}
}

func TestComputeDeltas_MultipleBuckets(t *testing.T) {
	prev := makeSnapshot(map[string]bolt.BucketStats{
		"alpha": {KeyN: 3, LeafInuse: 512},
		"beta":  {KeyN: 5, LeafInuse: 1024},
	})
	curr := makeSnapshot(map[string]bolt.BucketStats{
		"alpha": {KeyN: 6, LeafInuse: 512},
		"beta":  {KeyN: 5, LeafInuse: 1024},
		"gamma": {KeyN: 1, LeafInuse: 128},
	})

	deltas := ComputeDeltas(prev, curr)
	if len(deltas) != 2 {
		t.Fatalf("expected 2 deltas, got %d", len(deltas))
	}
}

func TestHasChanges(t *testing.T) {
	if HasChanges(nil) {
		t.Error("expected false for nil deltas")
	}
	if HasChanges([]Delta{}) {
		t.Error("expected false for empty deltas")
	}
	if !HasChanges([]Delta{{Bucket: "x", KeysDiff: 1}}) {
		t.Error("expected true for non-empty deltas")
	}
}
