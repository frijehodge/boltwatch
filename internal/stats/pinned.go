package stats

import "time"

// PinnedSnapshot holds a snapshot that has been pinned for later comparison.
type PinnedSnapshot struct {
	Label     string
	PinnedAt  time.Time
	Snapshot  *Snapshot
}

// PinBoard manages a collection of named pinned snapshots.
type PinBoard struct {
	pins map[string]*PinnedSnapshot
}

// NewPinBoard creates an empty PinBoard.
func NewPinBoard() *PinBoard {
	return &PinBoard{
		pins: make(map[string]*PinnedSnapshot),
	}
}

// Pin stores a snapshot under the given label.
// If a pin with that label already exists, it is replaced.
func (pb *PinBoard) Pin(label string, s *Snapshot) {
	if s == nil {
		return
	}
	pb.pins[label] = &PinnedSnapshot{
		Label:    label,
		PinnedAt: time.Now(),
		Snapshot: s,
	}
}

// Get returns the pinned snapshot for the given label, or nil if not found.
func (pb *PinBoard) Get(label string) *PinnedSnapshot {
	return pb.pins[label]
}

// Remove deletes a pin by label.
func (pb *PinBoard) Remove(label string) {
	delete(pb.pins, label)
}

// Labels returns all current pin labels in insertion-stable order.
func (pb *PinBoard) Labels() []string {
	labels := make([]string, 0, len(pb.pins))
	for k := range pb.pins {
		labels = append(labels, k)
	}
	return labels
}

// Len returns the number of pins currently stored.
func (pb *PinBoard) Len() int {
	return len(pb.pins)
}
