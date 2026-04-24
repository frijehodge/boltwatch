package stats

import "fmt"

// BucketStats holds statistics for a single BoltDB bucket.
type BucketStats struct {
	Name     string
	KeyCount int
	Size     int64
	Depth    int
}

// String returns a human-readable one-line summary of the bucket stats.
func (b BucketStats) String() string {
	return fmt.Sprintf("bucket=%q keys=%d size=%s depth=%d",
		b.Name, b.KeyCount, formatBytes(b.Size), b.Depth)
}

// formatBytes converts a byte count into a human-readable string.
func formatBytes(n int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)
	switch {
	case n >= GB:
		return fmt.Sprintf("%.2f GB", float64(n)/float64(GB))
	case n >= MB:
		return fmt.Sprintf("%.2f MB", float64(n)/float64(MB))
	case n >= KB:
		return fmt.Sprintf("%.2f KB", float64(n)/float64(KB))
	default:
		return fmt.Sprintf("%d B", n)
	}
}

// TotalKeys returns the sum of all key counts across the provided stats.
func TotalKeys(stats []BucketStats) int {
	total := 0
	for _, s := range stats {
		total += s.KeyCount
	}
	return total
}

// TotalSize returns the sum of all sizes across the provided stats.
func TotalSize(stats []BucketStats) int64 {
	var total int64
	for _, s := range stats {
		total += s.Size
	}
	return total
}
