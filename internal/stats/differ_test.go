package stats

import (
	"strings"
	"testing"
)

func makeDiffSnapshot(buckets map[string]BucketStats) *Snapshot {
	s := NewSnapshot()
	for name, bs := range buckets {
		s.Buckets[name] = bs
	}
	return s
}

func TestDiffSnapshots_BothNil(t *testing.T) {
	result := DiffSnapshots(nil, nil)
	if !result.IsEmpty() {
		t.Error("expected empty diff for nil inputs")
	}
}

func TestDiffSnapshots_NilBefore(t *testing.T) {
	after := makeDiffSnapshot(map[string]BucketStats{
		"alpha": {KeyCount: 10, Size: 100},
	})
	result := DiffSnapshots(nil, after)
	if len(result.Added) != 1 || result.Added[0] != "alpha" {
		t.Errorf("expected alpha in Added, got %v", result.Added)
	}
}

func TestDiffSnapshots_NilAfter(t *testing.T) {
	before := makeDiffSnapshot(map[string]BucketStats{
		"beta": {KeyCount: 5, Size: 50},
	})
	result := DiffSnapshots(before, nil)
	if len(result.Removed) != 1 || result.Removed[0] != "beta" {
		t.Errorf("expected beta in Removed, got %v", result.Removed)
	}
}

func TestDiffSnapshots_NoChanges(t *testing.T) {
	buckets := map[string]BucketStats{"gamma": {KeyCount: 3, Size: 30}}
	before := makeDiffSnapshot(buckets)
	after := makeDiffSnapshot(buckets)
	result := DiffSnapshots(before, after)
	if !result.IsEmpty() {
		t.Errorf("expected no diff, got added=%v removed=%v changed=%v", result.Added, result.Removed, result.Changed)
	}
	if len(result.Unchanged) != 1 {
		t.Errorf("expected 1 unchanged bucket, got %d", len(result.Unchanged))
	}
}

func TestDiffSnapshots_ChangedBucket(t *testing.T) {
	before := makeDiffSnapshot(map[string]BucketStats{"delta": {KeyCount: 2, Size: 20}})
	after := makeDiffSnapshot(map[string]BucketStats{"delta": {KeyCount: 8, Size: 80}})
	result := DiffSnapshots(before, after)
	if len(result.Changed) != 1 || result.Changed[0] != "delta" {
		t.Errorf("expected delta in Changed, got %v", result.Changed)
	}
}

func TestFormatDiff_Empty(t *testing.T) {
	d := &DiffResult{}
	out := FormatDiff(d)
	if out != "No differences found." {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormatDiff_WithChanges(t *testing.T) {
	d := &DiffResult{
		Added:   []string{"newbucket"},
		Removed: []string{"oldbucket"},
		Changed: []string{"modbucket"},
	}
	out := FormatDiff(d)
	if !strings.Contains(out, "+ newbucket") {
		t.Error("expected added bucket in output")
	}
	if !strings.Contains(out, "- oldbucket") {
		t.Error("expected removed bucket in output")
	}
	if !strings.Contains(out, "~ modbucket") {
		t.Error("expected changed bucket in output")
	}
}
