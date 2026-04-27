package stats

// LabelConfig defines thresholds for assigning labels to buckets.
type LabelConfig struct {
	LargeKeyCount  int
	MediumKeyCount int
	LargeSizeBytes int64
	MediumSizeBytes int64
}

// DefaultLabelConfig returns a sensible default label configuration.
func DefaultLabelConfig() LabelConfig {
	return LabelConfig{
		LargeKeyCount:   10000,
		MediumKeyCount:  1000,
		LargeSizeBytes:  10 * 1024 * 1024, // 10 MB
		MediumSizeBytes: 1 * 1024 * 1024,  // 1 MB
	}
}

// BucketLabel represents a classification label for a bucket.
type BucketLabel string

const (
	LabelLarge  BucketLabel = "large"
	LabelMedium BucketLabel = "medium"
	LabelSmall  BucketLabel = "small"
	LabelEmpty  BucketLabel = "empty"
)

// LabeledBucket pairs a bucket name with its assigned label.
type LabeledBucket struct {
	Name  string
	Label BucketLabel
	Stats BucketStats
}

// LabelSnapshot assigns labels to all buckets in a snapshot based on config.
func LabelSnapshot(snap *Snapshot, cfg LabelConfig) []LabeledBucket {
	if snap == nil || snap.IsEmpty() {
		return nil
	}

	labeled := make([]LabeledBucket, 0, len(snap.Buckets))
	for name, bs := range snap.Buckets {
		labeled = append(labeled, LabeledBucket{
			Name:  name,
			Label: classifyBucket(bs, cfg),
			Stats: bs,
		})
	}
	return labeled
}

func classifyBucket(bs BucketStats, cfg LabelConfig) BucketLabel {
	if bs.KeyCount == 0 {
		return LabelEmpty
	}
	if bs.KeyCount >= cfg.LargeKeyCount || bs.Size >= cfg.LargeSizeBytes {
		return LabelLarge
	}
	if bs.KeyCount >= cfg.MediumKeyCount || bs.Size >= cfg.MediumSizeBytes {
		return LabelMedium
	}
	return LabelSmall
}
