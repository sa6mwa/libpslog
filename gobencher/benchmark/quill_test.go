package benchmark

import (
	"encoding/json"
	"math"
	"strings"
	"testing"

	pslog "pkt.systems/pslog"
)

func TestQuillRenderPreservesJSONShape(t *testing.T) {
	if !HasQuill() {
		t.Skip("quill benchmark support is not enabled")
	}

	entry := Entry{
		Level:   pslog.WarnLevel,
		Message: "render test",
		Fields: []Field{
			{Key: "user", Kind: FieldString, String: "alice"},
			{Key: "cached", Kind: FieldBool, Bool: true},
			{Key: "attempt", Kind: FieldInt64, Int64: 3},
			{Key: "limit", Kind: FieldUint64, Uint64: 9},
			{Key: "ms", Kind: FieldFloat64, Float64: 12.5},
			{Key: "missing", Kind: FieldNull},
		},
	}

	line, err := RenderQuillEntry(entry)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(line, "\n") {
		t.Fatalf("rendered line missing newline: %q", line)
	}

	raw := make(map[string]any)
	if err := json.Unmarshal([]byte(strings.TrimSpace(line)), &raw); err != nil {
		t.Fatalf("rendered line is not valid json: %v\n%s", err, line)
	}

	if raw["lvl"] != "warn" {
		t.Fatalf("unexpected lvl: %#v", raw["lvl"])
	}
	if raw["msg"] != "render test" {
		t.Fatalf("unexpected msg: %#v", raw["msg"])
	}
	if raw["user"] != "alice" {
		t.Fatalf("unexpected user: %#v", raw["user"])
	}
	if raw["cached"] != true {
		t.Fatalf("unexpected cached: %#v", raw["cached"])
	}
	if raw["attempt"] != float64(3) {
		t.Fatalf("unexpected attempt: %#v", raw["attempt"])
	}
	if raw["limit"] != float64(9) {
		t.Fatalf("unexpected limit: %#v", raw["limit"])
	}
	if raw["ms"] != 12.5 {
		t.Fatalf("unexpected ms: %#v", raw["ms"])
	}
	if value, ok := raw["missing"]; !ok || value != nil {
		t.Fatalf("unexpected missing: %#v", raw["missing"])
	}
	if _, ok := raw["ts"].(string); !ok {
		t.Fatalf("unexpected ts: %#v", raw["ts"])
	}
}

func TestQuillRenderSupportsZeroAndSingleFieldEntries(t *testing.T) {
	if !HasQuill() {
		t.Skip("quill benchmark support is not enabled")
	}

	cases := []struct {
		name  string
		entry Entry
		check func(*testing.T, map[string]any)
	}{
		{
			name: "no fields",
			entry: Entry{
				Level:   pslog.InfoLevel,
				Message: "just a message",
			},
			check: func(t *testing.T, raw map[string]any) {
				t.Helper()
				if raw["lvl"] != "info" {
					t.Fatalf("unexpected lvl: %#v", raw["lvl"])
				}
				if raw["msg"] != "just a message" {
					t.Fatalf("unexpected msg: %#v", raw["msg"])
				}
				if len(raw) != 3 {
					t.Fatalf("unexpected field count: %#v", raw)
				}
			},
		},
		{
			name: "single string field",
			entry: Entry{
				Level:   pslog.DebugLevel,
				Message: "single field",
				Fields: []Field{
					{Key: "user", Kind: FieldString, String: "alice"},
				},
			},
			check: func(t *testing.T, raw map[string]any) {
				t.Helper()
				if raw["lvl"] != "debug" {
					t.Fatalf("unexpected lvl: %#v", raw["lvl"])
				}
				if raw["user"] != "alice" {
					t.Fatalf("unexpected user: %#v", raw["user"])
				}
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			line, err := RenderQuillEntry(tc.entry)
			if err != nil {
				t.Fatal(err)
			}

			raw := make(map[string]any)
			if err := json.Unmarshal([]byte(strings.TrimSpace(line)), &raw); err != nil {
				t.Fatalf("rendered line is not valid json: %v\n%s", err, line)
			}
			tc.check(t, raw)
		})
	}
}

func TestQuillRenderPreservesFatalAndPanicLevels(t *testing.T) {
	if !HasQuill() {
		t.Skip("quill benchmark support is not enabled")
	}

	cases := []struct {
		name     string
		level    pslog.Level
		expected string
	}{
		{name: "fatal", level: pslog.FatalLevel, expected: "fatal"},
		{name: "panic", level: pslog.PanicLevel, expected: "panic"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			line, err := RenderQuillEntry(Entry{
				Level:   tc.level,
				Message: tc.name,
				Fields: []Field{
					{Key: "code", Kind: FieldInt64, Int64: 42},
				},
			})
			if err != nil {
				t.Fatal(err)
			}

			raw := make(map[string]any)
			if err := json.Unmarshal([]byte(strings.TrimSpace(line)), &raw); err != nil {
				t.Fatalf("rendered line is not valid json: %v\n%s", err, line)
			}
			if raw["lvl"] != tc.expected {
				t.Fatalf("unexpected lvl: %#v", raw["lvl"])
			}
		})
	}
}

func TestQuillRejectsNonFiniteFloatFields(t *testing.T) {
	if !HasQuill() {
		t.Skip("quill benchmark support is not enabled")
	}

	cases := []struct {
		name  string
		value float64
	}{
		{name: "nan", value: math.NaN()},
		{name: "pos_inf", value: math.Inf(1)},
		{name: "neg_inf", value: math.Inf(-1)},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			_, err := RenderQuillEntry(Entry{
				Level:   pslog.InfoLevel,
				Message: "non-finite",
				Fields: []Field{
					{Key: "value", Kind: FieldFloat64, Float64: tc.value},
				},
			})
			if err == nil {
				t.Fatal("expected render to reject non-finite float")
			}
			if !strings.Contains(err.Error(), "non-finite float64") {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestQuillRenderRejectsOversizedOutput(t *testing.T) {
	if !HasQuill() {
		t.Skip("quill benchmark support is not enabled")
	}

	_, err := RenderQuillEntry(Entry{
		Level:   pslog.InfoLevel,
		Message: strings.Repeat("m", 9000),
	})
	if err == nil {
		t.Fatal("expected oversized quill render to fail")
	}
	if !strings.Contains(err.Error(), "render failed") {
		t.Fatalf("unexpected error: %v", err)
	}
}
