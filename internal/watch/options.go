package watch

import "time"

// Options controls the behaviour of a Watcher.
type Options struct {
	// Interval is how often the watcher polls the database.
	Interval time.Duration
	// Verbose enables extra output such as unchanged snapshots.
	Verbose bool
	// BucketPrefix restricts monitoring to buckets with this prefix.
	BucketPrefix string
	// MinKeys skips buckets with fewer keys than this value.
	MinKeys int
}

// DefaultOptions returns a sensible default Options configuration.
func DefaultOptions() Options {
	return Options{
		Interval: 5 * time.Second,
		Verbose:  false,
	}
}

// Option is a functional option for configuring a Watcher.
type Option func(*Options)

// WithInterval sets the polling interval.
func WithInterval(d time.Duration) Option {
	return func(o *Options) {
		o.Interval = d
	}
}

// WithVerbose enables or disables verbose output.
func WithVerbose(v bool) Option {
	return func(o *Options) {
		o.Verbose = v
	}
}

// WithBucketPrefix restricts watched buckets to those matching the prefix.
func WithBucketPrefix(prefix string) Option {
	return func(o *Options) {
		o.BucketPrefix = prefix
	}
}

// WithMinKeys skips buckets that have fewer than minKeys keys.
func WithMinKeys(minKeys int) Option {
	return func(o *Options) {
		o.MinKeys = minKeys
	}
}

// applyOptions starts from defaults and applies each Option in order.
func applyOptions(opts []Option) Options {
	cfg := DefaultOptions()
	for _, o := range opts {
		o(&cfg)
	}
	return cfg
}
