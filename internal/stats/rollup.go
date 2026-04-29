package stats

import "time"

// RollupPeriod defines the time window for a rollup.
type RollupPeriod string

const (
	RollupMinute RollupPeriod = "minute"
	RollupHour   RollupPeriod = "hour"
	RollupDay    RollupPeriod = "day"
)

// RollupEntry holds aggregated stats for a single time bucket.
type RollupEntry struct {
	PeriodStart time.Time
	PeriodEnd   time.Time
	Period      RollupPeriod
	AvgKeys     float64
	AvgSize     float64
	MaxKeys     int
	MaxSize     int64
	SampleCount int
}

// RollupResult maps bucket names to their rollup entries.
type RollupResult map[string]RollupEntry

// RollupSnapshots groups a slice of snapshots into time-bucketed rollups
// for a given period and computes per-bucket averages and maxima.
func RollupSnapshots(snapshots []*Snapshot, period RollupPeriod) RollupResult {
	if len(snapshots) == 0 {
		return RollupResult{}
	}

	type accumulator struct {
		sumKeys int
		sumSize int64
		maxKeys int
		maxSize int64
		count   int
	}

	accum := map[string]*accumulator{}

	var earliest, latest time.Time
	for i, snap := range snapshots {
		if snap == nil {
			continue
		}
		if i == 0 || snap.Timestamp.Before(earliest) {
			earliest = snap.Timestamp
		}
		if i == 0 || snap.Timestamp.After(latest) {
			latest = snap.Timestamp
		}
		for name, bs := range snap.Buckets {
			a, ok := accum[name]
			if !ok {
				a = &accumulator{}
				accum[name] = a
			}
			a.sumKeys += bs.KeyCount
			a.sumSize += bs.Size
			if bs.KeyCount > a.maxKeys {
				a.maxKeys = bs.KeyCount
			}
			if bs.Size > a.maxSize {
				a.maxSize = bs.Size
			}
			a.count++
		}
	}

	result := make(RollupResult, len(accum))
	for name, a := range accum {
		if a.count == 0 {
			continue
		}
		result[name] = RollupEntry{
			PeriodStart: earliest,
			PeriodEnd:   latest,
			Period:      period,
			AvgKeys:     float64(a.sumKeys) / float64(a.count),
			AvgSize:     float64(a.sumSize) / float64(a.count),
			MaxKeys:     a.maxKeys,
			MaxSize:     a.maxSize,
			SampleCount: a.count,
		}
	}
	return result
}
