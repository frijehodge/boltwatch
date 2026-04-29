package stats

import (
	"strings"
	"testing"
	"time"
)

func makeProfilerSnapshot() *Snapshot {
	snap := NewSnapshot()
	snap.Buckets = []BucketStats{
		{Name: "users", Keys: 500, Size: 51200},
		{Name: "sessions", Keys: 120, Size: 9600},
		{Name: "logs", Keys: 10, Size: 800},
	}
	snap.CapturedAt = time.Now()
	return snap
}

func TestProfileSnapshot_Nil(t *testing.T) {
	p := ProfileSnapshot(nil, DefaultScoreWeights(), DefaultLabelConfig())
	if p == nil {
		t.Fatal("expected non-nil profile for nil snapshot")
	}
	if len(p.Buckets) != 0 {
		t.Errorf("expected 0 buckets, got %d", len(p.Buckets))
	}
}

func TestProfileSnapshot_Empty(t *testing.T) {
	snap := NewSnapshot()
	p := ProfileSnapshot(snap, DefaultScoreWeights(), DefaultLabelConfig())
	if len(p.Buckets) != 0 {
		t.Errorf("expected 0 buckets for empty snapshot, got %d", len(p.Buckets))
	}
}

func TestProfileSnapshot_BucketCount(t *testing.T) {
	snap := makeProfilerSnapshot()
	p := ProfileSnapshot(snap, DefaultScoreWeights(), DefaultLabelConfig())
	if len(p.Buckets) != 3 {
		t.Errorf("expected 3 bucket profiles, got %d", len(p.Buckets))
	}
}

func TestProfileSnapshot_TotalsCorrect(t *testing.T) {
	snap := makeProfilerSnapshot()
	p := ProfileSnapshot(snap, DefaultScoreWeights(), DefaultLabelConfig())
	if p.TotalKeys != 630 {
		t.Errorf("expected TotalKeys=630, got %d", p.TotalKeys)
	}
	expectedSize := int64(51200 + 9600 + 800)
	if p.TotalSize != expectedSize {
		t.Errorf("expected TotalSize=%d, got %d", expectedSize, p.TotalSize)
	}
}

func TestProfileSnapshot_SortedByScore(t *testing.T) {
	snap := makeProfilerSnapshot()
	p := ProfileSnapshot(snap, DefaultScoreWeights(), DefaultLabelConfig())
	for i := 1; i < len(p.Buckets); i++ {
		if p.Buckets[i].Score > p.Buckets[i-1].Score {
			t.Errorf("profiles not sorted by score descending at index %d", i)
		}
	}
}

func TestProfileSnapshot_AvgKeySizeNonZero(t *testing.T) {
	snap := makeProfilerSnapshot()
	p := ProfileSnapshot(snap, DefaultScoreWeights(), DefaultLabelConfig())
	for _, b := range p.Buckets {
		if b.Keys > 0 && b.Size > 0 && b.AvgKeySize == 0 {
			t.Errorf("expected non-zero AvgKeySize for bucket %s", b.Name)
		}
	}
}

func TestFormatProfile_Nil(t *testing.T) {
	out := FormatProfile(nil)
	if !strings.Contains(out, "No profile") {
		t.Errorf("expected fallback message for nil profile, got: %s", out)
	}
}

func TestFormatProfile_ContainsHeaders(t *testing.T) {
	snap := makeProfilerSnapshot()
	p := ProfileSnapshot(snap, DefaultScoreWeights(), DefaultLabelConfig())
	out := FormatProfile(p)
	for _, header := range []string{"Bucket", "Keys", "Size", "Score", "Label"} {
		if !strings.Contains(out, header) {
			t.Errorf("expected header %q in output", header)
		}
	}
}

func TestFormatProfile_ContainsBucketNames(t *testing.T) {
	snap := makeProfilerSnapshot()
	p := ProfileSnapshot(snap, DefaultScoreWeights(), DefaultLabelConfig())
	out := FormatProfile(p)
	for _, name := range []string{"users", "sessions", "logs"} {
		if !strings.Contains(out, name) {
			t.Errorf("expected bucket name %q in output", name)
		}
	}
}
