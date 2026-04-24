package stats

import (
	"fmt"
	"io"
	"time"
)

// Reporter writes formatted snapshot reports to an output writer.
type Reporter struct {
	out     io.Writer
	verbose bool
}

// NewReporter creates a Reporter that writes to the given writer.
func NewReporter(out io.Writer, verbose bool) *Reporter {
	return &Reporter{out: out, verbose: verbose}
}

// Report writes a formatted report for the given snapshot to the reporter's output.
func (r *Reporter) Report(s Snapshot) {
	fmt.Fprintf(r.out, "\n[%s]\n", s.Timestamp.Format(time.RFC3339))
	if s.IsEmpty() {
		fmt.Fprintln(r.out, "No buckets found.")
		return
	}
	fmt.Fprintln(r.out, FormatTable(s.Buckets))
	if r.verbose {
		fmt.Fprintln(r.out, FormatSummary(s.Buckets))
	}
}

// ReportDelta writes a delta report comparing two snapshots.
func (r *Reporter) ReportDelta(prev, curr Snapshot) {
	deltas := ComputeDeltas(prev.Buckets, curr.Buckets)
	if !HasChanges(deltas) {
		if r.verbose {
			fmt.Fprintln(r.out, "(no changes since last snapshot)")
		}
		return
	}
	fmt.Fprintf(r.out, "\nChanges detected at %s:\n", curr.Timestamp.Format(time.RFC3339))
	for _, d := range deltas {
		if d.KeyDelta != 0 || d.SizeDelta != 0 {
			fmt.Fprintf(r.out, "  %-20s keys: %+d  size: %+s\n",
				d.Name,
				d.KeyDelta,
				formatDeltaBytes(d.SizeDelta),
			)
		}
	}
}

// formatDeltaBytes formats a signed byte delta with a + or - prefix.
func formatDeltaBytes(delta int64) string {
	if delta >= 0 {
		return "+" + formatBytes(uint64(delta))
	}
	return "-" + formatBytes(uint64(-delta))
}
