package stats

// GrowthClass represents the growth classification of a bucket.
type GrowthClass string

const (
	GrowthClassStable   GrowthClass = "stable"
	GrowthClassGrowing  GrowthClass = "growing"
	GrowthClassShrinking GrowthClass = "shrinking"
	GrowthClassNew      GrowthClass = "new"
	GrowthClassRemoved  GrowthClass = "removed"
)

// BucketClassification holds the growth class and supporting metadata for a bucket.
type BucketClassification struct {
	Bucket     string
	Class      GrowthClass
	KeyDelta   int
	SizeDelta  int64
	GrowthRate float64 // keys per second, approximate
}

// ClassifyGrowth classifies each bucket based on delta information between two snapshots.
// Returns a map of bucket name to BucketClassification.
func ClassifyGrowth(prev, curr *Snapshot) map[string]BucketClassification {
	result := make(map[string]BucketClassification)
	if prev == nil || curr == nil {
		return result
	}

	elapsed := curr.Timestamp.Sub(prev.Timestamp).Seconds()
	if elapsed <= 0 {
		elapsed = 1
	}

	for name, cs := range curr.Buckets {
		ps, existed := prev.Buckets[name]
		if !existed {
			result[name] = BucketClassification{
				Bucket: name,
				Class:  GrowthClassNew,
			}
			continue
		}
		keyDelta := cs.KeyCount - ps.KeyCount
		sizeDelta := cs.Size - ps.Size
		growthRate := float64(keyDelta) / elapsed

		class := GrowthClassStable
		if keyDelta > 0 {
			class = GrowthClassGrowing
		} else if keyDelta < 0 {
			class = GrowthClassShrinking
		}

		result[name] = BucketClassification{
			Bucket:     name,
			Class:      class,
			KeyDelta:   keyDelta,
			SizeDelta:  sizeDelta,
			GrowthRate: growthRate,
		}
	}

	for name := range prev.Buckets {
		if _, exists := curr.Buckets[name]; !exists {
			result[name] = BucketClassification{
				Bucket: name,
				Class:  GrowthClassRemoved,
			}
		}
	}

	return result
}
