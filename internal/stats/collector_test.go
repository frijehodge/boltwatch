package stats_test

import (
	"os"
	"testing"

	"go.etcd.io/bbolt"

	"github.com/yourorg/boltwatch/internal/stats"
)

func createTestDB(t *testing.T) (*bbolt.DB, func()) {
	t.Helper()
	f, err := os.CreateTemp("", "boltwatch-stats-*.db")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	f.Close()

	db, err := bbolt.Open(f.Name(), 0600, nil)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	return db, func() {
		db.Close()
		os.Remove(f.Name())
	}
}

func TestCollect_EmptyDB(t *testing.T) {
	db, cleanup := createTestDB(t)
	defer cleanup()

	c := stats.NewCollector(db)
	snap, err := c.Collect()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(snap.Buckets) != 0 {
		t.Errorf("expected 0 buckets, got %d", len(snap.Buckets))
	}
}

func TestCollect_WithBuckets(t *testing.T) {
	db, cleanup := createTestDB(t)
	defer cleanup()

	err := db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucket([]byte("alpha"))
		if err != nil {
			return err
		}
		for i := 0; i < 5; i++ {
			key := []byte(string(rune('a' + i)))
			if err := b.Put(key, []byte("value")); err != nil {
				return err
			}
		}
		_, err = tx.CreateBucket([]byte("beta"))
		return err
	})
	if err != nil {
		t.Fatalf("setup db: %v", err)
	}

	c := stats.NewCollector(db)
	snap, err := c.Collect()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(snap.Buckets) != 2 {
		t.Fatalf("expected 2 buckets, got %d", len(snap.Buckets))
	}

	for _, bs := range snap.Buckets {
		if bs.Name == "alpha" && bs.KeyCount != 5 {
			t.Errorf("alpha: expected 5 keys, got %d", bs.KeyCount)
		}
	}
}
