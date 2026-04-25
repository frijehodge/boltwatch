package stats

import (
	"time"
)

// Sample represents a single point-in-time measurement of a bucket.
type Sample struct {
	Bucket    string
	Keys      int
	Size      int64
	Timestamp time.Time
}

// Sampler collects periodic samples from snapshots for trend analysis.
type Sampler struct {
	max     int
	samples []Sample
}

// NewSampler creates a Sampler that retains at most max samples per bucket.
func NewSampler(max int) *Sampler {
	if max <= 0 {
		max = 60
	}
	return &Sampler{max: max}
}

// Record extracts samples from a snapshot and appends them to the internal
// buffer, evicting the oldest entries when the capacity is exceeded.
func (s *Sampler) Record(snap *Snapshot) {
	if snap == nil || snap.IsEmpty() {
		return
	}
	for bucket, bs := range snap.Buckets {
		s.samples = append(s.samples, Sample{
			Bucket:    bucket,
			Keys:      bs.KeyN,
			Size:      bs.Size(),
			Timestamp: snap.Timestamp,
		})
	}
	if len(s.samples) > s.max {
		s.samples = s.samples[len(s.samples)-s.max:]
	}
}

// All returns a copy of all retained samples.
func (s *Sampler) All() []Sample {
	out := make([]Sample, len(s.samples))
	copy(out, s.samples)
	return out
}

// ForBucket returns samples for a specific bucket name.
func (s *Sampler) ForBucket(name string) []Sample {
	var out []Sample
	for _, sm := range s.samples {
		if sm.Bucket == name {
			out = append(out, sm)
		}
	}
	return out
}

// Len returns the total number of samples currently stored.
func (s *Sampler) Len() int {
	return len(s.samples)
}

// Reset clears all stored samples.
func (s *Sampler) Reset() {
	s.samples = s.samples[:0]
}
