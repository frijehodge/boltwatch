package stats

import (
	"testing"
	"time"
)

func makeLabelerSnapshot(buckets map[string]BucketStats) *Snapshot {
	return &Snapshot{
		Timestamp: time.Now(),
		Buckets:   buckets,
	}
}

func TestLabelSnapshot_Nil(t *testing.T) {
	result := LabelSnapshot(nil, DefaultLabelConfig())
	if result != nil {
		t.Errorf("expected nil for nil snapshot, got %v", result)
	}
}

func TestLabelSnapshot_Empty(t *testing.T) {
	snap := makeLabelerSnapshot(map[string]BucketStats{})
	result := LabelSnapshot(snap, DefaultLabelConfig())
	if result != nil {
		t.Errorf("expected nil for empty snapshot, got %v", result)
	}
}

func TestLabelSnapshot_EmptyBucket(t *testing.T) {
	snap := makeLabelerSnapshot(map[string]BucketStats{
		"empty": {KeyCount: 0, Size: 0},
	})
	result := LabelSnapshot(snap, DefaultLabelConfig())
	if len(result) != 1 {
		t.Fatalf("expected 1 labeled bucket, got %d", len(result))
	}
	if result[0].Label != LabelEmpty {
		t.Errorf("expected LabelEmpty, got %s", result[0].Label)
	}
}

func TestLabelSnapshot_SmallBucket(t *testing.T) {
	snap := makeLabelerSnapshot(map[string]BucketStats{
		"small": {KeyCount: 50, Size: 512},
	})
	result := LabelSnapshot(snap, DefaultLabelConfig())
	if result[0].Label != LabelSmall {
		t.Errorf("expected LabelSmall, got %s", result[0].Label)
	}
}

func TestLabelSnapshot_MediumBucket(t *testing.T) {
	snap := makeLabelerSnapshot(map[string]BucketStats{
		"medium": {KeyCount: 5000, Size: 512 * 1024},
	})
	result := LabelSnapshot(snap, DefaultLabelConfig())
	if result[0].Label != LabelMedium {
		t.Errorf("expected LabelMedium, got %s", result[0].Label)
	}
}

func TestLabelSnapshot_LargeBucket(t *testing.T) {
	snap := makeLabelerSnapshot(map[string]BucketStats{
		"large": {KeyCount: 50000, Size: 20 * 1024 * 1024},
	})
	result := LabelSnapshot(snap, DefaultLabelConfig())
	if result[0].Label != LabelLarge {
		t.Errorf("expected LabelLarge, got %s", result[0].Label)
	}
}

func TestLabelSnapshot_MultipleBuckets(t *testing.T) {
	snap := makeLabelerSnapshot(map[string]BucketStats{
		"a": {KeyCount: 0},
		"b": {KeyCount: 100},
		"c": {KeyCount: 20000},
	})
	result := LabelSnapshot(snap, DefaultLabelConfig())
	if len(result) != 3 {
		t.Fatalf("expected 3 labeled buckets, got %d", len(result))
	}
}
