package stats

import (
	"testing"
	"time"

	bbolt "go.etcd.io/bbolt"
)

func makeDigestSnapshot(buckets map[string]int) *Snapshot {
	s := &Snapshot{
		Timestamp: time.Now(),
		Buckets:   make(map[string]bbolt.BucketStats),
	}
	for name, keys := range buckets {
		s.Buckets[name] = bbolt.BucketStats{
			KeyN:      keys,
			LeafInuse: keys * 64,
		}
	}
	return s
}

func TestDigestSnapshot_Nil(t *testing.T) {
	r := DigestSnapshot(nil)
	if r == nil {
		t.Fatal("expected non-nil result for nil snapshot")
	}
	if r.BucketCount != 0 || r.KeyCount != 0 {
		t.Errorf("expected zero counts for nil snapshot, got %+v", r)
	}
}

func TestDigestSnapshot_Empty(t *testing.T) {
	s := &Snapshot{Timestamp: time.Now(), Buckets: map[string]bbolt.BucketStats{}}
	r := DigestSnapshot(s)
	if r.BucketCount != 0 {
		t.Errorf("expected 0 buckets, got %d", r.BucketCount)
	}
}

func TestDigestSnapshot_CountsCorrect(t *testing.T) {
	s := makeDigestSnapshot(map[string]int{"alpha": 10, "beta": 20})
	r := DigestSnapshot(s)
	if r.BucketCount != 2 {
		t.Errorf("expected 2 buckets, got %d", r.BucketCount)
	}
	if r.KeyCount != 30 {
		t.Errorf("expected 30 keys, got %d", r.KeyCount)
	}
	if r.Hash == "" {
		t.Error("expected non-empty hash")
	}
}

func TestDigestSnapshot_StableHash(t *testing.T) {
	s1 := makeDigestSnapshot(map[string]int{"a": 5, "b": 10})
	s2 := makeDigestSnapshot(map[string]int{"b": 10, "a": 5})
	r1 := DigestSnapshot(s1)
	r2 := DigestSnapshot(s2)
	if r1.Hash != r2.Hash {
		t.Errorf("expected same hash for same data regardless of insertion order: %s vs %s", r1.Hash, r2.Hash)
	}
}

func TestDigestSnapshot_DifferentHash(t *testing.T) {
	s1 := makeDigestSnapshot(map[string]int{"a": 5})
	s2 := makeDigestSnapshot(map[string]int{"a": 6})
	r1 := DigestSnapshot(s1)
	r2 := DigestSnapshot(s2)
	if r1.Hash == r2.Hash {
		t.Error("expected different hashes for different data")
	}
}

func TestDigestsMatch_BothNil(t *testing.T) {
	if !DigestsMatch(nil, nil) {
		t.Error("expected nil digests to match")
	}
}

func TestDigestsMatch_OneNil(t *testing.T) {
	s := makeDigestSnapshot(map[string]int{"a": 1})
	r := DigestSnapshot(s)
	if DigestsMatch(r, nil) {
		t.Error("expected non-nil and nil to not match")
	}
}

func TestDigestsMatch_Equal(t *testing.T) {
	s1 := makeDigestSnapshot(map[string]int{"x": 3})
	s2 := makeDigestSnapshot(map[string]int{"x": 3})
	if !DigestsMatch(DigestSnapshot(s1), DigestSnapshot(s2)) {
		t.Error("expected matching digests for identical snapshots")
	}
}
