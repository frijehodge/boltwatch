package stats

import (
	"fmt"
	"time"
)

// MarkerKind classifies the type of a snapshot marker.
type MarkerKind string

const (
	MarkerKindManual    MarkerKind = "manual"
	MarkerKindAutomatic MarkerKind = "automatic"
	MarkerKindAlert     MarkerKind = "alert"
)

// Marker represents a named point-in-time annotation attached to a snapshot.
type Marker struct {
	Label     string
	Kind      MarkerKind
	Timestamp time.Time
	Snapshot  *Snapshot
	Note      string
}

// MarkerBoard holds a collection of named markers.
type MarkerBoard struct {
	markers map[string]*Marker
}

// NewMarkerBoard creates an empty MarkerBoard.
func NewMarkerBoard() *MarkerBoard {
	return &MarkerBoard{markers: make(map[string]*Marker)}
}

// Add stores a marker under the given label.
func (mb *MarkerBoard) Add(label string, snap *Snapshot, kind MarkerKind, note string) {
	if snap == nil {
		return
	}
	mb.markers[label] = &Marker{
		Label:     label,
		Kind:      kind,
		Timestamp: snap.Timestamp,
		Snapshot:  snap,
		Note:      note,
	}
}

// Get retrieves a marker by label. Returns nil if not found.
func (mb *MarkerBoard) Get(label string) *Marker {
	return mb.markers[label]
}

// Remove deletes a marker by label.
func (mb *MarkerBoard) Remove(label string) {
	delete(mb.markers, label)
}

// All returns all markers as a slice.
func (mb *MarkerBoard) All() []*Marker {
	out := make([]*Marker, 0, len(mb.markers))
	for _, m := range mb.markers {
		out = append(out, m)
	}
	return out
}

// Len returns the number of markers.
func (mb *MarkerBoard) Len() int {
	return len(mb.markers)
}

// String returns a short description of the marker.
func (m *Marker) String() string {
	return fmt.Sprintf("[%s] %s @ %s", m.Kind, m.Label, m.Timestamp.Format(time.RFC3339))
}
