package stats

import (
	"bytes"
	"strings"
	"testing"
)

func TestReporter_Report_Empty(t *testing.T) {
	var buf bytes.Buffer
	r := NewReporter(&buf, false)
	s := NewSnapshot([]BucketStats{})
	r.Report(s)

	out := buf.String()
	if !strings.Contains(out, "No buckets found") {
		t.Errorf("expected 'No buckets found' in output, got: %s", out)
	}
}

func TestReporter_Report_WithBuckets(t *testing.T) {
	var buf bytes.Buffer
	r := NewReporter(&buf, false)
	buckets := []BucketStats{
		{Name: "users", KeyN: 42, LeafSize: 2048},
	}
	s := NewSnapshot(buckets)
	r.Report(s)

	out := buf.String()
	if !strings.Contains(out, "users") {
		t.Errorf("expected bucket name 'users' in output, got: %s", out)
	}
}

func TestReporter_Report_Verbose_IncludesSummary(t *testing.T) {
	var buf bytes.Buffer
	r := NewReporter(&buf, true)
	buckets := []BucketStats{
		{Name: "logs", KeyN: 100, LeafSize: 4096},
	}
	s := NewSnapshot(buckets)
	r.Report(s)

	out := buf.String()
	if !strings.Contains(out, "logs") {
		t.Errorf("expected bucket name in verbose output, got: %s", out)
	}
}

func TestReporter_ReportDelta_NoChanges_Silent(t *testing.T) {
	var buf bytes.Buffer
	r := NewReporter(&buf, false)
	buckets := []BucketStats{{Name: "a", KeyN: 5, LeafSize: 512}}
	prev := NewSnapshot(buckets)
	curr := NewSnapshot(buckets)
	r.ReportDelta(prev, curr)

	if buf.Len() != 0 {
		t.Errorf("expected no output for unchanged snapshot (non-verbose), got: %s", buf.String())
	}
}

func TestReporter_ReportDelta_WithChanges(t *testing.T) {
	var buf bytes.Buffer
	r := NewReporter(&buf, false)
	prev := NewSnapshot([]BucketStats{{Name: "orders", KeyN: 10, LeafSize: 1024}})
	curr := NewSnapshot([]BucketStats{{Name: "orders", KeyN: 20, LeafSize: 2048}})
	r.ReportDelta(prev, curr)

	out := buf.String()
	if !strings.Contains(out, "orders") {
		t.Errorf("expected 'orders' in delta output, got: %s", out)
	}
	if !strings.Contains(out, "+10") {
		t.Errorf("expected key delta '+10' in output, got: %s", out)
	}
}

func TestReporter_ReportDelta_NoChanges_Verbose(t *testing.T) {
	var buf bytes.Buffer
	r := NewReporter(&buf, true)
	buckets := []BucketStats{{Name: "cache", KeyN: 3, LeafSize: 256}}
	prev := NewSnapshot(buckets)
	curr := NewSnapshot(buckets)
	r.ReportDelta(prev, curr)

	out := buf.String()
	if !strings.Contains(out, "no changes") {
		t.Errorf("expected 'no changes' message in verbose mode, got: %s", out)
	}
}
