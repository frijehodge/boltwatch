package stats

import (
	"strings"
	"testing"
	"time"
)

func makeFormatLabelSnapshot() []LabeledBucket {
	snap := &Snapshot{
		Timestamp: time.Now(),
		Buckets: map[string]BucketStats{
			"alpha": {KeyCount: 0, Size: 0},
			"beta":  {KeyCount: 500, Size: 256 * 1024},
			"gamma": {KeyCount: 5000, Size: 2 * 1024 * 1024},
			"delta": {KeyCount: 50000, Size: 20 * 1024 * 1024},
		},
	}
	return LabelSnapshot(snap, DefaultLabelConfig())
}

func TestFormatLabels_Empty(t *testing.T) {
	out := FormatLabels(nil)
	if !strings.Contains(out, "No buckets") {
		t.Errorf("expected 'No buckets' message, got: %s", out)
	}
}

func TestFormatLabels_ContainsHeader(t *testing.T) {
	labeled := makeFormatLabelSnapshot()
	out := FormatLabels(labeled)
	if !strings.Contains(out, "Bucket") {
		t.Errorf("expected header 'Bucket' in output, got: %s", out)
	}
	if !strings.Contains(out, "Label") {
		t.Errorf("expected header 'Label' in output, got: %s", out)
	}
}

func TestFormatLabels_ContainsBucketNames(t *testing.T) {
	labeled := makeFormatLabelSnapshot()
	out := FormatLabels(labeled)
	for _, name := range []string{"alpha", "beta", "gamma", "delta"} {
		if !strings.Contains(out, name) {
			t.Errorf("expected bucket name %q in output", name)
		}
	}
}

func TestFormatLabels_ContainsLabelValues(t *testing.T) {
	labeled := makeFormatLabelSnapshot()
	out := FormatLabels(labeled)
	for _, lbl := range []string{"empty", "small", "medium", "large"} {
		if !strings.Contains(out, lbl) {
			t.Errorf("expected label %q in output", lbl)
		}
	}
}

func TestFormatLabelSummary_Empty(t *testing.T) {
	out := FormatLabelSummary(nil)
	if !strings.Contains(out, "No buckets") {
		t.Errorf("expected 'No buckets' message, got: %s", out)
	}
}

func TestFormatLabelSummary_Counts(t *testing.T) {
	labeled := makeFormatLabelSnapshot()
	out := FormatLabelSummary(labeled)
	if !strings.Contains(out, "large=1") {
		t.Errorf("expected large=1 in summary, got: %s", out)
	}
	if !strings.Contains(out, "empty=1") {
		t.Errorf("expected empty=1 in summary, got: %s", out)
	}
}
