package watch

import (
	"context"
	"time"

	"github.com/user/boltwatch/internal/stats"
)

// Watcher polls a BoltDB file at a given interval and emits stat snapshots.
type Watcher struct {
	collector *stats.Collector
	interval  time.Duration
	OnUpdate  func([]stats.BucketStat)
	OnError   func(error)
}

// NewWatcher creates a Watcher with the given collector and poll interval.
func NewWatcher(collector *stats.Collector, interval time.Duration) *Watcher {
	return &Watcher{
		collector: collector,
		interval:  interval,
	}
}

// Start begins polling until ctx is cancelled.
func (w *Watcher) Start(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	// Emit an initial snapshot immediately.
	w.poll()

	for {
		select {
		case <-ticker.C:
			w.poll()
		case <-ctx.Done():
			return
		}
	}
}

func (w *Watcher) poll() {
	buckets, err := w.collector.Collect()
	if err != nil {
		if w.OnError != nil {
			w.OnError(err)
		}
		return
	}
	if w.OnUpdate != nil {
		w.OnUpdate(buckets)
	}
}
