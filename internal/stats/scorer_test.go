package stats

import (
	"testing"

	"go.etcd.io/bbolt"
)

func makeScorerSnapshot(data map[string]bbolt.BucketStats) *Snapshot {
	return NewSnapshot(data)
}

func TestScoreSnapshot_Nil(t *testing.T) {
	result := ScoreSnapshot(nil, DefaultScoreWeights(), nil)
	if result != nil {
		t.Errorf("expected nil for nil snapshot, got %v", result)
	}
}

func TestScoreSnapshot_Empty(t *testing.T) {
	snap := makeScorerSnapshot(map[string]bbolt.BucketStats{})
	result := ScoreSnapshot(snap, DefaultScoreWeights(), nil)
	if result != nil {
		t.Errorf("expected nil for empty snapshot, got %v", result)
	}
}

func TestScoreSnapshot_SingleBucket_ScoresOne(t *testing.T) {
	snap := makeScorerSnapshot(map[string]bbolt.BucketStats{
		"alpha": {KeyN: 100, LeafInuse: 4096},
	})
	weights := ScoreWeights{KeysWeight: 1.0, SizeWeight: 0.0, GrowthWeight: 0.0}
	result := ScoreSnapshot(snap, weights, nil)

	if len(result) != 1 {
		t.Fatalf("expected 1 score, got %d", len(result))
	}
	if result[0].Bucket != "alpha" {
		t.Errorf("expected bucket 'alpha', got %q", result[0].Bucket)
	}
	if result[0].Score != 1.0 {
		t.Errorf("expected score 1.0, got %f", result[0].Score)
	}
}

func TestScoreSnapshot_RelativeScoring(t *testing.T) {
	snap := makeScorerSnapshot(map[string]bbolt.BucketStats{
		"small": {KeyN: 10, LeafInuse: 512},
		"large": {KeyN: 100, LeafInuse: 4096},
	})
	weights := ScoreWeights{KeysWeight: 1.0, SizeWeight: 0.0, GrowthWeight: 0.0}
	result := ScoreSnapshot(snap, weights, nil)

	scores := make(map[string]float64)
	for _, s := range result {
		scores[s.Bucket] = s.Score
	}

	if scores["large"] <= scores["small"] {
		t.Errorf("expected 'large' to score higher than 'small', got large=%f small=%f",
			scores["large"], scores["small"])
	}
}

func TestScoreSnapshot_GrowthTrendBoost(t *testing.T) {
	snap := makeScorerSnapshot(map[string]bbolt.BucketStats{
		"growing": {KeyN: 50, LeafInuse: 2048},
		"stable":  {KeyN: 50, LeafInuse: 2048},
	})
	weights := ScoreWeights{KeysWeight: 0.0, SizeWeight: 0.0, GrowthWeight: 1.0}
	trends := map[string]string{
		"growing": "up",
		"stable":  "stable",
	}
	result := ScoreSnapshot(snap, weights, trends)

	scores := make(map[string]float64)
	for _, s := range result {
		scores[s.Bucket] = s.Score
	}

	if scores["growing"] <= scores["stable"] {
		t.Errorf("expected 'growing' to outscore 'stable', got growing=%f stable=%f",
			scores["growing"], scores["stable"])
	}
}

func TestDefaultScoreWeights_SumApproxOne(t *testing.T) {
	w := DefaultScoreWeights()
	sum := w.KeysWeight + w.SizeWeight + w.GrowthWeight
	if sum < 0.99 || sum > 1.01 {
		t.Errorf("expected weights to sum to ~1.0, got %f", sum)
	}
}
