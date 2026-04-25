package stats

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func makeExportSnapshot() *Snapshot {
	buckets := map[string]BucketStats{
		"users": {KeyCount: 10, SizeBytes: 1024},
		"sessions": {KeyCount: 5, SizeBytes: 512},
	}
	return NewSnapshot(buckets, time.Now())
}

func TestExportSnapshot_NilSnapshot(t *testing.T) {
	var buf bytes.Buffer
	err := ExportSnapshot(&buf, nil, FormatJSON)
	if err == nil {
		t.Fatal("expected error for nil snapshot")
	}
}

func TestExportSnapshot_UnsupportedFormat(t *testing.T) {
	snap := makeExportSnapshot()
	var buf bytes.Buffer
	err := ExportSnapshot(&buf, snap, ExportFormat("xml"))
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestExportSnapshot_JSON(t *testing.T) {
	snap := makeExportSnapshot()
	var buf bytes.Buffer
	err := ExportSnapshot(&buf, snap, FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var record ExportRecord
	if err := json.Unmarshal(buf.Bytes(), &record); err != nil {
		t.Fatalf("failed to unmarshal JSON output: %v", err)
	}
	if len(record.Buckets) != 2 {
		t.Errorf("expected 2 buckets, got %d", len(record.Buckets))
	}
	if record.TotalKeys != 15 {
		t.Errorf("expected total keys 15, got %d", record.TotalKeys)
	}
}

func TestExportSnapshot_CSV(t *testing.T) {
	snap := makeExportSnapshot()
	var buf bytes.Buffer
	err := ExportSnapshot(&buf, snap, FormatCSV)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	// 1 header + 2 bucket rows
	if len(lines) != 3 {
		t.Errorf("expected 3 lines in CSV, got %d", len(lines))
	}
	if lines[0] != "timestamp,bucket,keys,size_bytes" {
		t.Errorf("unexpected CSV header: %s", lines[0])
	}
}

func TestExportSnapshot_CSV_ContainsBuckets(t *testing.T) {
	snap := makeExportSnapshot()
	var buf bytes.Buffer
	_ = ExportSnapshot(&buf, snap, FormatCSV)
	output := buf.String()
	if !strings.Contains(output, "users") {
		t.Error("expected CSV to contain bucket 'users'")
	}
	if !strings.Contains(output, "sessions") {
		t.Error("expected CSV to contain bucket 'sessions'")
	}
}
