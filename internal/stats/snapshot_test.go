package stats

import (
	"testing"
	"time"
)

func makeBucketStats(name string, keys int, size int64) BucketStats {
	return BucketStats{
		Name:     name,
		KeyN:     keys,
		LeafSize: size,
	}
}

func TestNewSnapshot_SetsTimestamp(t *testing.T) {
	before := time.Now()
	s := NewSnapshot(nil)
	after := time.Now()

	if s.Timestamp.Before(before) || s.Timestamp.After(after) {
		t.Errorf("expected timestamp between %v and %v, got %v", before, after, s.Timestamp)
	}
}

func TestSnapshot_IsEmpty_True(t *testing.T) {
	s := NewSnapshot(nil)
	if !s.IsEmpty() {
		t.Error("expected snapshot with nil buckets to be empty")
	}
}

func TestSnapshot_IsEmpty_False(t *testing.T) {
	s := NewSnapshot([]BucketStats{makeBucketStats("b1", 5, 1024)})
	if s.IsEmpty() {
		t.Error("expected snapshot with buckets to not be empty")
	}
}

func TestSnapshot_TotalKeys(t *testing.T) {
	buckets := []BucketStats{
		makeBucketStats("a", 10, 512),
		makeBucketStats("b", 20, 1024),
	}
	s := NewSnapshot(buckets)
	got := s.TotalKeys()
	if got != 30 {
		t.Errorf("expected TotalKeys=30, got %d", got)
	}
}

func TestSnapshot_TotalSize(t *testing.T) {
	buckets := []BucketStats{
		makeBucketStats("a", 10, 512),
		makeBucketStats("b", 20, 1024),
	}
	s := NewSnapshot(buckets)
	got := s.TotalSize()
	if got != 1536 {
		t.Errorf("expected TotalSize=1536, got %d", got)
	}
}

func TestSnapshot_BucketNames(t *testing.T) {
	buckets := []BucketStats{
		makeBucketStats("alpha", 1, 100),
		makeBucketStats("beta", 2, 200),
		makeBucketStats("gamma", 3, 300),
	}
	s := NewSnapshot(buckets)
	names := s.BucketNames()
	if len(names) != 3 {
		t.Fatalf("expected 3 names, got %d", len(names))
	}
	expected := []string{"alpha", "beta", "gamma"}
	for i, name := range names {
		if name != expected[i] {
			t.Errorf("expected name[%d]=%q, got %q", i, expected[i], name)
		}
	}
}

func TestSnapshot_Empty_TotalKeys_Zero(t *testing.T) {
	s := NewSnapshot([]BucketStats{})
	if s.TotalKeys() != 0 {
		t.Error("expected TotalKeys=0 for empty snapshot")
	}
}
