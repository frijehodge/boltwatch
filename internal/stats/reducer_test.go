package stats

import (
	"strings"
	"testing"
	"time"
)

func makeReducerSnapshot(buckets []BucketStats) *Snapshot {
	return &Snapshot{
		Timestamp: time.Now(),
		Buckets:   buckets,
	}
}

func TestReduceSnapshot_Nil(t *testing.T) {
	result := ReduceSnapshot(nil, DefaultReduceConfig())
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

func TestReduceSnapshot_Empty(t *testing.T) {
	snap := makeReducerSnapshot(nil)
	result := ReduceSnapshot(snap, DefaultReduceConfig())
	if result != snap {
		t.Errorf("expected same snapshot returned for empty input")
	}
}

func TestReduceSnapshot_LimitsBuckets(t *testing.T) {
	buckets := []BucketStats{
		{Name: "a", KeyCount: 10},
		{Name: "b", KeyCount: 20},
		{Name: "c", KeyCount: 5},
		{Name: "d", KeyCount: 50},
	}
	snap := makeReducerSnapshot(buckets)
	cfg := ReduceConfig{MaxBuckets: 2, SortByKeys: true}
	result := ReduceSnapshot(snap, cfg)
	if len(result.Buckets) != 2 {
		t.Fatalf("expected 2 buckets, got %d", len(result.Buckets))
	}
	if result.Buckets[0].Name != "d" {
		t.Errorf("expected top bucket to be 'd', got '%s'", result.Buckets[0].Name)
	}
}

func TestReduceSnapshot_SortBySize(t *testing.T) {
	buckets := []BucketStats{
		{Name: "small", Size: 100},
		{Name: "large", Size: 9000},
		{Name: "medium", Size: 500},
	}
	snap := makeReducerSnapshot(buckets)
	cfg := ReduceConfig{MaxBuckets: 2, SortBySize: true}
	result := ReduceSnapshot(snap, cfg)
	if result.Buckets[0].Name != "large" {
		t.Errorf("expected 'large' first, got '%s'", result.Buckets[0].Name)
	}
}

func TestReduceSnapshot_NoLimitWhenZero(t *testing.T) {
	buckets := make([]BucketStats, 30)
	for i := range buckets {
		buckets[i] = BucketStats{Name: "b", KeyCount: i}
	}
	snap := makeReducerSnapshot(buckets)
	cfg := ReduceConfig{MaxBuckets: 0}
	result := ReduceSnapshot(snap, cfg)
	if len(result.Buckets) != 30 {
		t.Errorf("expected 30 buckets, got %d", len(result.Buckets))
	}
}

func TestFormatReduced_ContainsHeader(t *testing.T) {
	buckets := []BucketStats{{Name: "alpha", KeyCount: 3, Size: 1024}}
	snap := makeReducerSnapshot(buckets)
	out := FormatReduced(snap, 10)
	if !strings.Contains(out, "top 1 of 10") {
		t.Errorf("expected reduced header, got: %s", out)
	}
	if !strings.Contains(out, "alpha") {
		t.Errorf("expected bucket name in output")
	}
}

func TestFormatReduced_Nil(t *testing.T) {
	out := FormatReduced(nil, 0)
	if !strings.Contains(out, "no snapshot") {
		t.Errorf("expected nil message, got: %s", out)
	}
}
