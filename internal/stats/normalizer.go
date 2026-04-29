package stats

// NormalizeConfig holds configuration for snapshot normalization.
type NormalizeConfig struct {
	// CapSize sets the maximum allowed size per bucket (0 = no cap).
	CapSize int64
	// CapKeys sets the maximum allowed key count per bucket (0 = no cap).
	CapKeys int
	// ZeroNegative replaces any negative values with zero.
	ZeroNegative bool
}

// DefaultNormalizeConfig returns a NormalizeConfig with safe defaults.
func DefaultNormalizeConfig() NormalizeConfig {
	return NormalizeConfig{
		CapSize:      0,
		CapKeys:      0,
		ZeroNegative: true,
	}
}

// NormalizeSnapshot returns a new Snapshot with values adjusted according to cfg.
// Nil input returns nil. Buckets with invalid values are corrected in place.
func NormalizeSnapshot(snap *Snapshot, cfg NormalizeConfig) *Snapshot {
	if snap == nil {
		return nil
	}

	normalized := &Snapshot{
		Timestamp: snap.Timestamp,
		Buckets:   make(map[string]BucketStats, len(snap.Buckets)),
	}

	for name, bs := range snap.Buckets {
		keys := bs.KeyCount
		size := bs.Size

		if cfg.ZeroNegative {
			if keys < 0 {
				keys = 0
			}
			if size < 0 {
				size = 0
			}
		}

		if cfg.CapKeys > 0 && keys > cfg.CapKeys {
			keys = cfg.CapKeys
		}

		if cfg.CapSize > 0 && size > cfg.CapSize {
			size = cfg.CapSize
		}

		normalized.Buckets[name] = BucketStats{
			KeyCount: keys,
			Size:     size,
		}
	}

	return normalized
}
