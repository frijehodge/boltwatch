package watch_test

import (
	"testing"
	"time"

	"github.com/user/boltwatch/internal/watch"
)

func TestDefaultOptions(t *testing.T) {
	opts := watch.DefaultOptions()
	if opts.Interval != 2*time.Second {
		t.Errorf("expected 2s interval, got %v", opts.Interval)
	}
	if opts.Verbose {
		t.Error("expected Verbose to be false by default")
	}
}

func TestWithInterval(t *testing.T) {
	opts := watch.DefaultOptions().WithInterval(500 * time.Millisecond)
	if opts.Interval != 500*time.Millisecond {
		t.Errorf("expected 500ms, got %v", opts.Interval)
	}
}

func TestWithVerbose(t *testing.T) {
	opts := watch.DefaultOptions().WithVerbose(true)
	if !opts.Verbose {
		t.Error("expected Verbose to be true")
	}
}

func TestOptionsImmutability(t *testing.T) {
	base := watch.DefaultOptions()
	modified := base.WithInterval(10 * time.Second).WithVerbose(true)

	if base.Interval == modified.Interval {
		t.Error("base interval should not be modified")
	}
	if base.Verbose == modified.Verbose {
		t.Error("base verbose flag should not be modified")
	}
}
