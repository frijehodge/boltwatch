package stats

import (
	"strings"
	"testing"
)

func makeAnnotatorSnapshot(buckets map[string]BucketStats) *Snapshot {
	s := &Snapshot{Buckets: make(map[string]BucketStats)}
	for k, v := range buckets {
		s.Buckets[k] = v
	}
	return s
}

func TestAnnotateSnapshot_Nil(t *testing.T) {
	cfg := DefaultAlertConfig()
	result := AnnotateSnapshot(nil, cfg)
	if result != nil {
		t.Errorf("expected nil for nil snapshot, got %v", result)
	}
}

func TestAnnotateSnapshot_NoViolations(t *testing.T) {
	s := makeAnnotatorSnapshot(map[string]BucketStats{
		"users": {KeyCount: 10, Size: 1024},
	})
	cfg := AlertConfig{MaxKeys: 1000, WarnKeys: 800, MaxSize: 1 << 20, WarnSize: 800 * 1024}
	result := AnnotateSnapshot(s, cfg)
	if len(result) != 0 {
		t.Errorf("expected no annotations, got %d", len(result))
	}
}

func TestAnnotateSnapshot_CriticalKeys(t *testing.T) {
	s := makeAnnotatorSnapshot(map[string]BucketStats{
		"logs": {KeyCount: 5000, Size: 512},
	})
	cfg := AlertConfig{MaxKeys: 1000, WarnKeys: 800}
	result := AnnotateSnapshot(s, cfg)
	if len(result) != 1 {
		t.Fatalf("expected 1 annotation, got %d", len(result))
	}
	if result[0].Level != "critical" {
		t.Errorf("expected critical level, got %s", result[0].Level)
	}
	if result[0].Bucket != "logs" {
		t.Errorf("expected bucket 'logs', got %s", result[0].Bucket)
	}
}

func TestAnnotateSnapshot_WarnSize(t *testing.T) {
	s := makeAnnotatorSnapshot(map[string]BucketStats{
		"events": {KeyCount: 5, Size: 900 * 1024},
	})
	cfg := AlertConfig{MaxSize: 1 << 20, WarnSize: 800 * 1024}
	result := AnnotateSnapshot(s, cfg)
	if len(result) != 1 {
		t.Fatalf("expected 1 annotation, got %d", len(result))
	}
	if result[0].Level != "warn" {
		t.Errorf("expected warn level, got %s", result[0].Level)
	}
}

func TestFormatAnnotations_Empty(t *testing.T) {
	out := FormatAnnotations(nil)
	if out != "" {
		t.Errorf("expected empty string, got %q", out)
	}
}

func TestFormatAnnotations_WithData(t *testing.T) {
	anns := []Annotation{
		{Bucket: "users", Level: "warn", Message: "key count 800 approaching limit 1000"},
		{Bucket: "logs", Level: "critical", Message: "size 2.0 MiB exceeds limit 1.0 MiB"},
	}
	out := FormatAnnotations(anns)
	if !strings.Contains(out, "[warn]") {
		t.Errorf("expected '[warn]' in output, got: %s", out)
	}
	if !strings.Contains(out, "[critical]") {
		t.Errorf("expected '[critical]' in output, got: %s", out)
	}
	if !strings.Contains(out, "users") {
		t.Errorf("expected 'users' in output, got: %s", out)
	}
}
