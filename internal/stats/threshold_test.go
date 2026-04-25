package stats

import (
	"testing"

	"go.etcd.io/bbolt"
)

func makeSnapshotForThreshold(buckets map[string]bolt.BucketStats) *Snapshot {
	snap := &Snapshot{
		Buckets: make(map[string]bolt.BucketStats, len(buckets)),
	}
	for k, v := range buckets {
		snap.Buckets[k] = v
	}
	return snap
}

func TestCheckThresholds_NilSnapshot(t *testing.T) {
	violations := CheckThresholds(nil, []ThresholdConfig{{Bucket: "a", MaxKeys: 10}})
	if len(violations) != 0 {
		t.Errorf("expected no violations for nil snapshot, got %d", len(violations))
	}
}

func TestCheckThresholds_NoViolations(t *testing.T) {
	snap := makeSnapshotForThreshold(map[string]bolt.BucketStats{
		"users": {KeyN: 5, LeafInuse: 1024},
	})
	configs := []ThresholdConfig{
		{Bucket: "users", MaxKeys: 100, MaxSizeBytes: 1024 * 1024, Level: ThresholdWarning},
	}
	violations := CheckThresholds(snap, configs)
	if len(violations) != 0 {
		t.Errorf("expected 0 violations, got %d", len(violations))
	}
}

func TestCheckThresholds_KeysExceeded(t *testing.T) {
	snap := makeSnapshotForThreshold(map[string]bolt.BucketStats{
		"orders": {KeyN: 200, LeafInuse: 512},
	})
	configs := []ThresholdConfig{
		{Bucket: "orders", MaxKeys: 100, Level: ThresholdWarning},
	}
	violations := CheckThresholds(snap, configs)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Field != "keys" {
		t.Errorf("expected field 'keys', got %q", violations[0].Field)
	}
	if violations[0].Level != ThresholdWarning {
		t.Errorf("expected WARNING level, got %q", violations[0].Level)
	}
}

func TestCheckThresholds_SizeExceeded(t *testing.T) {
	snap := makeSnapshotForThreshold(map[string]bolt.BucketStats{
		"logs": {KeyN: 10, LeafInuse: 2048},
	})
	configs := []ThresholdConfig{
		{Bucket: "logs", MaxSizeBytes: 1024, Level: ThresholdCritical},
	}
	violations := CheckThresholds(snap, configs)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Field != "size" {
		t.Errorf("expected field 'size', got %q", violations[0].Field)
	}
}

func TestCheckThresholds_UnknownBucketSkipped(t *testing.T) {
	snap := makeSnapshotForThreshold(map[string]bolt.BucketStats{
		"users": {KeyN: 5},
	})
	configs := []ThresholdConfig{
		{Bucket: "missing", MaxKeys: 1, Level: ThresholdWarning},
	}
	violations := CheckThresholds(snap, configs)
	if len(violations) != 0 {
		t.Errorf("expected 0 violations for unknown bucket, got %d", len(violations))
	}
}
