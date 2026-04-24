package watch

import "time"

// Options holds configuration for the Watcher.
type Options struct {
	// Interval is how often the watcher polls the database.
	Interval time.Duration

	// Verbose enables printing raw stat output to stderr.
	Verbose bool
}

// DefaultOptions returns sensible defaults for the watcher.
func DefaultOptions() Options {
	return Options{
		Interval: 2 * time.Second,
		Verbose:  false,
	}
}

// WithInterval returns a copy of o with the interval overridden.
func (o Options) WithInterval(d time.Duration) Options {
	o.Interval = d
	return o
}

// WithVerbose returns a copy of o with verbose toggled.
func (o Options) WithVerbose(v bool) Options {
	o.Verbose = v
	return o
}
