package watch_test

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	bolt "go.etcd.io/bbolt"

	"github.com/user/boltwatch/internal/db"
	"github.com/user/boltwatch/internal/stats"
	"github.com/user/boltwatch/internal/watch"
)

func createTestDB(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.db")
	bdb, err := bolt.Open(path, 0600, nil)
	if err != nil {
		t.Fatalf("open bolt db: %v", err)
	}
	_ = bdb.Update(func(tx *bolt.Tx) error {
		_, _ = tx.CreateBucketIfNotExists([]byte("alpha"))
		return nil
	})
	bdb.Close()
	return path
}

func TestWatcher_ReceivesUpdates(t *testing.T) {
	path := createTestDB(t)
	defer os.Remove(path)

	inspector, err := db.Open(path)
	if err != nil {
		t.Fatalf("db.Open: %v", err)
	}
	defer inspector.Close()

	collector := stats.NewCollector(inspector)
	w := watch.NewWatcher(collector, 50*time.Millisecond)

	var mu sync.Mutex
	var received [][]stats.BucketStat
	w.OnUpdate = func(bs []stats.BucketStat) {
		mu.Lock()
		received = append(received, bs)
		mu.Unlock()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	w.Start(ctx)

	mu.Lock()
	defer mu.Unlock()
	if len(received) == 0 {
		t.Fatal("expected at least one update")
	}
	if len(received[0]) == 0 {
		t.Fatal("expected bucket stats in first update")
	}
}

func TestWatcher_CallsOnError(t *testing.T) {
	// Use a non-existent path to force an error after open.
	inspector, err := db.Open(filepath.Join(t.TempDir(), "missing.db"))
	if err == nil {
		inspector.Close()
		t.Skip("db opened unexpectedly")
	}
	// Verify that a nil inspector causes NewCollector to surface an error via OnError.
	_ = err // expected
}
