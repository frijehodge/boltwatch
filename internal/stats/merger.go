package stats

// MergeResult holds the combined output of merging two snapshots.
type MergeResult struct {
	Snapshot  *Snapshot
	Conflicts []string
}

// MergeStrategy controls how bucket conflicts are resolved.
type MergeStrategy int

const (
	// MergeStrategyPreferA keeps values from snapshot A on conflict.
	MergeStrategyPreferA MergeStrategy = iota
	// MergeStrategyPreferB keeps values from snapshot B on conflict.
	MergeStrategyPreferB
	// MergeStrategySumKeys sums key counts from both snapshots on conflict.
	MergeStrategySumKeys
)

// MergeSnapshots combines two snapshots into one according to the given strategy.
// Buckets present in only one snapshot are always included.
// Returns nil if both inputs are nil.
func MergeSnapshots(a, b *Snapshot, strategy MergeStrategy) *MergeResult {
	if a == nil && b == nil {
		return nil
	}
	if a == nil {
		return &MergeResult{Snapshot: b}
	}
	if b == nil {
		return &MergeResult{Snapshot: a}
	}

	merged := make(map[string]BucketStats)
	var conflicts []string

	// Add all buckets from A.
	for name, stat := range a.Buckets {
		merged[name] = stat
	}

	// Merge buckets from B.
	for name, statB := range b.Buckets {
		statA, exists := merged[name]
		if !exists {
			merged[name] = statB
			continue
		}

		conflicts = append(conflicts, name)

		switch strategy {
		case MergeStrategyPreferB:
			merged[name] = statB
		case MergeStrategySumKeys:
			merged[name] = BucketStats{
				Name:     statA.Name,
				KeyCount: statA.KeyCount + statB.KeyCount,
				Size:     statA.Size + statB.Size,
			}
		default: // MergeStrategyPreferA
			// keep statA as-is
		}
	}

	// Use the more recent timestamp.
	ts := a.Timestamp
	if b.Timestamp.After(ts) {
		ts = b.Timestamp
	}

	result := NewSnapshot(merged)
	result.Timestamp = ts

	return &MergeResult{
		Snapshot:  result,
		Conflicts: conflicts,
	}
}
