package stats

import (
	"testing"
	"time"
)

func makeReplayHistory(n int, interval time.Duration) *History {
	h := NewHistory(100)
	base := time.Now().Add(-time.Duration(n) * interval)
	for i := 0; i < n; i++ {
		snap := NewSnapshot()
		snap.Timestamp = base.Add(time.Duration(i) * interval)
		snap.Buckets = map[string]BucketStats{
			"bucket": {Keys: i + 1, Size: int64((i + 1) * 100)},
		}
		h.Add(snap)
	}
	return h
}

func TestReplayHistory_NilHistory(t *testing.T) {
	_, err := ReplayHistory(nil, DefaultReplayOptions())
	if err == nil {
		t.Fatal("expected error for nil history")
	}
}

func TestReplayHistory_EmptyHistory(t *testing.T) {
	h := NewHistory(10)
	r, err := ReplayHistory(h, DefaultReplayOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Snapshots) != 0 {
		t.Errorf("expected 0 snapshots, got %d", len(r.Snapshots))
	}
}

func TestReplayHistory_AllSnapshots(t *testing.T) {
	h := makeReplayHistory(5, time.Second)
	r, err := ReplayHistory(h, DefaultReplayOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Snapshots) != 5 {
		t.Errorf("expected 5 snapshots, got %d", len(r.Snapshots))
	}
}

func TestReplayHistory_WithStep(t *testing.T) {
	h := makeReplayHistory(10, 500*time.Millisecond)
	opts := DefaultReplayOptions()
	opts.Step = time.Second
	r, err := ReplayHistory(h, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Snapshots) > 6 {
		t.Errorf("expected ≤6 snapshots with 1s step, got %d", len(r.Snapshots))
	}
}

func TestReplayHistory_TimeRange(t *testing.T) {
	h := makeReplayHistory(10, time.Second)
	all, _ := ReplayHistory(h, DefaultReplayOptions())
	mid := all.Snapshots[len(all.Snapshots)/2].Timestamp

	opts := DefaultReplayOptions()
	opts.End = mid
	r, err := ReplayHistory(h, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, s := range r.Snapshots {
		if s.Timestamp.After(mid) {
			t.Errorf("snapshot %v is after End %v", s.Timestamp, mid)
		}
	}
}

func TestReplayHistory_Chronological(t *testing.T) {
	h := makeReplayHistory(6, time.Second)
	r, _ := ReplayHistory(h, DefaultReplayOptions())
	for i := 1; i < len(r.Snapshots); i++ {
		if r.Snapshots[i].Timestamp.Before(r.Snapshots[i-1].Timestamp) {
			t.Errorf("snapshots not in chronological order at index %d", i)
		}
	}
}
