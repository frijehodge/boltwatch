package stats

import "sort"

// RankEntry holds a bucket name and its rank score.
type RankEntry struct {
	Bucket string
	Keys   int
	Size   int64
	Rank   float64
}

// RankConfig controls how buckets are scored and ranked.
type RankConfig struct {
	// KeyWeight is the weight applied to normalized key count (0.0–1.0).
	KeyWeight float64
	// SizeWeight is the weight applied to normalized size (0.0–1.0).
	SizeWeight float64
	// TopN limits results; 0 means return all.
	TopN int
}

// DefaultRankConfig returns a balanced ranking configuration.
func DefaultRankConfig() RankConfig {
	return RankConfig{
		KeyWeight:  0.5,
		SizeWeight: 0.5,
		TopN:       0,
	}
}

// RankBuckets scores and ranks buckets from a snapshot using the given config.
// Buckets are ranked by a weighted combination of key count and size.
func RankBuckets(snap *Snapshot, cfg RankConfig) []RankEntry {
	if snap == nil || snap.IsEmpty() {
		return nil
	}

	var maxKeys int
	var maxSize int64

	for _, bs := range snap.Buckets {
		if bs.Keys > maxKeys {
			maxKeys = bs.Keys
		}
		if bs.Size > maxSize {
			maxSize = bs.Size
		}
	}

	entries := make([]RankEntry, 0, len(snap.Buckets))
	for name, bs := range snap.Buckets {
		var normKeys, normSize float64
		if maxKeys > 0 {
			normKeys = float64(bs.Keys) / float64(maxKeys)
		}
		if maxSize > 0 {
			normSize = float64(bs.Size) / float64(maxSize)
		}
		rank := cfg.KeyWeight*normKeys + cfg.SizeWeight*normSize
		entries = append(entries, RankEntry{
			Bucket: name,
			Keys:   bs.Keys,
			Size:   bs.Size,
			Rank:   rank,
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Rank != entries[j].Rank {
			return entries[i].Rank > entries[j].Rank
		}
		return entries[i].Bucket < entries[j].Bucket
	})

	if cfg.TopN > 0 && cfg.TopN < len(entries) {
		return entries[:cfg.TopN]
	}
	return entries
}
