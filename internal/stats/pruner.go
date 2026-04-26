package stats

import "time"

// PruneOptions configures snapshot pruning behavior.
type PruneOptions struct {
	// MaxAge removes snapshots older than this duration.
	MaxAge time.Duration
	// MaxCount keeps at most this many snapshots (0 = unlimited).
	MaxCount int
	// MinKeys removes buckets with fewer than this many keys.
	MinKeys int
}

// DefaultPruneOptions returns sensible pruning defaults.
func DefaultPruneOptions() PruneOptions {
	return PruneOptions{
		MaxAge:   24 * time.Hour,
		MaxCount: 100,
		MinKeys:  0,
	}
}

// PruneSnapshot removes bucket entries from a snapshot that do not meet
// the criteria defined in opts. It returns a new snapshot without
// modifying the original.
func PruneSnapshot(snap *Snapshot, opts PruneOptions) *Snapshot {
	if snap == nil {
		return nil
	}

	pruned := &Snapshot{
		Timestamp: snap.Timestamp,
		Buckets:   make(map[string]BucketStats),
	}

	for name, bs := range snap.Buckets {
		if opts.MinKeys > 0 && bs.KeyCount < opts.MinKeys {
			continue
		}
		pruned.Buckets[name] = bs
	}

	return pruned
}

// PruneHistory removes snapshots from a history slice that are older than
// opts.MaxAge and caps the result to opts.MaxCount (most recent kept).
// A nil or empty slice is returned as-is.
func PruneHistory(snapshots []*Snapshot, opts PruneOptions) []*Snapshot {
	if len(snapshots) == 0 {
		return snapshots
	}

	cutoff := time.Now().Add(-opts.MaxAge)

	var filtered []*Snapshot
	for _, s := range snapshots {
		if opts.MaxAge > 0 && s.Timestamp.Before(cutoff) {
			continue
		}
		filtered = append(filtered, s)
	}

	if opts.MaxCount > 0 && len(filtered) > opts.MaxCount {
		filtered = filtered[len(filtered)-opts.MaxCount:]
	}

	return filtered
}
