package stats

// Delta represents the change between two snapshots for a single bucket.
type Delta struct {
	Bucket   string
	KeysDiff int
	SizeDiff int64
}

// ComputeDeltas compares two snapshots and returns a slice of Deltas
// for buckets that have changed between them.
func ComputeDeltas(prev, curr *Snapshot) []Delta {
	if prev == nil || curr == nil {
		return nil
	}

	var deltas []Delta

	// Check buckets in current snapshot
	for name, cs := range curr.Buckets {
		ps, ok := prev.Buckets[name]
		if !ok {
			// New bucket
			deltas = append(deltas, Delta{
				Bucket:   name,
				KeysDiff: cs.KeyN,
				SizeDiff: int64(cs.LeafInuse),
			})
			continue
		}
		keysDiff := cs.KeyN - ps.KeyN
		sizeDiff := int64(cs.LeafInuse) - int64(ps.LeafInuse)
		if keysDiff != 0 || sizeDiff != 0 {
			deltas = append(deltas, Delta{
				Bucket:   name,
				KeysDiff: keysDiff,
				SizeDiff: sizeDiff,
			})
		}
	}

	// Check for removed buckets
	for name, ps := range prev.Buckets {
		if _, ok := curr.Buckets[name]; !ok {
			deltas = append(deltas, Delta{
				Bucket:   name,
				KeysDiff: -ps.KeyN,
				SizeDiff: -int64(ps.LeafInuse),
			})
		}
	}

	return deltas
}

// HasChanges returns true if any deltas are present.
func HasChanges(deltas []Delta) bool {
	return len(deltas) > 0
}
