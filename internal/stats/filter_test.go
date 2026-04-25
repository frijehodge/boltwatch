package stats

import (
	"testing"
	"time"
)

func makeFilterSnapshot() *Snapshot {
	return &Snapshot{
		Timestamp: time.Now(),
		Buckets: map[string]BucketStats{
			"users":    {KeyCount: 100, Size: 2048},
			"sessions": {KeyCount: 5, Size: 512},
			"logs":     {KeyCount: 50, Size: 8192},
			"log_arch": {KeyCount: 10, Size: 4096},
		},
	}
}

func TestFilterSnapshot_Nil(t *testing.T) {
	result := FilterSnapshot(nil, FilterConfig{})
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

func TestFilterSnapshot_NoFilter(t *testing.T) {
	snap := makeFilterSnapshot()
	result := FilterSnapshot(snap, FilterConfig{})
	if len(result.Buckets) != len(snap.Buckets) {
		t.Errorf("expected %d buckets, got %d", len(snap.Buckets), len(result.Buckets))
	}
}

func TestFilterSnapshot_ByPrefix(t *testing.T) {
	snap := makeFilterSnapshot()
	result := FilterSnapshot(snap, FilterConfig{Prefix: "log"})
	if len(result.Buckets) != 2 {
		t.Errorf("expected 2 buckets with prefix 'log', got %d", len(result.Buckets))
	}
	if _, ok := result.Buckets["logs"]; !ok {
		t.Error("expected 'logs' bucket to be present")
	}
	if _, ok := result.Buckets["log_arch"]; !ok {
		t.Error("expected 'log_arch' bucket to be present")
	}
}

func TestFilterSnapshot_ByMinKeys(t *testing.T) {
	snap := makeFilterSnapshot()
	result := FilterSnapshot(snap, FilterConfig{MinKeys: 20})
	if len(result.Buckets) != 2 {
		t.Errorf("expected 2 buckets with MinKeys=20, got %d", len(result.Buckets))
	}
}

func TestFilterSnapshot_ByMinSize(t *testing.T) {
	snap := makeFilterSnapshot()
	result := FilterSnapshot(snap, FilterConfig{MinSize: 2000})
	if len(result.Buckets) != 3 {
		t.Errorf("expected 3 buckets with MinSize=2000, got %d", len(result.Buckets))
	}
}

func TestFilterSnapshot_CombinedFilters(t *testing.T) {
	snap := makeFilterSnapshot()
	result := FilterSnapshot(snap, FilterConfig{Prefix: "log", MinKeys: 20})
	if len(result.Buckets) != 1 {
		t.Errorf("expected 1 bucket, got %d", len(result.Buckets))
	}
	if _, ok := result.Buckets["logs"]; !ok {
		t.Error("expected 'logs' bucket to be present")
	}
}

func TestMatchesFilter_Pass(t *testing.T) {
	bs := BucketStats{KeyCount: 50, Size: 1024}
	if !MatchesFilter("users", bs, FilterConfig{}) {
		t.Error("expected bucket to match empty filter")
	}
}

func TestMatchesFilter_FailPrefix(t *testing.T) {
	bs := BucketStats{KeyCount: 50, Size: 1024}
	if MatchesFilter("users", bs, FilterConfig{Prefix: "log"}) {
		t.Error("expected bucket to not match prefix filter")
	}
}
