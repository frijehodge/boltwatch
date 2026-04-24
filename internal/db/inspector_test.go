package db

import (
	"os"
	"testing"

	bolt "go.etcd.io/bbolt"
)

func createTestDB(t *testing.T) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "boltwatch-test-*.db")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	f.Close()

	db, err := bolt.Open(f.Name(), 0600, nil)
	if err != nil {
		t.Fatalf("failed to open bolt db: %v", err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte("users"))
		if err != nil {
			return err
		}
		_ = b.Put([]byte("alice"), []byte("data1"))
		_ = b.Put([]byte("bob"), []byte("data2"))

		_, err = tx.CreateBucket([]byte("sessions"))
		return err
	})
	if err != nil {
		t.Fatalf("failed to seed db: %v", err)
	}
	return f.Name()
}

func TestOpen(t *testing.T) {
	path := createTestDB(t)
	insp, err := Open(path)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer insp.Close()

	if insp.Path != path {
		t.Errorf("expected path %q, got %q", path, insp.Path)
	}
}

func TestCollectStats(t *testing.T) {
	path := createTestDB(t)
	insp, err := Open(path)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer insp.Close()

	stats, err := insp.CollectStats()
	if err != nil {
		t.Fatalf("CollectStats() error = %v", err)
	}

	if len(stats) != 2 {
		t.Fatalf("expected 2 buckets, got %d", len(stats))
	}

	bucketMap := make(map[string]BucketStats)
	for _, s := range stats {
		bucketMap[s.Name] = s
	}

	users, ok := bucketMap["users"]
	if !ok {
		t.Fatal("expected bucket 'users' not found")
	}
	if users.KeyCount != 2 {
		t.Errorf("expected 2 keys in 'users', got %d", users.KeyCount)
	}

	if _, ok := bucketMap["sessions"]; !ok {
		t.Fatal("expected bucket 'sessions' not found")
	}
}
