package benchmark

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"testing"
	"time"

	pslog "pkt.systems/pslog"
	"pkt.systems/pslog/ansi"
)

func TestParseProductionLine(t *testing.T) {
	entry, err := ParseProductionLine(EmbeddedProductionDataset[0])
	if err != nil {
		t.Fatalf("parse production line: %v", err)
	}
	if entry.Message == "" {
		t.Fatal("expected message")
	}
	if len(entry.Fields) == 0 {
		t.Fatal("expected fields")
	}
	if len(entry.Keyvals) == 0 {
		t.Fatal("expected keyvals")
	}
}

func TestProductionStaticFields(t *testing.T) {
	entries, err := LoadProductionEntries(64)
	if err != nil {
		t.Fatalf("load production entries: %v", err)
	}
	fields, keys := ProductionStaticFields(entries)
	if len(fields) == 0 || len(keys) == 0 {
		t.Fatal("expected static fields")
	}
	filtered := ProductionEntriesWithout(entries, keys)
	if len(filtered) != len(entries) {
		t.Fatalf("filtered entry count mismatch: got %d want %d", len(filtered), len(entries))
	}
}

func TestCLoggerWrites(t *testing.T) {
	prepared := NewCPreparedData(makeFixedEntries(4))
	kvfmt, err := NewCKVFmtData(makeFixedEntries(4))
	raw := NewCRawData(makeFixedEntries(4))
	defer prepared.Close()
	defer kvfmt.Close()
	defer raw.Close()

	logger, err := NewCLogger(Variant{Name: "json", Mode: pslog.ModeStructured})
	if err != nil {
		t.Fatalf("new C logger: %v", err)
	}
	defer logger.Close()

	logger.LogPrepared(prepared.Entries[0])
	if logger.BytesWritten() == 0 {
		t.Fatal("expected C logger to write bytes")
	}

	logger.Reset()
	result, err := logger.RunPrepared(prepared, 8)
	if err != nil {
		t.Fatalf("run prepared: %v", err)
	}
	if result.Ops != 8 {
		t.Fatalf("ops mismatch: got %d want 8", result.Ops)
	}
	if result.BytesWritten == 0 {
		t.Fatal("expected raw run to write bytes")
	}

	logger.Reset()
	rawResult, err := logger.RunRaw(raw, CRawFields{}, 8)
	if err != nil {
		t.Fatalf("run raw semantics: %v", err)
	}
	if rawResult.Ops != 8 {
		t.Fatalf("raw ops mismatch: got %d want 8", rawResult.Ops)
	}
	if rawResult.BytesWritten == 0 {
		t.Fatal("expected raw semantics run to write bytes")
	}

	logger.Reset()
	kvfmtResult, err := logger.RunKVFmt(kvfmt, CPreparedFields{}, 8)
	if err != nil {
		t.Fatalf("run kvfmt semantics: %v", err)
	}
	if kvfmtResult.Ops != 8 {
		t.Fatalf("kvfmt ops mismatch: got %d want 8", kvfmtResult.Ops)
	}
	if kvfmtResult.BytesWritten == 0 {
		t.Fatal("expected kvfmt run to write bytes")
	}
}

func TestCLoggerWithPrepared(t *testing.T) {
	entries := makeFixedEntries(8)
	staticFields := []Field{
		{Key: "app", Kind: FieldString, String: "demo"},
		{Key: "sys", Kind: FieldString, String: "bench"},
	}
	prepared := NewCPreparedData(entries)
	defer prepared.Close()
	withFields := prepared.PrepareFields(staticFields)

	logger, err := NewCLogger(Variant{Name: "json", Mode: pslog.ModeStructured})
	if err != nil {
		t.Fatalf("new C logger: %v", err)
	}
	defer logger.Close()

	withLogger, err := logger.WithPrepared(withFields)
	if err != nil {
		t.Fatalf("with prepared: %v", err)
	}
	defer withLogger.Close()

	withLogger.LogPrepared(prepared.Entries[0])
	if withLogger.BytesWritten() == 0 {
		t.Fatal("expected derived C logger to write bytes")
	}

	withLogger.Reset()
	result, err := withLogger.RunPrepared(prepared, 16)
	if err != nil {
		t.Fatalf("run prepared with derived logger: %v", err)
	}
	if result.Ops != 16 {
		t.Fatalf("ops mismatch: got %d want 16", result.Ops)
	}
	if result.BytesWritten == 0 {
		t.Fatal("expected derived raw run to write bytes")
	}
}

func TestCLoggerPublicWrites(t *testing.T) {
	entries := makeFixedEntries(4)
	staticFields := []Field{
		{Key: "app", Kind: FieldString, String: "demo"},
	}

	logger, err := NewCLogger(Variant{Name: "json", Mode: pslog.ModeStructured})
	if err != nil {
		t.Fatalf("new C logger: %v", err)
	}
	defer logger.Close()

	withLogger, err := logger.With(staticFields)
	if err != nil {
		t.Fatalf("with C logger: %v", err)
	}
	defer withLogger.Close()

	withLogger.Log(entries[0])
	if withLogger.BytesWritten() == 0 {
		t.Fatal("expected public C logger to write bytes")
	}
}

func TestBenchmarkVariantsEmitTimestamps(t *testing.T) {
	entry := makeFixedEntries(1)[0]

	for _, variant := range Variants() {
		t.Run(variant.Name, func(t *testing.T) {
			var goCaptured bytes.Buffer
			goLogger := NewGoLogger(&goCaptured, variant)
			goLogger.Log(entry.Level, entry.Message, entry.Keyvals...)
			assertBenchmarkTimestamp(t, "Go", variant, normalizeCapture(goCaptured.String(), variant.Color))

			cLogger, err := NewCLogger(variant)
			if err != nil {
				t.Fatalf("new C logger: %v", err)
			}
			defer cLogger.Close()
			if err := cLogger.EnableCapture(16 * 1024); err != nil {
				t.Fatalf("enable C capture: %v", err)
			}
			cLogger.Log(entry)
			assertBenchmarkTimestamp(t, "C", variant, normalizeCapture(cLogger.Captured(), variant.Color))
		})
	}
}

func newCParityLogger(variant Variant) (*CLogger, error) {
	return newCLoggerWithPaletteAndTimestamps(variant, "default", false)
}

func newCParityLoggerWithPalette(variant Variant, paletteName string) (*CLogger, error) {
	return newCLoggerWithPaletteAndTimestamps(variant, paletteName, false)
}

func TestGoCNoColorOutputParityFixed(t *testing.T) {
	entries, staticFields := makeGoCParityFixture()
	assertGoCNoColorParity(t, "fixed", entries, staticFields)
}

func assertBenchmarkTimestamp(t *testing.T, impl string, variant Variant, captured string) {
	t.Helper()

	lines := splitCapturedLines(captured)
	if len(lines) == 0 {
		t.Fatalf("%s %s: expected captured output", impl, variant.Name)
	}

	switch variant.Mode {
	case pslog.ModeStructured:
		assertJSONBenchmarkTimestamp(t, impl, variant.Name, lines[0])
	default:
		assertConsoleBenchmarkTimestamp(t, impl, variant.Name, lines[0])
	}
}

func assertJSONBenchmarkTimestamp(t *testing.T, impl, variantName, line string) {
	t.Helper()

	var payload map[string]any
	var tsValue any
	var ok bool

	if err := json.Unmarshal([]byte(line), &payload); err != nil {
		t.Fatalf("%s %s: unmarshal json output: %v", impl, variantName, err)
	}
	tsValue, ok = payload["ts"]
	if !ok {
		t.Fatalf("%s %s: missing ts field", impl, variantName)
	}
	ts, ok := tsValue.(string)
	if !ok || ts == "" {
		t.Fatalf("%s %s: invalid ts field %#v", impl, variantName, tsValue)
	}
	if _, err := time.Parse(time.RFC3339, ts); err != nil {
		t.Fatalf("%s %s: ts is not RFC3339: %q (%v)", impl, variantName, ts, err)
	}
}

func assertConsoleBenchmarkTimestamp(t *testing.T, impl, variantName, line string) {
	t.Helper()

	fields := strings.Fields(line)
	if len(fields) < 2 {
		t.Fatalf("%s %s: console output too short: %q", impl, variantName, line)
	}
	if !isDefaultConsoleTimestamp(fields[0]) {
		t.Fatalf("%s %s: missing default console timestamp prefix: %q", impl, variantName, line)
	}
}

func isDefaultConsoleTimestamp(text string) bool {
	var i int

	if len(text) != 6 {
		return false
	}
	for i = 0; i < len(text); i++ {
		if text[i] < '0' || text[i] > '9' {
			return false
		}
	}
	return true
}

func TestGoCPaletteTablesMatch(t *testing.T) {
	for _, palette := range goCParityPalettes() {
		cPalette, err := CPaletteByName(palette.Name)
		if err != nil {
			t.Fatalf("lookup C palette %s: %v", palette.Name, err)
		}
		assertGoCPaletteEqual(t, palette.Name, cPalette, palette.Go)
	}
}

func TestGoCColorOutputParityFixed(t *testing.T) {
	entries, staticFields := makeGoCParityFixture()

	variants := []Variant{
		{Name: "jsoncolor", Mode: pslog.ModeStructured, Color: true},
		{Name: "consolecolor", Mode: pslog.ModeConsole, Color: true},
	}

	for _, palette := range goCParityPalettes() {
		for _, variant := range variants {
			assertGoCColorParity(t, palette.Name, palette.Go, variant, entries, staticFields)
		}
	}
}

func assertGoCNoColorParity(t *testing.T, label string, entries []Entry, staticFields []Field) {
	t.Helper()

	variants := []Variant{
		{Name: "json", Mode: pslog.ModeStructured, Color: false},
		{Name: "console", Mode: pslog.ModeConsole, Color: false},
	}

	for _, variant := range variants {
		sink := NewSink()
		var captured bytes.Buffer
		sink.SetTee(&captured)
		goLogger := newGoLoggerWithTimestamp(sink, variant, false).With(FieldsToKeyvals(staticFields)...)

		rootCLogger, err := newCParityLogger(variant)
		if err != nil {
			t.Fatalf("%s: new C logger %s: %v", label, variant.Name, err)
		}
		if err := rootCLogger.EnableCapture(256 * 1024); err != nil {
			rootCLogger.Close()
			t.Fatalf("%s: enable capture %s: %v", label, variant.Name, err)
		}
		cLogger, err := rootCLogger.With(staticFields)
		if err != nil {
			rootCLogger.Close()
			t.Fatalf("%s: with C logger %s: %v", label, variant.Name, err)
		}

		for _, entry := range entries {
			entry.LogGo(goLogger)
			cLogger.Log(entry)
		}

		goCaptured := captured.String()
		cCaptured := cLogger.Captured()
		if goCaptured == "" || cCaptured == "" {
			cLogger.Close()
			rootCLogger.Close()
			t.Fatalf("%s: expected captured output for %s", label, variant.Name)
		}
		if !equivalentCapturedOutput(goCaptured, cCaptured, variant.Mode) {
			diff := explainCapturedMismatch(goCaptured, cCaptured, variant.Mode)
			cLogger.Close()
			rootCLogger.Close()
			t.Fatalf("%s: Go/C no-color output mismatch for %s: %s", label, variant.Name, diff)
		}

		cLogger.Close()
		rootCLogger.Close()
	}
}

func makeGoCParityFixture() ([]Entry, []Field) {
	entries := []Entry{
		{
			Level:   pslog.DebugLevel,
			Message: "request started",
			Fields: []Field{
				{Key: "path", Kind: FieldString, String: "/v1/messages"},
				{Key: "cached", Kind: FieldBool, Bool: true},
				{Key: "request_id", Kind: FieldInt64, Int64: 7},
			},
			Keyvals: []any{
				"path", "/v1/messages",
				"cached", true,
				"request_id", int64(7),
			},
		},
		{
			Level:   pslog.InfoLevel,
			Message: "request finished",
			Fields: []Field{
				{Key: "path", Kind: FieldString, String: "/v1/messages"},
				{Key: "cached", Kind: FieldBool, Bool: false},
				{Key: "request_id", Kind: FieldInt64, Int64: 8},
			},
			Keyvals: []any{
				"path", "/v1/messages",
				"cached", false,
				"request_id", int64(8),
			},
		},
		{
			Level:   pslog.WarnLevel,
			Message: "slow request",
			Fields: []Field{
				{Key: "path", Kind: FieldString, String: "/v1/search"},
				{Key: "cached", Kind: FieldBool, Bool: false},
				{Key: "request_id", Kind: FieldInt64, Int64: 9},
			},
			Keyvals: []any{
				"path", "/v1/search",
				"cached", false,
				"request_id", int64(9),
			},
		},
		{
			Level:   pslog.ErrorLevel,
			Message: "request failed",
			Fields: []Field{
				{Key: "path", Kind: FieldString, String: "/v1/search"},
				{Key: "cached", Kind: FieldBool, Bool: false},
				{Key: "request_id", Kind: FieldInt64, Int64: 10},
			},
			Keyvals: []any{
				"path", "/v1/search",
				"cached", false,
				"request_id", int64(10),
			},
		},
	}

	staticFields := []Field{
		{Key: "app", Kind: FieldString, String: "demo"},
		{Key: "component", Kind: FieldString, String: "parity"},
	}

	return entries, staticFields
}

func assertGoCPaletteEqual(t *testing.T, name string, cPalette CPalette, goPalette ansi.Palette) {
	t.Helper()

	if cPalette.Key != goPalette.Key {
		t.Fatalf("%s key mismatch: C=%q Go=%q", name, cPalette.Key, goPalette.Key)
	}
	if cPalette.String != goPalette.String {
		t.Fatalf("%s string mismatch: C=%q Go=%q", name, cPalette.String, goPalette.String)
	}
	if cPalette.Num != goPalette.Num {
		t.Fatalf("%s number mismatch: C=%q Go=%q", name, cPalette.Num, goPalette.Num)
	}
	if cPalette.Bool != goPalette.Bool {
		t.Fatalf("%s boolean mismatch: C=%q Go=%q", name, cPalette.Bool, goPalette.Bool)
	}
	if cPalette.Nil != goPalette.Nil {
		t.Fatalf("%s nil mismatch: C=%q Go=%q", name, cPalette.Nil, goPalette.Nil)
	}
	if cPalette.Trace != goPalette.Trace {
		t.Fatalf("%s trace mismatch: C=%q Go=%q", name, cPalette.Trace, goPalette.Trace)
	}
	if cPalette.Debug != goPalette.Debug {
		t.Fatalf("%s debug mismatch: C=%q Go=%q", name, cPalette.Debug, goPalette.Debug)
	}
	if cPalette.Info != goPalette.Info {
		t.Fatalf("%s info mismatch: C=%q Go=%q", name, cPalette.Info, goPalette.Info)
	}
	if cPalette.Warn != goPalette.Warn {
		t.Fatalf("%s warn mismatch: C=%q Go=%q", name, cPalette.Warn, goPalette.Warn)
	}
	if cPalette.Error != goPalette.Error {
		t.Fatalf("%s error mismatch: C=%q Go=%q", name, cPalette.Error, goPalette.Error)
	}
	if cPalette.Fatal != goPalette.Fatal {
		t.Fatalf("%s fatal mismatch: C=%q Go=%q", name, cPalette.Fatal, goPalette.Fatal)
	}
	if cPalette.Panic != goPalette.Panic {
		t.Fatalf("%s panic mismatch: C=%q Go=%q", name, cPalette.Panic, goPalette.Panic)
	}
	if cPalette.Timestamp != goPalette.Timestamp {
		t.Fatalf("%s timestamp mismatch: C=%q Go=%q", name, cPalette.Timestamp, goPalette.Timestamp)
	}
	if cPalette.Message != goPalette.Message {
		t.Fatalf("%s message mismatch: C=%q Go=%q", name, cPalette.Message, goPalette.Message)
	}
	if cPalette.MessageKey != goPalette.MessageKey {
		t.Fatalf("%s message_key mismatch: C=%q Go=%q", name, cPalette.MessageKey, goPalette.MessageKey)
	}
	if cPalette.Reset != ansi.Reset {
		t.Fatalf("%s reset mismatch: C=%q Go=%q", name, cPalette.Reset, ansi.Reset)
	}
}

func assertGoCColorParity(t *testing.T, paletteName string, palette ansi.Palette, variant Variant, entries []Entry, staticFields []Field) {
	t.Helper()

	var captured bytes.Buffer
	goLogger := pslog.NewWithOptions(nil, &captured, pslog.Options{
		Mode:             variant.Mode,
		DisableTimestamp: true,
		MinLevel:         pslog.TraceLevel,
		ForceColor:       true,
		NoColor:          false,
		Palette:          &palette,
	}).With(FieldsToKeyvals(staticFields)...)

	rootCLogger, err := newCParityLoggerWithPalette(variant, paletteName)
	if err != nil {
		t.Fatalf("%s %s: new C logger: %v", paletteName, variant.Name, err)
	}
	if err := rootCLogger.EnableCapture(256 * 1024); err != nil {
		rootCLogger.Close()
		t.Fatalf("%s %s: enable capture: %v", paletteName, variant.Name, err)
	}
	cLogger, err := rootCLogger.With(staticFields)
	if err != nil {
		rootCLogger.Close()
		t.Fatalf("%s %s: with C logger: %v", paletteName, variant.Name, err)
	}

	for _, entry := range entries {
		entry.LogGo(goLogger)
		cLogger.Log(entry)
	}

	goCaptured := captured.String()
	cCaptured := cLogger.Captured()
	if goCaptured == "" || cCaptured == "" {
		cLogger.Close()
		rootCLogger.Close()
		t.Fatalf("%s %s: expected captured output", paletteName, variant.Name)
	}
	if !equivalentColoredOutput(goCaptured, cCaptured, variant.Mode, palette) {
		diff := explainColoredMismatch(goCaptured, cCaptured, variant.Mode, palette)
		cLogger.Close()
		rootCLogger.Close()
		t.Fatalf("%s %s: Go/C raw color mismatch: %s", paletteName, variant.Name, diff)
	}

	cLogger.Close()
	rootCLogger.Close()
}

type goCParityPalette struct {
	Name string
	Go   ansi.Palette
}

func goCParityPalettes() []goCParityPalette {
	return []goCParityPalette{
		{Name: "default", Go: ansi.PaletteDefault},
		{Name: "outrun-electric", Go: ansi.PaletteOutrunElectric},
		{Name: "iosvkem", Go: ansi.PaletteDoomIosvkem},
		{Name: "gruvbox", Go: ansi.PaletteDoomGruvbox},
		{Name: "dracula", Go: ansi.PaletteDoomDracula},
		{Name: "nord", Go: ansi.PaletteDoomNord},
		{Name: "tokyo-night", Go: ansi.PaletteTokyoNight},
		{Name: "solarized-nightfall", Go: ansi.PaletteSolarizedNightfall},
		{Name: "catppuccin-mocha", Go: ansi.PaletteCatppuccinMocha},
		{Name: "gruvbox-light", Go: ansi.PaletteGruvboxLight},
		{Name: "monokai-vibrant", Go: ansi.PaletteMonokaiVibrant},
		{Name: "one-dark-aurora", Go: ansi.PaletteOneDarkAurora},
		{Name: "synthwave-84", Go: ansi.PaletteSynthwave84},
		{Name: "kanagawa", Go: ansi.PaletteKanagawa},
		{Name: "rose-pine", Go: ansi.PaletteRosePine},
		{Name: "rose-pine-dawn", Go: ansi.PaletteRosePineDawn},
		{Name: "everforest", Go: ansi.PaletteEverforest},
		{Name: "everforest-light", Go: ansi.PaletteEverforestLight},
		{Name: "night-owl", Go: ansi.PaletteNightOwl},
		{Name: "ayu-mirage", Go: ansi.PaletteAyuMirage},
		{Name: "ayu-light", Go: ansi.PaletteAyuLight},
		{Name: "one-light", Go: ansi.PaletteOneLight},
		{Name: "one-dark", Go: ansi.PaletteOneDark},
		{Name: "solarized-light", Go: ansi.PaletteSolarizedLight},
		{Name: "solarized-dark", Go: ansi.PaletteSolarizedDark},
		{Name: "github-light", Go: ansi.PaletteGithubLight},
		{Name: "github-dark", Go: ansi.PaletteGithubDark},
		{Name: "papercolor-light", Go: ansi.PalettePapercolorLight},
		{Name: "papercolor-dark", Go: ansi.PalettePapercolorDark},
		{Name: "oceanic-next", Go: ansi.PaletteOceanicNext},
		{Name: "horizon", Go: ansi.PaletteHorizon},
		{Name: "palenight", Go: ansi.PalettePalenight},
	}
}

func TestCPublicPreparedParityFixed(t *testing.T) {
	entries := makeFixedEntries(16)
	prepared := NewCPreparedData(entries)
	raw := NewCRawData(entries)
	defer prepared.Close()
	defer raw.Close()

	for _, variant := range Variants() {
		publicLogger, err := newCParityLogger(variant)
		if err != nil {
			t.Fatalf("new public C logger %s: %v", variant.Name, err)
		}
		if err := publicLogger.EnableCapture(128 * 1024); err != nil {
			t.Fatalf("enable public capture %s: %v", variant.Name, err)
		}

		preparedLogger, err := newCParityLogger(variant)
		if err != nil {
			publicLogger.Close()
			t.Fatalf("new prepared C logger %s: %v", variant.Name, err)
		}
		if err := preparedLogger.EnableCapture(128 * 1024); err != nil {
			publicLogger.Close()
			preparedLogger.Close()
			t.Fatalf("enable prepared capture %s: %v", variant.Name, err)
		}

		rawLogger, err := newCParityLogger(variant)
		if err != nil {
			publicLogger.Close()
			preparedLogger.Close()
			t.Fatalf("new raw C logger %s: %v", variant.Name, err)
		}
		if err := rawLogger.EnableCapture(128 * 1024); err != nil {
			publicLogger.Close()
			preparedLogger.Close()
			rawLogger.Close()
			t.Fatalf("enable raw capture %s: %v", variant.Name, err)
		}

		for i, entry := range entries {
			publicLogger.Log(entry)
			preparedLogger.LogPrepared(prepared.Entries[i])
		}

		publicCaptured := normalizeCapture(publicLogger.Captured(), variant.Color)
		preparedCaptured := normalizeCapture(preparedLogger.Captured(), variant.Color)
		rawResult, err := rawLogger.RunRaw(raw, CRawFields{}, len(entries))
		if err != nil {
			publicLogger.Close()
			preparedLogger.Close()
			rawLogger.Close()
			t.Fatalf("run raw %s: %v", variant.Name, err)
		}
		if rawResult.BytesWritten == 0 {
			publicLogger.Close()
			preparedLogger.Close()
			rawLogger.Close()
			t.Fatalf("expected raw bytes for %s", variant.Name)
		}
		rawCaptured := normalizeCapture(rawLogger.Captured(), variant.Color)
		if !equivalentCapturedOutput(publicCaptured, preparedCaptured, variant.Mode) {
			publicLogger.Close()
			preparedLogger.Close()
			rawLogger.Close()
			t.Fatalf("public/prepared mismatch for %s", variant.Name)
		}
		if !equivalentCapturedOutput(publicCaptured, rawCaptured, variant.Mode) {
			publicLogger.Close()
			preparedLogger.Close()
			rawLogger.Close()
			t.Fatalf("public/raw mismatch for %s", variant.Name)
		}

		publicLogger.Close()
		preparedLogger.Close()
		rawLogger.Close()
	}
}

func TestCKVFmtFixedOutputParity(t *testing.T) {
	entries := makeFixedEntries(16)
	kvfmt, err := NewCKVFmtData(entries)
	if err != nil {
		t.Fatalf("new kvfmt data: %v", err)
	}
	defer kvfmt.Close()

	for _, variant := range Variants() {
		publicLogger, err := newCParityLogger(variant)
		if err != nil {
			t.Fatalf("new public C logger %s: %v", variant.Name, err)
		}
		if err := publicLogger.EnableCapture(128 * 1024); err != nil {
			t.Fatalf("enable public capture %s: %v", variant.Name, err)
		}

		kvfmtLogger, err := newCParityLogger(variant)
		if err != nil {
			publicLogger.Close()
			t.Fatalf("new kvfmt C logger %s: %v", variant.Name, err)
		}
		if err := kvfmtLogger.EnableCapture(128 * 1024); err != nil {
			publicLogger.Close()
			kvfmtLogger.Close()
			t.Fatalf("enable kvfmt capture %s: %v", variant.Name, err)
		}

		for i, entry := range entries {
			publicLogger.Log(entry)
			if err := kvfmtLogger.LogKVFmt(kvfmt.Entries[i]); err != nil {
				publicLogger.Close()
				kvfmtLogger.Close()
				t.Fatalf("log kvfmt %s: %v", variant.Name, err)
			}
		}

		publicCaptured := normalizeCapture(publicLogger.Captured(), variant.Color)
		kvfmtCaptured := normalizeCapture(kvfmtLogger.Captured(), variant.Color)
		if !equivalentCapturedOutput(publicCaptured, kvfmtCaptured, variant.Mode) {
			publicLogger.Close()
			kvfmtLogger.Close()
			t.Fatalf("public/kvfmt mismatch for %s: %s", variant.Name, explainCapturedMismatch(publicCaptured, kvfmtCaptured, variant.Mode))
		}

		publicLogger.Close()
		kvfmtLogger.Close()
	}
}

func TestCProductionPreparedOutputParity(t *testing.T) {
	entries, err := LoadProductionEntries(64)
	if err != nil {
		t.Fatalf("load production entries: %v", err)
	}
	staticFields, staticKeys := ProductionStaticFields(entries)
	dynamicEntries := ProductionEntriesWithout(entries, staticKeys)
	preparedAll := NewCPreparedData(entries)
	defer preparedAll.Close()
	preparedDynamic := NewCPreparedData(dynamicEntries)
	defer preparedDynamic.Close()
	preparedStatic := preparedDynamic.PrepareFields(staticFields)

	for _, variant := range Variants() {
		fullLogger, err := newCParityLogger(variant)
		if err != nil {
			t.Fatalf("new full C logger %s: %v", variant.Name, err)
		}
		if err := fullLogger.EnableCapture(256 * 1024); err != nil {
			t.Fatalf("enable full capture %s: %v", variant.Name, err)
		}

		withRoot, err := newCParityLogger(variant)
		if err != nil {
			fullLogger.Close()
			t.Fatalf("new with-root C logger %s: %v", variant.Name, err)
		}
		if err := withRoot.EnableCapture(256 * 1024); err != nil {
			fullLogger.Close()
			withRoot.Close()
			t.Fatalf("enable with capture %s: %v", variant.Name, err)
		}
		withLogger, err := withRoot.WithPrepared(preparedStatic)
		if err != nil {
			fullLogger.Close()
			withRoot.Close()
			t.Fatalf("with prepared %s: %v", variant.Name, err)
		}

		for _, entry := range preparedAll.Entries {
			fullLogger.LogPrepared(entry)
		}
		for _, entry := range preparedDynamic.Entries {
			withLogger.LogPrepared(entry)
		}

		fullCaptured := normalizeCapture(fullLogger.Captured(), variant.Color)
		withCaptured := normalizeCapture(withLogger.Captured(), variant.Color)
		if !equivalentCapturedOutput(fullCaptured, withCaptured, variant.Mode) {
			fullLogger.Close()
			withLogger.Close()
			withRoot.Close()
			t.Fatalf("production output mismatch for %s", variant.Name)
		}
		if uint64(len(fullLogger.Captured())) != fullLogger.BytesWritten() {
			fullLogger.Close()
			withLogger.Close()
			withRoot.Close()
			t.Fatalf("byte count mismatch for %s full path", variant.Name)
		}
		if uint64(len(withLogger.Captured())) != withLogger.BytesWritten() {
			fullLogger.Close()
			withLogger.Close()
			withRoot.Close()
			t.Fatalf("byte count mismatch for %s with path", variant.Name)
		}

		fullLogger.Close()
		withLogger.Close()
		withRoot.Close()
	}
}

func TestCProductionPublicOutputParity(t *testing.T) {
	entries, err := LoadProductionEntries(64)
	if err != nil {
		t.Fatalf("load production entries: %v", err)
	}
	staticFields, staticKeys := ProductionStaticFields(entries)
	dynamicEntries := ProductionEntriesWithout(entries, staticKeys)
	rawAll := NewCRawData(entries)
	rawDynamic := NewCRawData(dynamicEntries)
	rawStatic := rawDynamic.PrepareFields(staticFields)
	defer rawAll.Close()
	defer rawDynamic.Close()

	for _, variant := range Variants() {
		fullLogger, err := newCParityLogger(variant)
		if err != nil {
			t.Fatalf("new full public C logger %s: %v", variant.Name, err)
		}
		if err := fullLogger.EnableCapture(256 * 1024); err != nil {
			t.Fatalf("enable full public capture %s: %v", variant.Name, err)
		}

		withRoot, err := newCParityLogger(variant)
		if err != nil {
			fullLogger.Close()
			t.Fatalf("new with-root public C logger %s: %v", variant.Name, err)
		}
		if err := withRoot.EnableCapture(256 * 1024); err != nil {
			fullLogger.Close()
			withRoot.Close()
			t.Fatalf("enable with public capture %s: %v", variant.Name, err)
		}
		withLogger, err := withRoot.With(staticFields)
		if err != nil {
			fullLogger.Close()
			withRoot.Close()
			t.Fatalf("with public fields %s: %v", variant.Name, err)
		}

		for _, entry := range entries {
			fullLogger.Log(entry)
		}
		for _, entry := range dynamicEntries {
			withLogger.Log(entry)
		}

		fullCaptured := normalizeCapture(fullLogger.Captured(), variant.Color)
		withCaptured := normalizeCapture(withLogger.Captured(), variant.Color)
		if !equivalentCapturedOutput(fullCaptured, withCaptured, variant.Mode) {
			fullLogger.Close()
			withLogger.Close()
			withRoot.Close()
			t.Fatalf("public production output mismatch for %s", variant.Name)
		}

		rawLogger, err := newCParityLogger(variant)
		if err != nil {
			fullLogger.Close()
			withLogger.Close()
			withRoot.Close()
			t.Fatalf("new raw C logger %s: %v", variant.Name, err)
		}
		if err := rawLogger.EnableCapture(256 * 1024); err != nil {
			fullLogger.Close()
			withLogger.Close()
			withRoot.Close()
			rawLogger.Close()
			t.Fatalf("enable raw capture %s: %v", variant.Name, err)
		}
		rawResult, err := rawLogger.RunRaw(rawDynamic, rawStatic, len(dynamicEntries))
		if err != nil {
			fullLogger.Close()
			withLogger.Close()
			withRoot.Close()
			rawLogger.Close()
			t.Fatalf("run raw production %s: %v", variant.Name, err)
		}
		if rawResult.BytesWritten == 0 {
			fullLogger.Close()
			withLogger.Close()
			withRoot.Close()
			rawLogger.Close()
			t.Fatalf("expected raw production bytes for %s", variant.Name)
		}
		rawCaptured := normalizeCapture(rawLogger.Captured(), variant.Color)
		if !equivalentCapturedOutput(fullCaptured, rawCaptured, variant.Mode) {
			fullLogger.Close()
			withLogger.Close()
			withRoot.Close()
			rawLogger.Close()
			t.Fatalf("raw production output mismatch for %s", variant.Name)
		}

		fullLogger.Close()
		withLogger.Close()
		withRoot.Close()
		rawLogger.Close()
	}
}

func TestCKVFmtProductionOutputParity(t *testing.T) {
	entries, err := LoadProductionEntries(64)
	if err != nil {
		t.Fatalf("load production entries: %v", err)
	}
	staticFields, staticKeys := ProductionStaticFields(entries)
	dynamicEntries := ProductionEntriesWithout(entries, staticKeys)
	kvfmtDynamic, err := NewCKVFmtData(dynamicEntries)
	if err != nil {
		t.Fatalf("new production kvfmt data: %v", err)
	}
	defer kvfmtDynamic.Close()
	preparedDynamic := NewCPreparedData(dynamicEntries)
	defer preparedDynamic.Close()
	preparedStatic := preparedDynamic.PrepareFields(staticFields)

	for _, variant := range Variants() {
		publicRoot, err := newCParityLogger(variant)
		if err != nil {
			t.Fatalf("new public root %s: %v", variant.Name, err)
		}
		if err := publicRoot.EnableCapture(256 * 1024); err != nil {
			t.Fatalf("enable public root capture %s: %v", variant.Name, err)
		}
		publicLogger, err := publicRoot.With(staticFields)
		if err != nil {
			publicRoot.Close()
			t.Fatalf("with public logger %s: %v", variant.Name, err)
		}

		kvfmtRoot, err := newCParityLogger(variant)
		if err != nil {
			publicLogger.Close()
			publicRoot.Close()
			t.Fatalf("new kvfmt root %s: %v", variant.Name, err)
		}
		if err := kvfmtRoot.EnableCapture(256 * 1024); err != nil {
			publicLogger.Close()
			publicRoot.Close()
			kvfmtRoot.Close()
			t.Fatalf("enable kvfmt root capture %s: %v", variant.Name, err)
		}
		kvfmtLogger, err := kvfmtRoot.WithPrepared(preparedStatic)
		if err != nil {
			publicLogger.Close()
			publicRoot.Close()
			kvfmtRoot.Close()
			t.Fatalf("with prepared kvfmt logger %s: %v", variant.Name, err)
		}

		for i, entry := range dynamicEntries {
			publicLogger.Log(entry)
			if err := kvfmtLogger.LogKVFmt(kvfmtDynamic.Entries[i]); err != nil {
				publicLogger.Close()
				publicRoot.Close()
				kvfmtLogger.Close()
				kvfmtRoot.Close()
				t.Fatalf("log kvfmt production %s: %v", variant.Name, err)
			}
		}

		publicCaptured := normalizeCapture(publicLogger.Captured(), variant.Color)
		kvfmtCaptured := normalizeCapture(kvfmtLogger.Captured(), variant.Color)
		if !equivalentCapturedOutput(publicCaptured, kvfmtCaptured, variant.Mode) {
			publicLogger.Close()
			publicRoot.Close()
			kvfmtLogger.Close()
			kvfmtRoot.Close()
			t.Fatalf("production public/kvfmt mismatch for %s: %s", variant.Name, explainCapturedMismatch(publicCaptured, kvfmtCaptured, variant.Mode))
		}

		runLogger, err := newCParityLogger(variant)
		if err != nil {
			publicLogger.Close()
			publicRoot.Close()
			kvfmtLogger.Close()
			kvfmtRoot.Close()
			t.Fatalf("new kvfmt run logger %s: %v", variant.Name, err)
		}
		if err := runLogger.EnableCapture(256 * 1024); err != nil {
			publicLogger.Close()
			publicRoot.Close()
			kvfmtLogger.Close()
			kvfmtRoot.Close()
			runLogger.Close()
			t.Fatalf("enable kvfmt run capture %s: %v", variant.Name, err)
		}
		result, err := runLogger.RunKVFmt(kvfmtDynamic, preparedStatic, len(dynamicEntries))
		if err != nil {
			publicLogger.Close()
			publicRoot.Close()
			kvfmtLogger.Close()
			kvfmtRoot.Close()
			runLogger.Close()
			t.Fatalf("run kvfmt production %s: %v", variant.Name, err)
		}
		if result.BytesWritten != uint64(len(runLogger.Captured())) {
			publicLogger.Close()
			publicRoot.Close()
			kvfmtLogger.Close()
			kvfmtRoot.Close()
			runLogger.Close()
			t.Fatalf("kvfmt production byte count mismatch for %s", variant.Name)
		}

		publicLogger.Close()
		publicRoot.Close()
		kvfmtLogger.Close()
		kvfmtRoot.Close()
		runLogger.Close()
	}
}

func TestCProductionRunPreparedCountsCapturedBytes(t *testing.T) {
	entries, err := LoadProductionEntries(64)
	if err != nil {
		t.Fatalf("load production entries: %v", err)
	}
	prepared := NewCPreparedData(entries)
	defer prepared.Close()

	logger, err := NewCLogger(Variant{Name: "json", Mode: pslog.ModeStructured})
	if err != nil {
		t.Fatalf("new C logger: %v", err)
	}
	defer logger.Close()
	if err := logger.EnableCapture(256 * 1024); err != nil {
		t.Fatalf("enable capture: %v", err)
	}

	result, err := logger.RunPrepared(prepared, len(prepared.Entries))
	if err != nil {
		t.Fatalf("run prepared: %v", err)
	}
	if result.Ops != uint64(len(prepared.Entries)) {
		t.Fatalf("ops mismatch: got %d want %d", result.Ops, len(prepared.Entries))
	}
	if result.BytesWritten != uint64(len(logger.Captured())) {
		t.Fatalf("captured byte mismatch: got %d want %d", result.BytesWritten, len(logger.Captured()))
	}
}

func normalizeCapture(text string, colored bool) string {
	if !colored {
		return text
	}
	var out strings.Builder
	out.Grow(len(text))
	for i := 0; i < len(text); {
		if text[i] == '\x1b' && i+1 < len(text) && text[i+1] == '[' {
			i += 2
			for i < len(text) && text[i] != 'm' {
				i++
			}
			if i < len(text) {
				i++
			}
			continue
		}
		out.WriteByte(text[i])
		i++
	}
	return out.String()
}

func equivalentCapturedOutput(left, right string, mode pslog.Mode) bool {
	leftLines := splitCapturedLines(left)
	rightLines := splitCapturedLines(right)
	if len(leftLines) != len(rightLines) {
		return false
	}
	for i := range leftLines {
		switch mode {
		case pslog.ModeStructured:
			if canonicalJSONLine(leftLines[i]) != canonicalJSONLine(rightLines[i]) {
				return false
			}
		default:
			if canonicalConsoleLine(leftLines[i]) != canonicalConsoleLine(rightLines[i]) {
				return false
			}
		}
	}
	return true
}

func explainCapturedMismatch(left, right string, mode pslog.Mode) string {
	leftLines := splitCapturedLines(left)
	rightLines := splitCapturedLines(right)
	if len(leftLines) != len(rightLines) {
		return fmt.Sprintf("line count mismatch: left=%d right=%d", len(leftLines), len(rightLines))
	}
	for i := range leftLines {
		leftCanonical := leftLines[i]
		rightCanonical := rightLines[i]
		switch mode {
		case pslog.ModeStructured:
			leftCanonical = canonicalJSONLine(leftLines[i])
			rightCanonical = canonicalJSONLine(rightLines[i])
		default:
			leftCanonical = canonicalConsoleLine(leftLines[i])
			rightCanonical = canonicalConsoleLine(rightLines[i])
		}
		if leftCanonical != rightCanonical {
			return fmt.Sprintf("line %d mismatch: left=%q right=%q", i+1, leftCanonical, rightCanonical)
		}
	}
	return "unknown mismatch"
}

func explainRawCapturedMismatch(left, right string) string {
	leftLines := splitCapturedLines(left)
	rightLines := splitCapturedLines(right)
	if len(leftLines) != len(rightLines) {
		return fmt.Sprintf("line count mismatch: left=%d right=%d", len(leftLines), len(rightLines))
	}
	for i := range leftLines {
		if leftLines[i] != rightLines[i] {
			return fmt.Sprintf("line %d mismatch: left=%q right=%q", i+1, leftLines[i], rightLines[i])
		}
	}
	return "unknown mismatch"
}

func equivalentColoredOutput(left, right string, mode pslog.Mode, palette ansi.Palette) bool {
	leftLines := splitCapturedLines(left)
	rightLines := splitCapturedLines(right)
	if len(leftLines) != len(rightLines) {
		return false
	}
	for i := range leftLines {
		leftLine := leftLines[i]
		rightLine := rightLines[i]
		if mode == pslog.ModeStructured {
			leftLine = canonicalColoredJSONLine(leftLine, palette)
			rightLine = canonicalColoredJSONLine(rightLine, palette)
		}
		if leftLine != rightLine {
			return false
		}
	}
	return true
}

func explainColoredMismatch(left, right string, mode pslog.Mode, palette ansi.Palette) string {
	leftLines := splitCapturedLines(left)
	rightLines := splitCapturedLines(right)
	if len(leftLines) != len(rightLines) {
		return fmt.Sprintf("line count mismatch: left=%d right=%d", len(leftLines), len(rightLines))
	}
	for i := range leftLines {
		leftLine := leftLines[i]
		rightLine := rightLines[i]
		if mode == pslog.ModeStructured {
			leftLine = canonicalColoredJSONLine(leftLine, palette)
			rightLine = canonicalColoredJSONLine(rightLine, palette)
		}
		if leftLine != rightLine {
			return fmt.Sprintf("line %d mismatch: left=%q right=%q", i+1, leftLine, rightLine)
		}
	}
	return "unknown mismatch"
}

func canonicalColoredJSONLine(line string, palette ansi.Palette) string {
	line = canonicalColoredJSONKeyBoundary(line, palette.Key)
	line = canonicalColoredJSONKeyBoundary(line, palette.MessageKey)
	return line
}

func canonicalColoredJSONKeyBoundary(line, color string) string {
	if color == "" {
		return line
	}
	pattern := color + "\""
	var out strings.Builder
	out.Grow(len(line) + 32)
	for i := 0; i < len(line); {
		start := strings.Index(line[i:], pattern)
		if start < 0 {
			out.WriteString(line[i:])
			break
		}
		start += i
		out.WriteString(line[i:start])
		keyStart := start + len(pattern)
		keyEnd := strings.IndexByte(line[keyStart:], '"')
		if keyEnd < 0 {
			out.WriteString(line[start:])
			break
		}
		keyEnd += keyStart
		if keyEnd+2 <= len(line) && line[keyEnd+1] == ':' {
			out.WriteString(line[start : keyEnd+1])
			out.WriteString(ansi.Reset)
			out.WriteByte(':')
			i = keyEnd + 2
			continue
		}
		out.WriteString(line[start : keyEnd+1])
		i = keyEnd + 1
	}
	return out.String()
}

func splitCapturedLines(text string) []string {
	text = strings.TrimSuffix(text, "\n")
	if text == "" {
		return nil
	}
	return strings.Split(text, "\n")
}

func canonicalJSONLine(line string) string {
	var value map[string]any
	if err := json.Unmarshal([]byte(line), &value); err != nil {
		return line
	}
	data, err := json.Marshal(value)
	if err != nil {
		return line
	}
	return string(data)
}

func canonicalConsoleLine(line string) string {
	tokens := splitConsoleTokens(line)
	if len(tokens) <= 2 {
		return line
	}
	fields := append([]string(nil), tokens[2:]...)
	sort.Strings(fields)
	return fmt.Sprintf("%s %s %s", tokens[0], tokens[1], strings.Join(fields, " "))
}

func splitConsoleTokens(line string) []string {
	var tokens []string
	var builder strings.Builder
	inQuotes := false
	escaped := false

	for i := 0; i < len(line); i++ {
		ch := line[i]
		if escaped {
			builder.WriteByte(ch)
			escaped = false
			continue
		}
		if ch == '\\' {
			builder.WriteByte(ch)
			escaped = true
			continue
		}
		if ch == '"' {
			inQuotes = !inQuotes
			builder.WriteByte(ch)
			continue
		}
		if ch == ' ' && !inQuotes {
			if builder.Len() > 0 {
				tokens = append(tokens, builder.String())
				builder.Reset()
			}
			continue
		}
		builder.WriteByte(ch)
	}
	if builder.Len() > 0 {
		tokens = append(tokens, builder.String())
	}
	return tokens
}
