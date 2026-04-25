package stats

import "strings"

// FilterConfig holds configuration for filtering bucket stats.
type FilterConfig struct {
	// Prefix filters buckets whose names start with the given prefix.
	Prefix string
	// MinKeys filters out buckets with fewer keys than this threshold.
	MinKeys int
	// MinSize filters out buckets with total size below this threshold (bytes).
	MinSize int64
}

// FilterSnapshot returns a new Snapshot containing only buckets that match
// the given FilterConfig. The original snapshot is not modified.
func FilterSnapshot(snap *Snapshot, cfg FilterConfig) *Snapshot {
	if snap == nil {
		return nil
	}

	filtered := make(map[string]BucketStats)
	for name, bs := range snap.Buckets {
		if cfg.Prefix != "" && !strings.HasPrefix(name, cfg.Prefix) {
			continue
		}
		if cfg.MinKeys > 0 && bs.KeyCount < cfg.MinKeys {
			continue
		}
		if cfg.MinSize > 0 && bs.Size < cfg.MinSize {
			continue
		}
		filtered[name] = bs
	}

	return &Snapshot{
		Timestamp: snap.Timestamp,
		Buckets:   filtered,
	}
}

// MatchesFilter reports whether a single bucket name and stats satisfy cfg.
func MatchesFilter(name string, bs BucketStats, cfg FilterConfig) bool {
	if cfg.Prefix != "" && !strings.HasPrefix(name, cfg.Prefix) {
		return false
	}
	if cfg.MinKeys > 0 && bs.KeyCount < cfg.MinKeys {
		return false
	}
	if cfg.MinSize > 0 && bs.Size < cfg.MinSize {
		return false
	}
	return true
}
