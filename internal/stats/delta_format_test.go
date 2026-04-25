package stats

import (
	"strings"
	"testing"
)

func TestFormatDeltas_NoChanges(t *testing.T) {
	deltas := map[string]Delta{}
	result := FormatDeltas(deltas)
	if result != "No changes detected." {
		t.Errorf("expected no-change message, got: %q", result)
	}
}

func TestFormatDeltas_WithKeyIncrease(t *testing.T) {
	deltas := map[string]Delta{
		"users": {KeysDelta: 5, SizeDelta: 1024},
	}
	result := FormatDeltas(deltas)
	if !strings.Contains(result, "users") {
		t.Errorf("expected bucket name in output, got: %q", result)
	}
	if !strings.Contains(result, "+5") {
		t.Errorf("expected positive keys delta, got: %q", result)
	}
}

func TestFormatDeltas_WithKeyDecrease(t *testing.T) {
	deltas := map[string]Delta{
		"sessions": {KeysDelta: -3, SizeDelta: -512},
	}
	result := FormatDeltas(deltas)
	if !strings.Contains(result, "-3") {
		t.Errorf("expected negative keys delta, got: %q", result)
	}
}

func TestFormatDeltas_MultipleBuckets(t *testing.T) {
	deltas := map[string]Delta{
		"alpha": {KeysDelta: 2, SizeDelta: 256},
		"beta":  {KeysDelta: -1, SizeDelta: -128},
	}
	result := FormatDeltas(deltas)
	if !strings.Contains(result, "alpha") {
		t.Errorf("expected alpha bucket in output")
	}
	if !strings.Contains(result, "beta") {
		t.Errorf("expected beta bucket in output")
	}
}

func TestFormatDeltaBytes_Positive(t *testing.T) {
	out := formatDeltaBytes(2048)
	if !strings.HasPrefix(out, "+") {
		t.Errorf("expected positive prefix, got: %q", out)
	}
}

func TestFormatDeltaBytes_Negative(t *testing.T) {
	out := formatDeltaBytes(-1024)
	if !strings.HasPrefix(out, "-") {
		t.Errorf("expected negative prefix, got: %q", out)
	}
}

func TestFormatDeltaBytes_Zero(t *testing.T) {
	out := formatDeltaBytes(0)
	if out != "+0 B" {
		t.Errorf("expected zero representation, got: %q", out)
	}
}
