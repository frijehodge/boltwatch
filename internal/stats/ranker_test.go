package stats

import (
	"testing"
	"time"
)

func makeRankerSnapshot(buckets map[string]BucketStats) *Snapshot {
	return &Snapshot{
		Timestamp: time.Now(),
		Buckets:   buckets,
	}
}

func TestRankBuckets_NilSnapshot(t *testing.T) {
	result := RankBuckets(nil, DefaultRankConfig())
	if result != nil {
		t.Errorf("expected nil for nil snapshot, got %v", result)
	}
}

func TestRankBuckets_EmptySnapshot(t *testing.T) {
	snap := makeRankerSnapshot(map[string]BucketStats{})
	result := RankBuckets(snap, DefaultRankConfig())
	if result != nil {
		t.Errorf("expected nil for empty snapshot, got %v", result)
	}
}

func TestRankBuckets_OrderByScore(t *testing.T) {
	snap := makeRankerSnapshot(map[string]BucketStats{
		"small":  {Keys: 10, Size: 1024},
		"large":  {Keys: 1000, Size: 1048576},
		"medium": {Keys: 100, Size: 102400},
	})

	result := RankBuckets(snap, DefaultRankConfig())
	if len(result) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(result))
	}
	if result[0].Bucket != "large" {
		t.Errorf("expected 'large' first, got %q", result[0].Bucket)
	}
	if result[len(result)-1].Bucket != "small" {
		t.Errorf("expected 'small' last, got %q", result[len(result)-1].Bucket)
	}
}

func TestRankBuckets_TopN(t *testing.T) {
	snap := makeRankerSnapshot(map[string]BucketStats{
		"a": {Keys: 300, Size: 3000},
		"b": {Keys: 200, Size: 2000},
		"c": {Keys: 100, Size: 1000},
		"d": {Keys: 50, Size: 500},
	})

	cfg := DefaultRankConfig()
	cfg.TopN = 2
	result := RankBuckets(snap, cfg)
	if len(result) != 2 {
		t.Fatalf("expected 2 entries with TopN=2, got %d", len(result))
	}
}

func TestRankBuckets_KeyWeightOnly(t *testing.T) {
	snap := makeRankerSnapshot(map[string]BucketStats{
		"manykeys": {Keys: 500, Size: 100},
		"bigsize":  {Keys: 10, Size: 999999},
	})

	cfg := RankConfig{KeyWeight: 1.0, SizeWeight: 0.0}
	result := RankBuckets(snap, cfg)
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
	if result[0].Bucket != "manykeys" {
		t.Errorf("expected 'manykeys' first with key-only weight, got %q", result[0].Bucket)
	}
}

func TestRankBuckets_RankScoreRange(t *testing.T) {
	snap := makeRankerSnapshot(map[string]BucketStats{
		"alpha": {Keys: 100, Size: 2048},
		"beta":  {Keys: 50, Size: 1024},
	})

	result := RankBuckets(snap, DefaultRankConfig())
	for _, e := range result {
		if e.Rank < 0 || e.Rank > 1.0 {
			t.Errorf("rank %f out of [0,1] range for bucket %q", e.Rank, e.Bucket)
		}
	}
}
