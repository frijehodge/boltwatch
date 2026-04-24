package stats

import (
	"fmt"
	"time"

	"go.etcd.io/bbolt"
)

// BucketStat holds statistics for a single BoltDB bucket.
type BucketStat struct {
	Name      string
	KeyCount  int
	Depth     int
	Branches  int
	Leaves    int
	Timestamp time.Time
}

// Snapshot holds a point-in-time collection of bucket stats.
type Snapshot struct {
	CollectedAt time.Time
	Buckets     []BucketStat
}

// Collector gathers bucket statistics from an open BoltDB instance.
type Collector struct {
	db *bbolt.DB
}

// NewCollector creates a new Collector for the given database.
func NewCollector(db *bbolt.DB) *Collector {
	return &Collector{db: db}
}

// Collect performs a read-only transaction and gathers stats for all top-level buckets.
func (c *Collector) Collect() (*Snapshot, error) {
	snap := &Snapshot{
		CollectedAt: time.Now(),
	}

	err := c.db.View(func(tx *bbolt.Tx) error {
		return tx.ForEach(func(name []byte, b *bbolt.Bucket) error {
			bs := b.Stats()
			stat := BucketStat{
				Name:      string(name),
				KeyCount:  bs.KeyN,
				Depth:     bs.Depth,
				Branches:  bs.BranchPageN,
				Leaves:    bs.LeafPageN,
				Timestamp: snap.CollectedAt,
			}
			snap.Buckets = append(snap.Buckets, stat)
			return nil
		})
	})
	if err != nil {
		return nil, fmt.Errorf("collect stats: %w", err)
	}

	return snap, nil
}
