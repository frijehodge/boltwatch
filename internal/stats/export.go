package stats

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// ExportFormat defines the output format for exported stats.
type ExportFormat string

const (
	FormatJSON ExportFormat = "json"
	FormatCSV  ExportFormat = "csv"
)

// ExportRecord represents a single exportable snapshot record.
type ExportRecord struct {
	Timestamp  time.Time            `json:"timestamp"`
	Buckets    map[string]BucketStats `json:"buckets"`
	TotalKeys  int                  `json:"total_keys"`
	TotalBytes int64                `json:"total_bytes"`
}

// ExportSnapshot writes a snapshot to the given writer in the specified format.
func ExportSnapshot(w io.Writer, snap *Snapshot, format ExportFormat) error {
	if snap == nil {
		return fmt.Errorf("snapshot is nil")
	}
	record := ExportRecord{
		Timestamp:  snap.Timestamp(),
		Buckets:    snap.Buckets(),
		TotalKeys:  snap.TotalKeys(),
		TotalBytes: snap.TotalSize(),
	}
	switch format {
	case FormatJSON:
		return exportJSON(w, record)
	case FormatCSV:
		return exportCSV(w, record)
	default:
		return fmt.Errorf("unsupported export format: %s", format)
	}
}

func exportJSON(w io.Writer, record ExportRecord) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(record)
}

func exportCSV(w io.Writer, record ExportRecord) error {
	_, err := fmt.Fprintf(w, "timestamp,bucket,keys,size_bytes\n")
	if err != nil {
		return err
	}
	ts := record.Timestamp.Format(time.RFC3339)
	for name, bs := range record.Buckets {
		_, err := fmt.Fprintf(w, "%s,%s,%d,%d\n", ts, name, bs.KeyCount, bs.SizeBytes)
		if err != nil {
			return err
		}
	}
	return nil
}
