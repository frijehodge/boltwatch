package stats

import "fmt"

// DiffResult holds the result of diffing two snapshots.
type DiffResult struct {
	Added   []string
	Removed []string
	Changed []string
	Unchanged []string
}

// IsEmpty returns true if there are no differences.
func (d *DiffResult) IsEmpty() bool {
	return len(d.Added) == 0 && len(d.Removed) == 0 && len(d.Changed) == 0
}

// DiffSnapshots compares two snapshots and returns a DiffResult describing
// which buckets were added, removed, changed, or unchanged.
func DiffSnapshots(before, after *Snapshot) *DiffResult {
	result := &DiffResult{}

	if before == nil && after == nil {
		return result
	}
	if before == nil {
		for name := range after.Buckets {
			result.Added = append(result.Added, name)
		}
		return result
	}
	if after == nil {
		for name := range before.Buckets {
			result.Removed = append(result.Removed, name)
		}
		return result
	}

	for name, afterStats := range after.Buckets {
		beforeStats, exists := before.Buckets[name]
		if !exists {
			result.Added = append(result.Added, name)
			continue
		}
		if beforeStats.KeyCount != afterStats.KeyCount || beforeStats.Size != afterStats.Size {
			result.Changed = append(result.Changed, name)
		} else {
			result.Unchanged = append(result.Unchanged, name)
		}
	}

	for name := range before.Buckets {
		if _, exists := after.Buckets[name]; !exists {
			result.Removed = append(result.Removed, name)
		}
	}

	return result
}

// FormatDiff returns a human-readable summary of a DiffResult.
func FormatDiff(d *DiffResult) string {
	if d == nil || d.IsEmpty() {
		return "No differences found."
	}
	out := ""
	for _, name := range d.Added {
		out += fmt.Sprintf("  + %s (added)\n", name)
	}
	for _, name := range d.Removed {
		out += fmt.Sprintf("  - %s (removed)\n", name)
	}
	for _, name := range d.Changed {
		out += fmt.Sprintf("  ~ %s (changed)\n", name)
	}
	return out
}
