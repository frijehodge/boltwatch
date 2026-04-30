package stats

import (
	"strings"
	"testing"
)

func makeFormatMarkerSnapshot(keys int) *Snapshot {
	return makeMarkerSnapshot(keys)
}

func TestFormatMarkerBoard_Nil(t *testing.T) {
	out := FormatMarkerBoard(nil)
	if !strings.Contains(out, "No markers") {
		t.Fatalf("unexpected: %s", out)
	}
}

func TestFormatMarkerBoard_ContainsHeaders(t *testing.T) {
	mb := NewMarkerBoard()
	mb.Add("v2", makeFormatMarkerSnapshot(20), MarkerKindManual, "note")
	out := FormatMarkerBoard(mb)
	if !strings.Contains(out, "Label") {
		t.Fatalf("missing header: %s", out)
	}
	if !strings.Contains(out, "Kind") {
		t.Fatalf("missing kind header: %s", out)
	}
}

func TestFormatMarkerBoard_ContainsKind(t *testing.T) {
	mb := NewMarkerBoard()
	mb.Add("alert-marker", makeFormatMarkerSnapshot(5), MarkerKindAlert, "")
	out := FormatMarkerBoard(mb)
	if !strings.Contains(out, string(MarkerKindAlert)) {
		t.Fatalf("expected kind in output: %s", out)
	}
}

func TestFormatMarkerDiff_NilMarker(t *testing.T) {
	snap := makeFormatMarkerSnapshot(10)
	out := FormatMarkerDiff(snap, nil)
	if !strings.Contains(out, "not found") {
		t.Fatalf("unexpected: %s", out)
	}
}

func TestFormatMarkerDiff_NoChanges(t *testing.T) {
	snap := makeFormatMarkerSnapshot(10)
	mb := NewMarkerBoard()
	mb.Add("stable", snap, MarkerKindManual, "")
	m := mb.Get("stable")
	out := FormatMarkerDiff(snap, m)
	if !strings.Contains(out, "No changes") {
		t.Fatalf("expected no-changes message: %s", out)
	}
}

func TestFormatMarkerDiff_ContainsLabel(t *testing.T) {
	base := makeFormatMarkerSnapshot(10)
	mb := NewMarkerBoard()
	mb.Add("my-marker", base, MarkerKindManual, "")
	m := mb.Get("my-marker")
	current := makeFormatMarkerSnapshot(10)
	out := FormatMarkerDiff(current, m)
	if !strings.Contains(out, "my-marker") {
		t.Fatalf("expected label in diff output: %s", out)
	}
}
