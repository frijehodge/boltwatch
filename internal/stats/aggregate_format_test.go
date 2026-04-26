package stats

import (
	"strings"
	"testing"
)

func makeAggregateFormatSnapshot() *Snapshot {
	return &Snapshot{
		Buckets: map[string]BucketStats{
			"users":  {KeyCount: 500, Size: 512000},
			"events": {KeyCount: 1200, Size: 2048000},
			"config": {KeyCount: 10, Size: 1024},
		},
	}
}

func TestFormatAggregate_NilSnapshot(t *testing.T) {
	result := FormatAggregate(nil)
	if result == "" {
		t.Error("expected non-empty string for nil snapshot")
	}
	if !strings.Contains(result, "no data") {
		t.Errorf("expected 'no data' message, got: %s", result)
	}
}

func TestFormatAggregate_ContainsHeaders(t *testing.T) {
	snap := makeAggregateFormatSnapshot()
	agg := Aggregate(snap)
	result := FormatAggregate(agg)

	for _, header := range []string{"Bucket", "Keys", "Size"} {
		if !strings.Contains(result, header) {
			t.Errorf("expected header %q in output: %s", header, result)
		}
	}
}

func TestFormatAggregate_TopBucketFirst(t *testing.T) {
	snap := makeAggregateFormatSnapshot()
	agg := Aggregate(snap)
	result := FormatAggregate(agg)

	eventsIdx := strings.Index(result, "events")
	usersIdx := strings.Index(result, "users")
	if eventsIdx == -1 || usersIdx == -1 {
		t.Fatal("expected both 'events' and 'users' in output")
	}
	if eventsIdx > usersIdx {
		t.Error("expected 'events' (highest key count) to appear before 'users'")
	}
}

func TestFormatAggregate_ContainsTotals(t *testing.T) {
	snap := makeAggregateFormatSnapshot()
	agg := Aggregate(snap)
	result := FormatAggregate(agg)

	if !strings.Contains(result, "Total") {
		t.Errorf("expected 'Total' row in aggregate output: %s", result)
	}
}

func TestFormatAggregate_EmptySnapshot(t *testing.T) {
	empty := &Snapshot{Buckets: map[string]BucketStats{}}
	agg := Aggregate(empty)
	result := FormatAggregate(agg)
	if result == "" {
		t.Error("expected non-empty output for empty snapshot aggregate")
	}
}
