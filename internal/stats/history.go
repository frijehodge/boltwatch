package stats

import (
	"sync"
	"time"
)

// SnapshotEntry holds a point-in-time snapshot of bucket statistics.
type SnapshotEntry struct {
	Timestamp time.Time
	Buckets   []BucketStats
}

// History stores a rolling window of stat snapshots for trend analysis.
type History struct {
	mu       sync.RWMutex
	snapshots []SnapshotEntry
	maxSize  int
}

// NewHistory creates a History that retains at most maxSize snapshots.
func NewHistory(maxSize int) *History {
	if maxSize <= 0 {
		maxSize = 60
	}
	return &History{maxSize: maxSize}
}

// Add appends a new snapshot, evicting the oldest if capacity is exceeded.
func (h *History) Add(buckets []BucketStats) {
	h.mu.Lock()
	defer h.mu.Unlock()

	entry := SnapshotEntry{
		Timestamp: time.Now(),
		Buckets:   buckets,
	}
	h.snapshots = append(h.snapshots, entry)
	if len(h.snapshots) > h.maxSize {
		h.snapshots = h.snapshots[len(h.snapshots)-h.maxSize:]
	}
}

// Latest returns the most recent snapshot, or nil if history is empty.
func (h *History) Latest() *SnapshotEntry {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if len(h.snapshots) == 0 {
		return nil
	}
	copy := h.snapshots[len(h.snapshots)-1]
	return &copy
}

// All returns a copy of all stored snapshots in chronological order.
func (h *History) All() []SnapshotEntry {
	h.mu.RLock()
	defer h.mu.RUnlock()
	out := make([]SnapshotEntry, len(h.snapshots))
	copy(out, h.snapshots)
	return out
}

// Len returns the current number of stored snapshots.
func (h *History) Len() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.snapshots)
}
