package stats

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// FormatTable writes a human-readable table of bucket statistics to w.
func FormatTable(w io.Writer, snap *Snapshot) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	fmt.Fprintf(tw, "Collected at: %s\n\n", snap.CollectedAt.Format("2006-01-02 15:04:05"))
	fmt.Fprintln(tw, "BUCKET\tKEYS\tDEPTH\tBRANCH PAGES\tLEAF PAGES")
	fmt.Fprintln(tw, "------\t----\t-----\t------------\t----------")

	for _, b := range snap.Buckets {
		fmt.Fprintf(tw, "%s\t%d\t%d\t%d\t%d\n",
			b.Name,
			b.KeyCount,
			b.Depth,
			b.Branches,
			b.Leaves,
		)
	}

	return tw.Flush()
}

// FormatSummary writes a single-line summary of the snapshot to w.
func FormatSummary(w io.Writer, snap *Snapshot) {
	totalKeys := 0
	for _, b := range snap.Buckets {
		totalKeys += b.KeyCount
	}
	fmt.Fprintf(w, "[%s] buckets=%d total_keys=%d\n",
		snap.CollectedAt.Format("15:04:05"),
		len(snap.Buckets),
		totalKeys,
	)
}
