package stats

import (
	"strings"
	"testing"
	"time"
)

func makeFormatPinSnapshot(buckets map[string]int) *Snapshot {
	b := make(map[string]BucketStats, len(buckets))
	for name, keys := range buckets {
		b[name] = BucketStats{KeyCount: keys, Size: int64(keys * 128)}
	}
	return &Snapshot{Timestamp: time.Now(), Buckets: b}
}

func TestFormatPinBoard_Nil(t *testing.T) {
	out := FormatPinBoard(nil)
	if !strings.Contains(out, "No pinned") {
		t.Errorf("expected empty message, got: %s", out)
	}
}

func TestFormatPinBoard_Empty(t *testing.T) {
	out := FormatPinBoard(NewPinBoard())
	if !strings.Contains(out, "No pinned") {
		t.Errorf("expected empty message, got: %s", out)
	}
}

func TestFormatPinBoard_ContainsLabel(t *testing.T) {
	pb := NewPinBoard()
	pb.Pin("release-v1", makeFormatPinSnapshot(map[string]int{"users": 50}))
	out := FormatPinBoard(pb)
	if !strings.Contains(out, "release-v1") {
		t.Errorf("expected label in output, got:\n%s", out)
	}
}

func TestFormatPinBoard_ContainsHeaders(t *testing.T) {
	pb := NewPinBoard()
	pb.Pin("snap", makeFormatPinSnapshot(map[string]int{"a": 1}))
	out := FormatPinBoard(pb)
	for _, hdr := range []string{"LABEL", "PINNED AT", "BUCKETS", "TOTAL KEYS"} {
		if !strings.Contains(out, hdr) {
			t.Errorf("missing header %q in:\n%s", hdr, out)
		}
	}
}

func TestFormatPinDiff_NilInputs(t *testing.T) {
	out := FormatPinDiff(nil, nil)
	if !strings.Contains(out, "Cannot diff") {
		t.Errorf("expected nil message, got: %s", out)
	}
}

func TestFormatPinDiff_ShowsDelta(t *testing.T) {
	pb := NewPinBoard()
	pb.Pin("old", makeFormatPinSnapshot(map[string]int{"orders": 100}))
	pin := pb.Get("old")
	live := makeFormatPinSnapshot(map[string]int{"orders": 150})
	out := FormatPinDiff(pin, live)
	if !strings.Contains(out, "+50") {
		t.Errorf("expected +50 delta in output:\n%s", out)
	}
}

func TestFormatPinDiff_ShowsNewBucket(t *testing.T) {
	pb := NewPinBoard()
	pb.Pin("base", makeFormatPinSnapshot(map[string]int{"a": 10}))
	pin := pb.Get("base")
	live := makeFormatPinSnapshot(map[string]int{"a": 10, "b": 5})
	out := FormatPinDiff(pin, live)
	if !strings.Contains(out, "new since pin") {
		t.Errorf("expected new-bucket note in output:\n%s", out)
	}
}

func TestFormatPinDiff_ShowsRemovedBucket(t *testing.T) {
	pb := NewPinBoard()
	pb.Pin("base", makeFormatPinSnapshot(map[string]int{"a": 10, "gone": 7}))
	pin := pb.Get("base")
	live := makeFormatPinSnapshot(map[string]int{"a": 10})
	out := FormatPinDiff(pin, live)
	if !strings.Contains(out, "removed since pin") {
		t.Errorf("expected removed-bucket note in output:\n%s", out)
	}
}
