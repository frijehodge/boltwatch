package stats

// TrendDirection indicates whether a metric is increasing, decreasing, or stable.
type TrendDirection int

const (
	TrendStable   TrendDirection = 0
	TrendUp       TrendDirection = 1
	TrendDown     TrendDirection = -1
)

// TrendSymbol returns a visual symbol for the trend direction.
func (t TrendDirection) TrendSymbol() string {
	switch t {
	case TrendUp:
		return "↑"
	case TrendDown:
		return "↓"
	default:
		return "→"
	}
}

// BucketTrend holds trend information for a single bucket.
type BucketTrend struct {
	Bucket    string
	KeysTrend TrendDirection
	SizeTrend TrendDirection
}

// ComputeTrends analyzes a slice of snapshots and returns trend directions
// for each bucket based on the oldest and newest snapshots.
func ComputeTrends(history []Snapshot) []BucketTrend {
	if len(history) < 2 {
		return nil
	}

	oldest := history[0]
	newest := history[len(history)-1]

	trends := make([]BucketTrend, 0)

	for bucket, newStats := range newest.Buckets {
		trend := BucketTrend{Bucket: bucket}

		if oldStats, ok := oldest.Buckets[bucket]; ok {
			trend.KeysTrend = directionOf(int64(oldStats.KeyN), int64(newStats.KeyN))
			trend.SizeTrend = directionOf(oldStats.LeafInuse, newStats.LeafInuse)
		} else {
			// Bucket is new — treat as growing
			trend.KeysTrend = TrendUp
			trend.SizeTrend = TrendUp
		}

		trends = append(trends, trend)
	}

	return trends
}

func directionOf(old, new int64) TrendDirection {
	switch {
	case new > old:
		return TrendUp
	case new < old:
		return TrendDown
	default:
		return TrendStable
	}
}
