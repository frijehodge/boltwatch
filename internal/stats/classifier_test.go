package stats

import (
	"testing"
	"time"
)

func makeClassifierSnapshot(ts time.Time, buckets map[string]BucketStats) *Snapshot {
	return &Snapshot{
		Timestamp: ts,
		Buckets:   buckets,
	}
}

func TestClassifyGrowth_NilInputs(t *testing.T) {
	result := ClassifyGrowth(nil, nil)
	if len(result) != 0 {
		t.Errorf("expected empty result for nil inputs, got %d", len(result))
	}
}

func TestClassifyGrowth_NewBucket(t *testing.T) {
	now := time.Now()
	prev := makeClassifierSnapshot(now.Add(-time.Second), map[string]BucketStats{})
	curr := makeClassifierSnapshot(now, map[string]BucketStats{
		"users": {KeyCount: 10, Size: 1024},
	})

	result := ClassifyGrowth(prev, curr)
	if result["users"].Class != GrowthClassNew {
		t.Errorf("expected new, got %s", result["users"].Class)
	}
}

func TestClassifyGrowth_RemovedBucket(t *testing.T) {
	now := time.Now()
	prev := makeClassifierSnapshot(now.Add(-time.Second), map[string]BucketStats{
		"logs": {KeyCount: 5, Size: 512},
	})
	curr := makeClassifierSnapshot(now, map[string]BucketStats{})

	result := ClassifyGrowth(prev, curr)
	if result["logs"].Class != GrowthClassRemoved {
		t.Errorf("expected removed, got %s", result["logs"].Class)
	}
}

func TestClassifyGrowth_GrowingBucket(t *testing.T) {
	now := time.Now()
	prev := makeClassifierSnapshot(now.Add(-time.Second), map[string]BucketStats{
		"events": {KeyCount: 100, Size: 2048},
	})
	curr := makeClassifierSnapshot(now, map[string]BucketStats{
		"events": {KeyCount: 150, Size: 3072},
	})

	result := ClassifyGrowth(prev, curr)
	c := result["events"]
	if c.Class != GrowthClassGrowing {
		t.Errorf("expected growing, got %s", c.Class)
	}
	if c.KeyDelta != 50 {
		t.Errorf("expected KeyDelta=50, got %d", c.KeyDelta)
	}
	if c.GrowthRate <= 0 {
		t.Errorf("expected positive growth rate, got %f", c.GrowthRate)
	}
}

func TestClassifyGrowth_StableBucket(t *testing.T) {
	now := time.Now()
	prev := makeClassifierSnapshot(now.Add(-time.Second), map[string]BucketStats{
		"cache": {KeyCount: 50, Size: 1024},
	})
	curr := makeClassifierSnapshot(now, map[string]BucketStats{
		"cache": {KeyCount: 50, Size: 1024},
	})

	result := ClassifyGrowth(prev, curr)
	if result["cache"].Class != GrowthClassStable {
		t.Errorf("expected stable, got %s", result["cache"].Class)
	}
}

func TestFormatClassifications_Empty(t *testing.T) {
	out := FormatClassifications(map[string]BucketClassification{})
	if out == "" {
		t.Error("expected non-empty output for empty classifications")
	}
}

func TestFormatClassifications_ContainsHeaders(t *testing.T) {
	classifications := map[string]BucketClassification{
		"users": {Bucket: "users", Class: GrowthClassGrowing, KeyDelta: 5},
	}
	out := FormatClassifications(classifications)
	for _, header := range []string{"Bucket", "Class", "KeyDelta"} {
		if !containsStr(out, header) {
			t.Errorf("expected output to contain %q", header)
		}
	}
}
