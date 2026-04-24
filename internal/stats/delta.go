package stats

// Delta represents the change in statistics between two snapshots.
type Delta struct {
	Name       string
	KeyDelta   int
	SizeDelta  int64
}

// ComputeDeltas compares two slices of BucketStats and returns per-bucket
// deltas. Buckets present only in 'after' show their full values as growth;
// buckets removed from 'after' are omitted.
func ComputeDeltas(before, after []BucketStats) []Delta {
	prev := make(map[string]BucketStats, len(before))
	for _, b := range before {
		prev[b.Name] = b
	}

	deltas := make([]Delta, 0, len(after))
	for _, b := range after {
		d := Delta{Name: b.Name}
		if old, ok := prev[b.Name]; ok {
			d.KeyDelta = b.KeyCount - old.KeyCount
			d.SizeDelta = b.Size - old.Size
		} else {
			d.KeyDelta = b.KeyCount
			d.SizeDelta = b.Size
		}
		deltas = append(deltas, d)
	}
	return deltas
}

// HasChanges returns true if any delta contains a non-zero key or size change.
func HasChanges(deltas []Delta) bool {
	for _, d := range deltas {
		if d.KeyDelta != 0 || d.SizeDelta != 0 {
			return true
		}
	}
	return false
}
