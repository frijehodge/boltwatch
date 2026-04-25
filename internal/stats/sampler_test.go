package stats

import (
	"testing"
	"time"

	"go.etcd.io/bbolt"
)

func makeSamplerSnapshot(keys map[string]int) *Snapshot {
	buckets := make(map[string]BucketStats)
	for name, k := range keys {
		buckets[name] = BucketStats{
			KeyN:     k,
			BoltStats: bbolt.BucketStats{LeafInuse: k * 64},
		}
	}
	return &Snapshot{Buckets: buckets, Timestamp: time.Now()}
}

func TestNewSampler_DefaultMax(t *testing.T) {
	s := NewSampler(0)
	if s.max != 60 {
		t.Fatalf("expected default max 60, got %d", s.max)
	}
}

func TestSampler_Record_Nil(t *testing.T) {
	s := NewSampler(10)
	s.Record(nil)
	if s.Len() != 0 {
		t.Fatal("expected 0 samples after recording nil snapshot")
	}
}

func TestSampler_Record_Empty(t *testing.T) {
	s := NewSampler(10)
	s.Record(&Snapshot{})
	if s.Len() != 0 {
		t.Fatal("expected 0 samples for empty snapshot")
	}
}

func TestSampler_Record_AddsSamples(t *testing.T) {
	s := NewSampler(100)
	snap := makeSamplerSnapshot(map[string]int{"users": 5, "orders": 3})
	s.Record(snap)
	if s.Len() != 2 {
		t.Fatalf("expected 2 samples, got %d", s.Len())
	}
}

func TestSampler_ForBucket(t *testing.T) {
	s := NewSampler(100)
	s.Record(makeSamplerSnapshot(map[string]int{"users": 5, "orders": 3}))
	s.Record(makeSamplerSnapshot(map[string]int{"users": 7}))

	samples := s.ForBucket("users")
	if len(samples) != 2 {
		t.Fatalf("expected 2 user samples, got %d", len(samples))
	}
	if samples[0].Keys != 5 || samples[1].Keys != 7 {
		t.Fatal("unexpected key counts in user samples")
	}
}

func TestSampler_Eviction(t *testing.T) {
	s := NewSampler(3)
	for i := 0; i < 5; i++ {
		s.Record(makeSamplerSnapshot(map[string]int{"b": i}))
	}
	if s.Len() != 3 {
		t.Fatalf("expected 3 samples after eviction, got %d", s.Len())
	}
	all := s.All()
	if all[2].Keys != 4 {
		t.Fatalf("expected newest sample keys=4, got %d", all[2].Keys)
	}
}

func TestSampler_Reset(t *testing.T) {
	s := NewSampler(10)
	s.Record(makeSamplerSnapshot(map[string]int{"a": 1}))
	s.Reset()
	if s.Len() != 0 {
		t.Fatal("expected 0 samples after reset")
	}
}
