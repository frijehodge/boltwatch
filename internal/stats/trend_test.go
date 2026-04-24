package stats

import (
	"testing"
	"time"

	bolt "go.etcd.io/bbolt"
)

func makeSnapshotWithKeys(bucket string, keys int, size int64, ts time.Time) Snapshot {
	return Snapshot{
		Timestamp: ts,
		Buckets: map[string]bolt.BucketStats{
			bucket: {
				KeyN:      keys,
				LeafInuse: size,
			},
		},
	}
}

func TestComputeTrends_NilOrSingle(t *testing.T) {
	result := ComputeTrends(nil)
	if result != nil {
		t.Errorf("expected nil for nil input, got %v", result)
	}

	single := []Snapshot{makeSnapshotWithKeys("a", 5, 100, time.Now())}
	result = ComputeTrends(single)
	if result != nil {
		t.Errorf("expected nil for single snapshot, got %v", result)
	}
}

func TestComputeTrends_StableBucket(t *testing.T) {
	old := makeSnapshotWithKeys("users", 10, 512, time.Now().Add(-time.Minute))
	new := makeSnapshotWithKeys("users", 10, 512, time.Now())

	trends := ComputeTrends([]Snapshot{old, new})
	if len(trends) != 1 {
		t.Fatalf("expected 1 trend, got %d", len(trends))
	}
	if trends[0].KeysTrend != TrendStable {
		t.Errorf("expected stable keys trend, got %v", trends[0].KeysTrend)
	}
	if trends[0].SizeTrend != TrendStable {
		t.Errorf("expected stable size trend, got %v", trends[0].SizeTrend)
	}
}

func TestComputeTrends_GrowingBucket(t *testing.T) {
	old := makeSnapshotWithKeys("orders", 5, 256, time.Now().Add(-time.Minute))
	new := makeSnapshotWithKeys("orders", 20, 1024, time.Now())

	trends := ComputeTrends([]Snapshot{old, new})
	if len(trends) != 1 {
		t.Fatalf("expected 1 trend, got %d", len(trends))
	}
	if trends[0].KeysTrend != TrendUp {
		t.Errorf("expected keys trend up, got %v", trends[0].KeysTrend)
	}
	if trends[0].SizeTrend != TrendUp {
		t.Errorf("expected size trend up, got %v", trends[0].SizeTrend)
	}
}

func TestComputeTrends_ShrinkingBucket(t *testing.T) {
	old := makeSnapshotWithKeys("cache", 100, 4096, time.Now().Add(-time.Minute))
	new := makeSnapshotWithKeys("cache", 40, 1024, time.Now())

	trends := ComputeTrends([]Snapshot{old, new})
	if trends[0].KeysTrend != TrendDown {
		t.Errorf("expected keys trend down, got %v", trends[0].KeysTrend)
	}
}

func TestTrendDirection_Symbol(t *testing.T) {
	if TrendUp.TrendSymbol() != "↑" {
		t.Errorf("unexpected symbol for TrendUp")
	}
	if TrendDown.TrendSymbol() != "↓" {
		t.Errorf("unexpected symbol for TrendDown")
	}
	if TrendStable.TrendSymbol() != "→" {
		t.Errorf("unexpected symbol for TrendStable")
	}
}
