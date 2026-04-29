package stats

import (
	"strings"
	"testing"
)

func TestLevelLabel(t *testing.T) {
	tests := []struct {
		level WatchdogLevel
		want  string
	}{
		{WatchdogInfo, "INFO"},
		{WatchdogWarn, "WARN"},
		{WatchdogCritical, "CRIT"},
	}
	for _, tt := range tests {
		if got := levelLabel(tt.level); got != tt.want {
			t.Errorf("levelLabel(%v) = %q, want %q", tt.level, got, tt.want)
		}
	}
}

func TestLevelIcon_NotEmpty(t *testing.T) {
	levels := []WatchdogLevel{WatchdogInfo, WatchdogWarn, WatchdogCritical}
	for _, l := range levels {
		if icon := levelIcon(l); icon == "" {
			t.Errorf("levelIcon(%v) returned empty string", l)
		}
	}
}

func TestFormatWatchdogEvents_ContainsDivider(t *testing.T) {
	events := []WatchdogEvent{
		{Bucket: "alpha", Level: WatchdogWarn, Reason: "grew too fast"},
	}
	out := FormatWatchdogEvents(events)
	if !strings.Contains(out, "-") {
		t.Error("expected divider line in output")
	}
	if !strings.Contains(out, "alpha") {
		t.Error("expected bucket name in output")
	}
}

func TestFormatWatchdogEvents_CriticalIcon(t *testing.T) {
	events := []WatchdogEvent{
		{Bucket: "x", Level: WatchdogCritical, Reason: "over limit"},
	}
	out := FormatWatchdogEvents(events)
	if !strings.Contains(out, "CRIT") {
		t.Error("expected CRIT label in output")
	}
}

func TestWatchdogSummary_MultipleCritical(t *testing.T) {
	events := []WatchdogEvent{
		{Level: WatchdogCritical},
		{Level: WatchdogCritical},
		{Level: WatchdogWarn},
	}
	s := WatchdogSummary(events)
	if !strings.Contains(s, "2 critical") {
		t.Errorf("expected '2 critical' in %q", s)
	}
	if !strings.Contains(s, "1 warning") {
		t.Errorf("expected '1 warning' in %q", s)
	}
}
