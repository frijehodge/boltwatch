package stats

import "time"

// Baseline holds a reference snapshot used to compare current state against
// an established baseline for drift detection.
type Baseline struct {
	snapshot  *Snapshot
	CreatedAt time.Time
	Label     string
}

// NewBaseline creates a new Baseline from the given snapshot.
func NewBaseline(snap *Snapshot, label string) *Baseline {
	if snap == nil {
		return nil
	}
	return &Baseline{
		snapshot:  snap,
		CreatedAt: time.Now(),
		Label:     label,
	}
}

// Snapshot returns the underlying reference snapshot.
func (b *Baseline) Snapshot() *Snapshot {
	if b == nil {
		return nil
	}
	return b.snapshot
}

// DriftResult describes how a bucket has drifted from the baseline.
type DriftResult struct {
	Bucket   string
	KeyDrift int    // positive = grown, negative = shrunk
	SizeDrift int64 // positive = grown, negative = shrunk
	IsNew    bool
	IsGone   bool
}

// ComputeDrift compares a current snapshot against the baseline and returns
// per-bucket drift results. Returns nil if either input is nil.
func ComputeDrift(base *Baseline, current *Snapshot) []DriftResult {
	if base == nil || base.snapshot == nil || current == nil {
		return nil
	}

	baseStats := base.snapshot.Buckets()
	currStats := current.Buckets()

	seen := make(map[string]bool)
	var results []DriftResult

	for name, cur := range currStats {
		seen[name] = true
		if bst, ok := baseStats[name]; ok {
			results = append(results, DriftResult{
				Bucket:    name,
				KeyDrift:  cur.KeyN - bst.KeyN,
				SizeDrift: int64(cur.LeafInuse) - int64(bst.LeafInuse),
			})
		} else {
			results = append(results, DriftResult{
				Bucket: name,
				IsNew:  true,
			})
		}
	}

	for name := range baseStats {
		if !seen[name] {
			results = append(results, DriftResult{
				Bucket: name,
				IsGone: true,
			})
		}
	}

	return results
}
