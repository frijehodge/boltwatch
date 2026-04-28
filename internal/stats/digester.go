package stats

import (
	"crypto/md5"
	"fmt"
	"sort"
	"strings"
)

// DigestResult holds a fingerprint summary of a snapshot.
type DigestResult struct {
	Hash       string
	BucketCount int
	KeyCount   int64
	ByteSize   int64
}

// DigestSnapshot computes a stable hash fingerprint for a snapshot,
// useful for detecting whether anything meaningful has changed between
// two polling cycles without doing a full deep comparison.
func DigestSnapshot(s *Snapshot) *DigestResult {
	if s == nil || s.IsEmpty() {
		return &DigestResult{Hash: emptyHash()}
	}

	// Sort bucket names for a stable hash.
	names := make([]string, 0, len(s.Buckets))
	for name := range s.Buckets {
		names = append(names, name)
	}
	sort.Strings(names)

	var sb strings.Builder
	var totalKeys int64
	var totalSize int64

	for _, name := range names {
		bs := s.Buckets[name]
		fmt.Fprintf(&sb, "%s:%d:%d;", name, bs.KeyN, bs.LeafInuse)
		totalKeys += int64(bs.KeyN)
		totalSize += int64(bs.LeafInuse)
	}

	h := md5.Sum([]byte(sb.String()))

	return &DigestResult{
		Hash:        fmt.Sprintf("%x", h),
		BucketCount: len(names),
		KeyCount:    totalKeys,
		ByteSize:    totalSize,
	}
}

// DigestsMatch returns true if two digest results represent identical snapshots.
func DigestsMatch(a, b *DigestResult) bool {
	if a == nil || b == nil {
		return a == b
	}
	return a.Hash == b.Hash
}

func emptyHash() string {
	h := md5.Sum([]byte(""))
	return fmt.Sprintf("%x", h)
}
