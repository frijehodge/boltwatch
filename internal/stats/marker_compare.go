package stats

// CompareToMarker computes deltas between a marker's snapshot and the current snapshot.
// Returns nil if either the marker or current snapshot is nil.
func CompareToMarker(current *Snapshot, m *Marker) []BucketDelta {
	if m == nil || current == nil {
		return nil
	}
	return ComputeDeltas(m.Snapshot, current)
}

// HasMarkerChanges returns true if there are any changes between the marker and current snapshot.
func HasMarkerChanges(current *Snapshot, m *Marker) bool {
	deltas := CompareToMarker(current, m)
	return HasChanges(deltas)
}

// MarkerDriftSummary returns a brief string summarising drift from a marker.
func MarkerDriftSummary(current *Snapshot, m *Marker) string {
	if m == nil || current == nil {
		return "unavailable"
	}
	deltas := CompareToMarker(current, m)
	if !HasChanges(deltas) {
		return "no drift"
	}
	changed := 0
	for _, d := range deltas {
		if d.KeysDelta != 0 || d.SizeDelta != 0 {
			changed++
		}
	}
	if changed == 1 {
		return "1 bucket changed"
	}
	return formatInt(changed) + " buckets changed"
}

func formatInt(n int) string {
	return itoa(n)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	result := ""
	for n > 0 {
		result = string(rune('0'+n%10)) + result
		n /= 10
	}
	return result
}
