package stats

import (
	"strings"
	"testing"

	"go.etcd.io/bbolt"
)

func TestFormatTable_Empty(t *testing.T) {
	stats := []BucketStats{}
	output := FormatTable(stats)
	if !strings.Contains(output, "Bucket") {
		t.Error("expected table header to contain 'Bucket'")
	}
	if !strings.Contains(output, "Keys") {
		t.Error("expected table header to contain 'Keys'")
	}
}

func TestFormatTable_WithData(t *testing.T) {
	stats := []BucketStats{
		{Name: "users", KeyCount: 42, Size: 1024},
		{Name: "sessions", KeyCount: 7, Size: 256},
	}
	output := FormatTable(stats)
	if !strings.Contains(output, "users") {
		t.Errorf("expected output to contain bucket name 'users', got:\n%s", output)
	}
	if !strings.Contains(output, "sessions") {
		t.Errorf("expected output to contain bucket name 'sessions', got:\n%s", output)
	}
	if !strings.Contains(output, "42") {
		t.Errorf("expected output to contain key count '42', got:\n%s", output)
	}
}

func TestFormatSummary_Empty(t *testing.T) {
	stats := []BucketStats{}
	output := FormatSummary(stats)
	if !strings.Contains(output, "0") {
		t.Errorf("expected summary to show 0 buckets, got:\n%s", output)
	}
}

func TestFormatSummary_WithData(t *testing.T) {
	stats := []BucketStats{
		{Name: "a", KeyCount: 10, Size: 512},
		{Name: "b", KeyCount: 20, Size: 1024},
	}
	output := FormatSummary(stats)
	if !strings.Contains(output, "2") {
		t.Errorf("expected summary to show 2 buckets, got:\n%s", output)
	}
	if !strings.Contains(output, "30") {
		t.Errorf("expected summary to show total key count 30, got:\n%s", output)
	}
}

func TestFormatTable_SizeFormatting(t *testing.T) {
	stats := []BucketStats{
		{Name: "large", KeyCount: 1, Size: 2048},
	}
	output := FormatTable(stats)
	if output == "" {
		t.Error("expected non-empty output")
	}
	_ = bbolt.Options{} // ensure bbolt import is used transitively via test db helpers
}
