package stats

import (
	"strings"
	"testing"
	"time"
)

func makeFormatReplayResult(n int) *ReplayResult {
	snaps := make([]*Snapshot, n)
	base := time.Now()
	for i := 0; i < n; i++ {
		s := NewSnapshot()
		s.Timestamp = base.Add(time.Duration(i) * time.Second)
		s.Buckets = map[string]BucketStats{
			"alpha": {Keys: i + 1, Size: int64((i + 1) * 512)},
		}
		snaps[i] = s
	}
	result := &ReplayResult{Snapshots: snaps}
	if n > 0 {
		result.From = snaps[0].Timestamp
		result.To = snaps[n-1].Timestamp
	}
	return result
}

func TestFormatReplay_Nil(t *testing.T) {
	out := FormatReplay(nil)
	if !strings.Contains(out, "no snapshots") {
		t.Errorf("expected 'no snapshots', got: %s", out)
	}
}

func TestFormatReplay_Empty(t *testing.T) {
	out := FormatReplay(&ReplayResult{})
	if !strings.Contains(out, "no snapshots") {
		t.Errorf("expected 'no snapshots', got: %s", out)
	}
}

func TestFormatReplay_ContainsHeader(t *testing.T) {
	out := FormatReplay(makeFormatReplayResult(3))
	if !strings.Contains(out, "Replay:") {
		t.Errorf("expected 'Replay:' header, got: %s", out)
	}
}

func TestFormatReplay_ContainsSnapshotCount(t *testing.T) {
	out := FormatReplay(makeFormatReplayResult(4))
	if !strings.Contains(out, "4 snapshot(s)") {
		t.Errorf("expected '4 snapshot(s)', got: %s", out)
	}
}

func TestFormatReplayDiff_TooFew(t *testing.T) {
	out := FormatReplayDiff(makeFormatReplayResult(1))
	if !strings.Contains(out, "need at least 2") {
		t.Errorf("expected 'need at least 2', got: %s", out)
	}
}

func TestFormatReplayDiff_ContainsSteps(t *testing.T) {
	out := FormatReplayDiff(makeFormatReplayResult(3))
	if !strings.Contains(out, "step") {
		t.Errorf("expected step entries, got: %s", out)
	}
}
