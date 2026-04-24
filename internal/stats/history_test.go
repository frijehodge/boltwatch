package stats

import (
	"testing"
	"time"
)

func TestNewHistory_DefaultMaxSize(t *testing.T) {
	h := NewHistory(0)
	if h.maxSize != 60 {
		t.Errorf("expected default maxSize 60, got %d", h.maxSize)
	}
}

func TestHistory_AddAndLen(t *testing.T) {
	h := NewHistory(5)
	if h.Len() != 0 {
		t.Fatal("expected empty history")
	}

	h.Add([]BucketStats{{Name: "a", KeyCount: 1}})
	h.Add([]BucketStats{{Name: "b", KeyCount: 2}})

	if h.Len() != 2 {
		t.Errorf("expected length 2, got %d", h.Len())
	}
}

func TestHistory_Eviction(t *testing.T) {
	h := NewHistory(3)
	for i := 0; i < 5; i++ {
		h.Add([]BucketStats{{Name: "bucket", KeyCount: i}})
	}

	if h.Len() != 3 {
		t.Errorf("expected length 3 after eviction, got %d", h.Len())
	}

	// Oldest entry should have KeyCount == 2 (entries 0,1 evicted)
	all := h.All()
	if all[0].Buckets[0].KeyCount != 2 {
		t.Errorf("expected oldest KeyCount 2, got %d", all[0].Buckets[0].KeyCount)
	}
}

func TestHistory_Latest_Empty(t *testing.T) {
	h := NewHistory(10)
	if h.Latest() != nil {
		t.Error("expected nil for empty history")
	}
}

func TestHistory_Latest_ReturnsNewest(t *testing.T) {
	h := NewHistory(10)
	h.Add([]BucketStats{{Name: "old", KeyCount: 1}})
	time.Sleep(time.Millisecond)
	h.Add([]BucketStats{{Name: "new", KeyCount: 99}})

	latest := h.Latest()
	if latest == nil {
		t.Fatal("expected non-nil latest")
	}
	if latest.Buckets[0].KeyCount != 99 {
		t.Errorf("expected KeyCount 99, got %d", latest.Buckets[0].KeyCount)
	}
}

func TestHistory_All_ReturnsCopy(t *testing.T) {
	h := NewHistory(10)
	h.Add([]BucketStats{{Name: "x", KeyCount: 7}})

	all := h.All()
	all[0].Buckets[0].KeyCount = 999

	latest := h.Latest()
	if latest.Buckets[0].KeyCount == 999 {
		t.Error("All() should return an independent copy")
	}
}
