package stats

// ReduceConfig controls how snapshots are reduced.
type ReduceConfig struct {
	// MaxBuckets limits the number of buckets retained. 0 means no limit.
	MaxBuckets int
	// SortByKeys sorts buckets by key count descending before reducing.
	SortByKeys bool
	// SortBySize sorts buckets by size descending before reducing.
	SortBySize bool
}

// DefaultReduceConfig returns sensible defaults for snapshot reduction.
func DefaultReduceConfig() ReduceConfig {
	return ReduceConfig{
		MaxBuckets: 20,
		SortByKeys: true,
		SortBySize: false,
	}
}

// ReduceSnapshot returns a new snapshot containing at most cfg.MaxBuckets
// buckets. When sorting is enabled the most significant buckets are kept.
// If snap is nil or empty it is returned as-is.
func ReduceSnapshot(snap *Snapshot, cfg ReduceConfig) *Snapshot {
	if snap == nil || snap.IsEmpty() {
		return snap
	}

	buckets := make([]BucketStats, len(snap.Buckets))
	copy(buckets, snap.Buckets)

	if cfg.SortByKeys {
		sortBucketsByKeys(buckets)
	} else if cfg.SortBySize {
		sortBucketsBySize(buckets)
	}

	if cfg.MaxBuckets > 0 && len(buckets) > cfg.MaxBuckets {
		buckets = buckets[:cfg.MaxBuckets]
	}

	return &Snapshot{
		Timestamp: snap.Timestamp,
		Buckets:   buckets,
	}
}

func sortBucketsByKeys(buckets []BucketStats) {
	for i := 1; i < len(buckets); i++ {
		for j := i; j > 0 && buckets[j].KeyCount > buckets[j-1].KeyCount; j-- {
			buckets[j], buckets[j-1] = buckets[j-1], buckets[j]
		}
	}
}

func sortBucketsBySize(buckets []BucketStats) {
	for i := 1; i < len(buckets); i++ {
		for j := i; j > 0 && buckets[j].Size > buckets[j-1].Size; j-- {
			buckets[j], buckets[j-1] = buckets[j-1], buckets[j]
		}
	}
}
