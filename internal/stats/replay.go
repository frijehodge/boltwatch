package stats

import (
	"errors"
	"sort"
	"time"
)

// ReplayOptions controls how a history replay is performed.
type ReplayOptions struct {
	// Start filters snapshots at or after this time. Zero means no lower bound.
	Start time.Time
	// End filters snapshots at or before this time. Zero means no upper bound.
	End time.Time
	// Step is the minimum duration between replayed snapshots. Zero means all.
	Step time.Duration
}

// DefaultReplayOptions returns sensible replay defaults.
func DefaultReplayOptions() ReplayOptions {
	return ReplayOptions{}
}

// ReplayResult holds an ordered slice of snapshots selected from history.
type ReplayResult struct {
	Snapshots []*Snapshot
	From      time.Time
	To        time.Time
}

// ReplayHistory replays snapshots from a History according to the given options.
// Snapshots are returned in chronological order.
func ReplayHistory(h *History, opts ReplayOptions) (*ReplayResult, error) {
	if h == nil {
		return nil, errors.New("replay: nil history")
	}

	all := h.All()
	if len(all) == 0 {
		return &ReplayResult{}, nil
	}

	// Ensure chronological order.
	sort.Slice(all, func(i, j int) bool {
		return all[i].Timestamp.Before(all[j].Timestamp)
	})

	var selected []*Snapshot
	var lastKept time.Time

	for _, snap := range all {
		if !opts.Start.IsZero() && snap.Timestamp.Before(opts.Start) {
			continue
		}
		if !opts.End.IsZero() && snap.Timestamp.After(opts.End) {
			continue
		}
		if opts.Step > 0 && !lastKept.IsZero() && snap.Timestamp.Sub(lastKept) < opts.Step {
			continue
		}
		selected = append(selected, snap)
		lastKept = snap.Timestamp
	}

	result := &ReplayResult{Snapshots: selected}
	if len(selected) > 0 {
		result.From = selected[0].Timestamp
		result.To = selected[len(selected)-1].Timestamp
	}
	return result, nil
}
