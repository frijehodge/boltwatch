package stats

import (
	"testing"
	"time"
)

func makeNormalizerSnapshot() *Snapshot {
	return &Snapshot{
		Timestamp: time.Now(),
		Buckets: map[string]BucketStats{
			"alpha": {KeyCount: 10, Size: 1024},
			"beta":  {KeyCount: -5, Size: -200},
			"gamma": {KeyCount: 9000, Size: 5_000_000},
		},
	}
}

func TestNormalizeSnapshot_Nil(t *testing.T) {
	result := NormalizeSnapshot(nil, DefaultNormalizeConfig())
	if result != nil {
		t.Error("expected nil for nil input")
	}
}

func TestNormalizeSnapshot_ZeroNegative(t *testing.T) {
	snap := makeNormalizerSnapshot()
	cfg := DefaultNormalizeConfig()

	result := NormalizeSnapshot(snap, cfg)

	if result.Buckets["beta"].KeyCount != 0 {
		t.Errorf("expected KeyCount=0 for negative, got %d", result.Buckets["beta"].KeyCount)
	}
	if result.Buckets["beta"].Size != 0 {
		t.Errorf("expected Size=0 for negative, got %d", result.Buckets["beta"].Size)
	}
}

func TestNormalizeSnapshot_CapKeys(t *testing.T) {
	snap := makeNormalizerSnapshot()
	cfg := DefaultNormalizeConfig()
	cfg.CapKeys = 100

	result := NormalizeSnapshot(snap, cfg)

	if result.Buckets["gamma"].KeyCount != 100 {
		t.Errorf("expected KeyCount=100 after cap, got %d", result.Buckets["gamma"].KeyCount)
	}
	if result.Buckets["alpha"].KeyCount != 10 {
		t.Errorf("expected KeyCount=10 unchanged, got %d", result.Buckets["alpha"].KeyCount)
	}
}

func TestNormalizeSnapshot_CapSize(t *testing.T) {
	snap := makeNormalizerSnapshot()
	cfg := DefaultNormalizeConfig()
	cfg.CapSize = 2048

	result := NormalizeSnapshot(snap, cfg)

	if result.Buckets["gamma"].Size != 2048 {
		t.Errorf("expected Size=2048 after cap, got %d", result.Buckets["gamma"].Size)
	}
	if result.Buckets["alpha"].Size != 1024 {
		t.Errorf("expected Size=1024 unchanged, got %d", result.Buckets["alpha"].Size)
	}
}

func TestNormalizeSnapshot_PreservesTimestamp(t *testing.T) {
	snap := makeNormalizerSnapshot()
	result := NormalizeSnapshot(snap, DefaultNormalizeConfig())

	if !result.Timestamp.Equal(snap.Timestamp) {
		t.Error("expected timestamp to be preserved")
	}
}

func TestNormalizeSnapshot_DoesNotMutateOriginal(t *testing.T) {
	snap := makeNormalizerSnapshot()
	origBeta := snap.Buckets["beta"]

	NormalizeSnapshot(snap, DefaultNormalizeConfig())

	if snap.Buckets["beta"].KeyCount != origBeta.KeyCount {
		t.Error("original snapshot was mutated")
	}
}
