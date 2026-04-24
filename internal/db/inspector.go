package db

import (
	"fmt"
	"time"

	bolt "go.etcd.io/bbolt"
)

// BucketStats holds statistics for a single BoltDB bucket.
type BucketStats struct {
	Name       string
	KeyCount   int
	Depth      int
	BranchPages int
	LeafPages   int
	OverflowPages int
	LastUpdated time.Time
}

// Inspector wraps a BoltDB file and provides stat collection.
type Inspector struct {
	Path string
	db   *bolt.DB
}

// Open opens the BoltDB file in read-only mode.
func Open(path string) (*Inspector, error) {
	db, err := bolt.Open(path, 0600, &bolt.Options{
		ReadOnly: true,
		Timeout:  2 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open db %q: %w", path, err)
	}
	return &Inspector{Path: path, db: db}, nil
}

// Close closes the underlying database.
func (i *Inspector) Close() error {
	return i.db.Close()
}

// CollectStats returns statistics for all top-level buckets.
func (i *Inspector) CollectStats() ([]BucketStats, error) {
	var stats []BucketStats

	err := i.db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			bs := b.Stats()
			stats = append(stats, BucketStats{
				Name:          string(name),
				KeyCount:      bs.KeyN,
				Depth:         bs.Depth,
				BranchPages:   bs.BranchPageN,
				LeafPages:     bs.LeafPageN,
				OverflowPages: bs.LeafOverflowN,
				LastUpdated:   time.Now(),
			})
			return nil
		})
	})
	if err != nil {
		return nil, fmt.Errorf("failed to collect stats: %w", err)
	}
	return stats, nil
}
