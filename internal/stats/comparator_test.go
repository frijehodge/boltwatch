package stats

import (
	"testing"
)

func makeCompSnapshot(buckets map[string]BucketStats) *Snapshot {
	s := &Snapshot{Buckets: make(map[string]BucketStats)}
	for k, v := range buckets {
		s.Buckets[k] = v
	}
	return s
}

func TestCompareSnapshots_NilInputs(t *testing.T) {
	result := CompareSnapshots(nil, nil)
	if len(result.Comparisons) != 0 {
		t.Errorf("expected empty comparisons for nil inputs")
	}

	snap := makeCompSnapshot(map[string]BucketStats{"a": {KeyCount: 1}})
	result = CompareSnapshots(nil, snap)
	if len(result.Comparisons) != 0 {
		t.Errorf("expected empty comparisons when a is nil")
	}
}

func TestCompareSnapshots_NoChanges(t *testing.T) {
	bs := map[string]BucketStats{"users": {KeyCount: 10, Size: 1024}}
	a := makeCompSnapshot(bs)
	b := makeCompSnapshot(bs)

	result := CompareSnapshots(a, b)
	if result.AddedCount != 0 || result.RemovedCount != 0 || result.ChangedCount != 0 {
		t.Errorf("expected no changes, got added=%d removed=%d changed=%d",
			result.AddedCount, result.RemovedCount, result.ChangedCount)
	}
}

func TestCompareSnapshots_NewBucket(t *testing.T) {
	a := makeCompSnapshot(map[string]BucketStats{})
	b := makeCompSnapshot(map[string]BucketStats{"logs": {KeyCount: 5, Size: 512}})

	result := CompareSnapshots(a, b)
	if result.AddedCount != 1 {
		t.Errorf("expected 1 added bucket, got %d", result.AddedCount)
	}
	if len(result.Comparisons) != 1 || !result.Comparisons[0].IsNew {
		t.Errorf("expected comparison to be marked as new")
	}
}

func TestCompareSnapshots_RemovedBucket(t *testing.T) {
	a := makeCompSnapshot(map[string]BucketStats{"sessions": {KeyCount: 3, Size: 256}})
	b := makeCompSnapshot(map[string]BucketStats{})

	result := CompareSnapshots(a, b)
	if result.RemovedCount != 1 {
		t.Errorf("expected 1 removed bucket, got %d", result.RemovedCount)
	}
	if len(result.Comparisons) != 1 || !result.Comparisons[0].IsRemoved {
		t.Errorf("expected comparison to be marked as removed")
	}
}

func TestCompareSnapshots_ChangedBucket(t *testing.T) {
	a := makeCompSnapshot(map[string]BucketStats{"events": {KeyCount: 10, Size: 1000}})
	b := makeCompSnapshot(map[string]BucketStats{"events": {KeyCount: 15, Size: 1500}})

	result := CompareSnapshots(a, b)
	if result.ChangedCount != 1 {
		t.Errorf("expected 1 changed bucket, got %d", result.ChangedCount)
	}
	c := result.Comparisons[0]
	if c.KeysDiff != 5 {
		t.Errorf("expected KeysDiff=5, got %d", c.KeysDiff)
	}
	if c.SizeDiff != 500 {
		t.Errorf("expected SizeDiff=500, got %d", c.SizeDiff)
	}
}

func TestCompareSnapshots_SortedByName(t *testing.T) {
	bs := map[string]BucketStats{
		"zebra": {KeyCount: 1},
		"alpha": {KeyCount: 2},
		"mango": {KeyCount: 3},
	}
	a := makeCompSnapshot(bs)
	b := makeCompSnapshot(bs)

	result := CompareSnapshots(a, b)
	names := []string{result.Comparisons[0].Bucket, result.Comparisons[1].Bucket, result.Comparisons[2].Bucket}
	if names[0] != "alpha" || names[1] != "mango" || names[2] != "zebra" {
		t.Errorf("expected sorted bucket names, got %v", names)
	}
}
