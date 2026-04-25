package stats

import "sort"

// BucketRank holds a bucket name and its rank score.
type BucketRank struct {
	Bucket string
	Keys   int
	Size   int64
}

// AggregateResult holds the output of an aggregation pass.
type AggregateResult struct {
	TopByKeys []BucketRank
	TopBySize []BucketRank
	TotalKeys int
	TotalSize int64
	Buckets   int
}

// Aggregate computes ranked bucket statistics from a snapshot.
// topN controls how many top buckets are returned; 0 means all.
func Aggregate(snap *Snapshot, topN int) *AggregateResult {
	if snap == nil || snap.IsEmpty() {
		return &AggregateResult{}
	}

	ranks := make([]BucketRank, 0, len(snap.Buckets))
	for name, bs := range snap.Buckets {
		ranks = append(ranks, BucketRank{
			Bucket: name,
			Keys:   bs.KeyN,
			Size:   bs.LeafInuse,
		})
	}

	byKeys := make([]BucketRank, len(ranks))
	copy(byKeys, ranks)
	sort.Slice(byKeys, func(i, j int) bool {
		return byKeys[i].Keys > byKeys[j].Keys
	})

	bySize := make([]BucketRank, len(ranks))
	copy(bySize, ranks)
	sort.Slice(bySize, func(i, j int) bool {
		return bySize[i].Size > bySize[j].Size
	})

	if topN > 0 && topN < len(byKeys) {
		byKeys = byKeys[:topN]
		bySize = bySize[:topN]
	}

	return &AggregateResult{
		TopByKeys: byKeys,
		TopBySize: bySize,
		TotalKeys: snap.TotalKeys(),
		TotalSize: snap.TotalSize(),
		Buckets:   len(snap.Buckets),
	}
}
