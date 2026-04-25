package stats

import (
	"time"
)

// Snapshot holds a point-in-time view of all bucket statistics.
type Snapshot struct {
	timestamp time.Time
	buckets   map[string]BucketStats
}

// NewSnapshot creates a new Snapshot from the given bucket stats and timestamp.
func NewSnapshot(buckets map[string]BucketStats, ts time.Time) *Snapshot {
	copy := make(map[string]BucketStats, len(buckets))
	for k, v := range buckets {
		copy[k] = v
	}
	return &Snapshot{
		timestamp: ts,
		buckets:   copy,
	}
}

// Timestamp returns when this snapshot was taken.
func (s *Snapshot) Timestamp() time.Time {
	return s.timestamp
}

// Buckets returns a copy of the bucket stats map.
func (s *Snapshot) Buckets() map[string]BucketStats {
	copy := make(map[string]BucketStats, len(s.buckets))
	for k, v := range s.buckets {
		copy[k] = v
	}
	return copy
}

// IsEmpty returns true if the snapshot contains no buckets.
func (s *Snapshot) IsEmpty() bool {
	return len(s.buckets) == 0
}

// TotalKeys returns the sum of all keys across all buckets.
func (s *Snapshot) TotalKeys() int {
	total := 0
	for _, b := range s.buckets {
		total += b.KeyCount
	}
	return total
}

// TotalSize returns the sum of all size bytes across all buckets.
func (s *Snapshot) TotalSize() int64 {
	var total int64
	for _, b := range s.buckets {
		total += b.SizeBytes
	}
	return total
}

// Get returns the BucketStats for a named bucket and whether it exists.
func (s *Snapshot) Get(name string) (BucketStats, bool) {
	bs, ok := s.buckets[name]
	return bs, ok
}
