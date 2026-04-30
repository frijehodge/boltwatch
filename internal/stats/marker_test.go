package stats

import (
	"strings"
	"testing"
	"time"
)

func makeMarkerSnapshot(keys int) *Snapshot {
	bs := map[string]BucketStats{
		"bucket-a": {Name: "bucket-a", KeyCount: keys, Size: int64(keys * 64)},
	}
	snap := NewSnapshot(bs)
	snap.Timestamp = time.Now()
	return snap
}

func TestNewMarkerBoard_Empty(t *testing.T) {
	mb := NewMarkerBoard()
	if mb.Len() != 0 {
		t.Fatalf("expected 0, got %d", mb.Len())
	}
}

func TestMarkerBoard_Add_NilSnapshot(t *testing.T) {
	mb := NewMarkerBoard()
	mb.Add("x", nil, MarkerKindManual, "")
	if mb.Len() != 0 {
		t.Fatal("nil snapshot should not be stored")
	}
}

func TestMarkerBoard_Add_StoresMarker(t *testing.T) {
	mb := NewMarkerBoard()
	snap := makeMarkerSnapshot(10)
	mb.Add("release-1", snap, MarkerKindManual, "first release")
	if mb.Len() != 1 {
		t.Fatalf("expected 1, got %d", mb.Len())
	}
}

func TestMarkerBoard_Get_ReturnsMarker(t *testing.T) {
	mb := NewMarkerBoard()
	snap := makeMarkerSnapshot(5)
	mb.Add("v1", snap, MarkerKindAutomatic, "auto")
	m := mb.Get("v1")
	if m == nil {
		t.Fatal("expected marker, got nil")
	}
	if m.Label != "v1" {
		t.Fatalf("expected label v1, got %s", m.Label)
	}
}

func TestMarkerBoard_Remove(t *testing.T) {
	mb := NewMarkerBoard()
	snap := makeMarkerSnapshot(3)
	mb.Add("temp", snap, MarkerKindAlert, "")
	mb.Remove("temp")
	if mb.Len() != 0 {
		t.Fatal("expected empty board after remove")
	}
}

func TestFormatMarkerBoard_Empty(t *testing.T) {
	mb := NewMarkerBoard()
	out := FormatMarkerBoard(mb)
	if !strings.Contains(out, "No markers") {
		t.Fatalf("unexpected output: %s", out)
	}
}

func TestFormatMarkerBoard_ContainsLabel(t *testing.T) {
	mb := NewMarkerBoard()
	snap := makeMarkerSnapshot(7)
	mb.Add("checkpoint", snap, MarkerKindManual, "test note")
	out := FormatMarkerBoard(mb)
	if !strings.Contains(out, "checkpoint") {
		t.Fatalf("expected label in output: %s", out)
	}
}

func TestCompareToMarker_NilInputs(t *testing.T) {
	if CompareToMarker(nil, nil) != nil {
		t.Fatal("expected nil")
	}
}

func TestHasMarkerChanges_NoChanges(t *testing.T) {
	snap := makeMarkerSnapshot(10)
	mb := NewMarkerBoard()
	mb.Add("base", snap, MarkerKindManual, "")
	m := mb.Get("base")
	if HasMarkerChanges(snap, m) {
		t.Fatal("expected no changes")
	}
}

func TestMarkerDriftSummary_NoDrift(t *testing.T) {
	snap := makeMarkerSnapshot(10)
	mb := NewMarkerBoard()
	mb.Add("base", snap, MarkerKindManual, "")
	m := mb.Get("base")
	summary := MarkerDriftSummary(snap, m)
	if summary != "no drift" {
		t.Fatalf("expected 'no drift', got %q", summary)
	}
}
