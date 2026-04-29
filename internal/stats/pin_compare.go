package stats

// PinComparison holds the result of comparing a live snapshot to a pin.
type PinComparison struct {
	Label       string
	NewBuckets  []string
	Dropped     []string
	Grown       []string
	Shrunk      []string
	Unchanged   []string
}

// CompareToPinned compares a live snapshot against a named pin on the board.
// Returns nil if the pin does not exist or either snapshot is nil.
func CompareToPinned(pb *PinBoard, label string, live *Snapshot) *PinComparison {
	if pb == nil || live == nil {
		return nil
	}
	pin := pb.Get(label)
	if pin == nil || pin.Snapshot == nil {
		return nil
	}

	cmp := &PinComparison{Label: label}

	for name, liveStat := range live.Buckets {
		pinStat, existed := pin.Snapshot.Buckets[name]
		if !existed {
			cmp.NewBuckets = append(cmp.NewBuckets, name)
			continue
		}
		switch {
		case liveStat.KeyCount > pinStat.KeyCount:
			cmp.Grown = append(cmp.Grown, name)
		case liveStat.KeyCount < pinStat.KeyCount:
			cmp.Shrunk = append(cmp.Shrunk, name)
		default:
			cmp.Unchanged = append(cmp.Unchanged, name)
		}
	}

	for name := range pin.Snapshot.Buckets {
		if _, ok := live.Buckets[name]; !ok {
			cmp.Dropped = append(cmp.Dropped, name)
		}
	}

	return cmp
}

// HasPinChanges returns true if the comparison found any growth, shrinkage,
// new buckets, or dropped buckets relative to the pin.
func HasPinChanges(cmp *PinComparison) bool {
	if cmp == nil {
		return false
	}
	return len(cmp.Grown) > 0 || len(cmp.Shrunk) > 0 ||
		len(cmp.NewBuckets) > 0 || len(cmp.Dropped) > 0
}
