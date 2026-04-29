package stats

import (
	"testing"
	"time"
)

func makePinnedSnapshot(keys int) *Snapshot {
	buckets := map[string]BucketStats{
		"alpha": {KeyCount: keys, Size: int64(keys * 64)},
	}
	return &Snapshot{
		Timestamp: time.Now(),
		Buckets:   buckets,
	}
}

func TestNewPinBoard_Empty(t *testing.T) {
	pb := NewPinBoard()
	if pb.Len() != 0 {
		t.Fatalf("expected 0 pins, got %d", pb.Len())
	}
}

func TestPinBoard_Pin_NilSnapshot(t *testing.T) {
	pb := NewPinBoard()
	pb.Pin("test", nil)
	if pb.Len() != 0 {
		t.Fatal("nil snapshot should not be pinned")
	}
}

func TestPinBoard_Pin_StoresSnapshot(t *testing.T) {
	pb := NewPinBoard()
	s := makePinnedSnapshot(10)
	pb.Pin("baseline", s)
	if pb.Len() != 1 {
		t.Fatalf("expected 1 pin, got %d", pb.Len())
	}
}

func TestPinBoard_Get_ReturnsPin(t *testing.T) {
	pb := NewPinBoard()
	s := makePinnedSnapshot(5)
	pb.Pin("snap1", s)
	p := pb.Get("snap1")
	if p == nil {
		t.Fatal("expected pinned snapshot, got nil")
	}
	if p.Label != "snap1" {
		t.Errorf("expected label snap1, got %s", p.Label)
	}
	if p.Snapshot != s {
		t.Error("snapshot pointer mismatch")
	}
}

func TestPinBoard_Get_MissingLabel(t *testing.T) {
	pb := NewPinBoard()
	if pb.Get("missing") != nil {
		t.Fatal("expected nil for missing label")
	}
}

func TestPinBoard_Pin_Replaces(t *testing.T) {
	pb := NewPinBoard()
	pb.Pin("x", makePinnedSnapshot(1))
	s2 := makePinnedSnapshot(99)
	pb.Pin("x", s2)
	if pb.Len() != 1 {
		t.Fatalf("expected 1 pin after replace, got %d", pb.Len())
	}
	if pb.Get("x").Snapshot != s2 {
		t.Error("expected replaced snapshot")
	}
}

func TestPinBoard_Remove(t *testing.T) {
	pb := NewPinBoard()
	pb.Pin("del", makePinnedSnapshot(3))
	pb.Remove("del")
	if pb.Len() != 0 {
		t.Fatal("expected 0 pins after remove")
	}
}

func TestPinBoard_PinnedAt_Set(t *testing.T) {
	before := time.Now()
	pb := NewPinBoard()
	pb.Pin("ts", makePinnedSnapshot(2))
	after := time.Now()
	p := pb.Get("ts")
	if p.PinnedAt.Before(before) || p.PinnedAt.After(after) {
		t.Error("PinnedAt timestamp out of expected range")
	}
}
