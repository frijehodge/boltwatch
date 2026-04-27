package stats

// ScoreWeights controls how each metric contributes to the overall bucket score.
type ScoreWeights struct {
	KeysWeight float64
	SizeWeight float64
	GrowthWeight float64
}

// DefaultScoreWeights returns sensible default weights for bucket scoring.
func DefaultScoreWeights() ScoreWeights {
	return ScoreWeights{
		KeysWeight:   0.5,
		SizeWeight:   0.3,
		GrowthWeight: 0.2,
	}
}

// BucketScore holds the computed score and contributing factors for a bucket.
type BucketScore struct {
	Bucket      string
	Score       float64
	KeysScore   float64
	SizeScore   float64
	GrowthScore float64
}

// ScoreSnapshot computes a normalized score for each bucket in the snapshot.
// Scores are in the range [0.0, 1.0] relative to the max observed values.
// An optional Trend map (bucket -> direction) boosts the growth component.
func ScoreSnapshot(snap *Snapshot, weights ScoreWeights, trends map[string]string) []BucketScore {
	if snap == nil || snap.IsEmpty() {
		return nil
	}

	buckets := snap.Buckets()
	if len(buckets) == 0 {
		return nil
	}

	// Find max values for normalization.
	var maxKeys, maxSize float64
	for _, b := range buckets {
		if k := float64(b.KeyN); k > maxKeys {
			maxKeys = k
		}
		if s := float64(b.LeafInuse); s > maxSize {
			maxSize = s
		}
	}

	scores := make([]BucketScore, 0, len(buckets))
	for name, b := range buckets {
		var ks, ss, gs float64

		if maxKeys > 0 {
			ks = float64(b.KeyN) / maxKeys
		}
		if maxSize > 0 {
			ss = float64(b.LeafInuse) / maxSize
		}
		if trends != nil {
			switch trends[name] {
			case "up":
				gs = 1.0
			case "stable":
				gs = 0.5
			case "down":
				gs = 0.0
			}
		}

		total := ks*weights.KeysWeight + ss*weights.SizeWeight + gs*weights.GrowthWeight
		scores = append(scores, BucketScore{
			Bucket:      name,
			Score:       total,
			KeysScore:   ks,
			SizeScore:   ss,
			GrowthScore: gs,
		})
	}

	return scores
}
