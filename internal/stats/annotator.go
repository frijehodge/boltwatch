package stats

import "fmt"

// Annotation holds a human-readable note attached to a bucket.
type Annotation struct {
	Bucket  string
	Level   string // "info", "warn", "critical"
	Message string
}

// AnnotateSnapshot inspects a snapshot and returns a list of annotations
// describing notable conditions for each bucket (e.g. large size, high key
// count, rapid growth).
func AnnotateSnapshot(s *Snapshot, cfg AlertConfig) []Annotation {
	if s == nil || s.IsEmpty() {
		return nil
	}

	var annotations []Annotation

	for name, bs := range s.Buckets {
		if cfg.MaxKeys > 0 && bs.KeyCount >= cfg.MaxKeys {
			annotations = append(annotations, Annotation{
				Bucket:  name,
				Level:   "critical",
				Message: fmt.Sprintf("key count %d exceeds limit %d", bs.KeyCount, cfg.MaxKeys),
			})
		} else if cfg.WarnKeys > 0 && bs.KeyCount >= cfg.WarnKeys {
			annotations = append(annotations, Annotation{
				Bucket:  name,
				Level:   "warn",
				Message: fmt.Sprintf("key count %d approaching limit %d", bs.KeyCount, cfg.MaxKeys),
			})
		}

		if cfg.MaxSize > 0 && bs.Size >= cfg.MaxSize {
			annotations = append(annotations, Annotation{
				Bucket:  name,
				Level:   "critical",
				Message: fmt.Sprintf("size %s exceeds limit %s", formatBytes(bs.Size), formatBytes(cfg.MaxSize)),
			})
		} else if cfg.WarnSize > 0 && bs.Size >= cfg.WarnSize {
			annotations = append(annotations, Annotation{
				Bucket:  name,
				Level:   "warn",
				Message: fmt.Sprintf("size %s approaching limit %s", formatBytes(bs.Size), formatBytes(cfg.MaxSize)),
			})
		}
	}

	return annotations
}

// FormatAnnotations returns a human-readable string for a slice of annotations.
func FormatAnnotations(annotations []Annotation) string {
	if len(annotations) == 0 {
		return ""
	}

	var out string
	for _, a := range annotations {
		out += fmt.Sprintf("[%s] %s: %s\n", a.Level, a.Bucket, a.Message)
	}
	return out
}
