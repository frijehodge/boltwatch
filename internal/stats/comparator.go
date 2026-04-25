package stats

import "sort"

// BucketComparison holds the comparison result between two snapshots for a single bucket.
type BucketComparison struct {
	Bucket    string
	KeysA     int
	KeysB     int
	SizeA     int64
	SizeB     int64
	KeysDiff  int
	SizeDiff  int64
	IsNew     bool
	IsRemoved bool
}

// CompareResult holds all bucket comparisons between two snapshots.
type CompareResult struct {
	Comparisons []BucketComparison
	AddedCount  int
	RemovedCount int
	ChangedCount int
}

// CompareSnapshots compares two snapshots and returns a structured diff.
// If either snapshot is nil, an empty CompareResult is returned.
func CompareSnapshots(a, b *Snapshot) CompareResult {
	if a == nil || b == nil {
		return CompareResult{}
	}

	seen := make(map[string]bool)
	var comparisons []BucketComparison

	for name, bsB := range b.Buckets {
		seen[name] = true
		comp := BucketComparison{
			Bucket: name,
			KeysB:  bsB.KeyCount,
			SizeB:  bsB.Size,
		}
		if bsA, ok := a.Buckets[name]; ok {
			comp.KeysA = bsA.KeyCount
			comp.SizeA = bsA.Size
			comp.KeysDiff = bsB.KeyCount - bsA.KeyCount
			comp.SizeDiff = bsB.Size - bsA.Size
			if comp.KeysDiff != 0 || comp.SizeDiff != 0 {
				// changed
			}
		} else {
			comp.IsNew = true
			comp.KeysDiff = bsB.KeyCount
			comp.SizeDiff = bsB.Size
		}
		comparisons = append(comparisons, comp)
	}

	for name, bsA := range a.Buckets {
		if seen[name] {
			continue
		}
		comparisons = append(comparisons, BucketComparison{
			Bucket:    name,
			KeysA:     bsA.KeyCount,
			SizeA:     bsA.Size,
			KeysDiff:  -bsA.KeyCount,
			SizeDiff:  -bsA.Size,
			IsRemoved: true,
		})
	}

	sort.Slice(comparisons, func(i, j int) bool {
		return comparisons[i].Bucket < comparisons[j].Bucket
	})

	result := CompareResult{Comparisons: comparisons}
	for _, c := range comparisons {
		switch {
		case c.IsNew:
			result.AddedCount++
		case c.IsRemoved:
			result.RemovedCount++
		case c.KeysDiff != 0 || c.SizeDiff != 0:
			result.ChangedCount++
		}
	}
	return result
}
