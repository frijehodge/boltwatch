package stats

import "time"

// Snapshot represents a point-in-time capture of bucket statistics.
type Snapshot struct {
	Timestamp time.Time
	Buckets   []BucketStats
}

// NewSnapshot creates a new Snapshot with the current time and provided buckets.
func NewSnapshot(buckets []BucketStats) Snapshot {
	return Snapshot{
		Timestamp: time.Now(),
		Buckets:   buckets,
	}
}

// IsEmpty returns true if the snapshot contains no bucket data.
func (s Snapshot) IsEmpty() bool {
	return len(s.Buckets) == 0
}

// TotalKeys returns the total number of keys across all buckets in the snapshot.
func (s Snapshot) TotalKeys() int {
	total := 0
	for _, b := range s.Buckets {
		total += TotalKeys(b)
	}
	return total
}

// TotalSize returns the total size in bytes across all buckets in the snapshot.
func (s Snapshot) TotalSize() int64 {
	var total int64
	for _, b := range s.Buckets {
		total += TotalSize(b)
	}
	return total
}

// BucketNames returns a slice of all bucket names in the snapshot.
func (s Snapshot) BucketNames() []string {
	names := make([]string, 0, len(s.Buckets))
	for _, b := range s.Buckets {
		names = append(names, b.Name)
	}
	return names
}
