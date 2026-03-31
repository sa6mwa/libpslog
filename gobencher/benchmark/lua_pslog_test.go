//go:build cgo

package benchmark

import (
	"testing"

	pslog "pkt.systems/pslog"
)

func TestLuaBenchmarkBridge(t *testing.T) {
	if !HasLuaPSLog() {
		t.Skip("lua pslog benchmark support is not available")
	}

	entries := []Entry{
		{
			Level:   pslog.InfoLevel,
			Message: "ready",
			Fields: []Field{
				{Key: "service", Kind: FieldString, String: "api"},
				{Key: "ok", Kind: FieldBool, Bool: true},
				{Key: "attempt", Kind: FieldInt64, Int64: 3},
				{Key: "latency_ms", Kind: FieldFloat64, Float64: 12.5},
				{Key: "big", Kind: FieldUint64, Uint64: ^uint64(0)},
			},
		},
	}
	staticFields := []Field{
		{Key: "component", Kind: FieldString, String: "worker"},
		{Key: "region", Kind: FieldNull},
	}

	data := NewLuaData(entries)
	if data == nil {
		t.Fatal("NewLuaData returned nil")
	}
	defer data.Close()

	logger, err := NewLuaLogger(staticFields, false)
	if err != nil {
		t.Fatalf("NewLuaLogger: %v", err)
	}
	defer logger.Close()

	result, err := RunLua(logger, data, 8)
	if err != nil {
		t.Fatalf("RunLua: %v", err)
	}
	if result.Ops != 8 {
		t.Fatalf("expected 8 ops, got %d", result.Ops)
	}
	if result.Writes == 0 {
		t.Fatal("expected Lua bridge to report writes")
	}
	if result.BytesWritten == 0 {
		t.Fatal("expected Lua bridge to report bytes")
	}
}

func TestLuaBenchmarkBridgeWithMatchesFullEntries(t *testing.T) {
	if !HasLuaPSLog() {
		t.Skip("lua pslog benchmark support is not available")
	}

	fullEntries := []Entry{
		{
			Level:   pslog.InfoLevel,
			Message: "request",
			Fields: []Field{
				{Key: "service", Kind: FieldString, String: "checkout"},
				{Key: "component", Kind: FieldString, String: "worker"},
				{Key: "attempt", Kind: FieldInt64, Int64: 1},
				{Key: "ok", Kind: FieldBool, Bool: true},
			},
		},
		{
			Level:   pslog.WarnLevel,
			Message: "retry",
			Fields: []Field{
				{Key: "service", Kind: FieldString, String: "checkout"},
				{Key: "component", Kind: FieldString, String: "worker"},
				{Key: "attempt", Kind: FieldInt64, Int64: 2},
				{Key: "ok", Kind: FieldBool, Bool: false},
			},
		},
	}
	staticFields, staticKeys := ProductionStaticFields(fullEntries)
	if len(staticFields) != 2 {
		t.Fatalf("expected 2 static fields, got %d", len(staticFields))
	}
	dynamicEntries := ProductionEntriesWithout(fullEntries, staticKeys)

	fullData := NewLuaData(fullEntries)
	if fullData == nil {
		t.Fatal("NewLuaData(fullEntries) returned nil")
	}
	defer fullData.Close()

	dynamicData := NewLuaData(dynamicEntries)
	if dynamicData == nil {
		t.Fatal("NewLuaData(dynamicEntries) returned nil")
	}
	defer dynamicData.Close()

	fullLogger, err := NewLuaLogger(nil, false)
	if err != nil {
		t.Fatalf("NewLuaLogger(full): %v", err)
	}
	defer fullLogger.Close()

	withLogger, err := NewLuaLogger(staticFields, false)
	if err != nil {
		t.Fatalf("NewLuaLogger(with): %v", err)
	}
	defer withLogger.Close()

	fullResult, err := RunLua(fullLogger, fullData, 32)
	if err != nil {
		t.Fatalf("RunLua(full): %v", err)
	}
	withResult, err := RunLua(withLogger, dynamicData, 32)
	if err != nil {
		t.Fatalf("RunLua(with): %v", err)
	}

	if fullResult.BytesWritten != withResult.BytesWritten {
		t.Fatalf("expected with-bytes %d to match full-bytes %d", withResult.BytesWritten, fullResult.BytesWritten)
	}
	if fullResult.Writes != withResult.Writes {
		t.Fatalf("expected with-writes %d to match full-writes %d", withResult.Writes, fullResult.Writes)
	}
	if fullResult.Ops != withResult.Ops {
		t.Fatalf("expected with-ops %d to match full-ops %d", withResult.Ops, fullResult.Ops)
	}
}

func TestLuaPreparedBenchmarkBridgeMatchesRawRun(t *testing.T) {
	if !HasLuaPSLog() {
		t.Skip("lua pslog benchmark support is not available")
	}

	entries := []Entry{
		{
			Level:   pslog.InfoLevel,
			Message: "ready",
			Fields: []Field{
				{Key: "service", Kind: FieldString, String: "api"},
				{Key: "ok", Kind: FieldBool, Bool: true},
				{Key: "attempt", Kind: FieldInt64, Int64: 3},
			},
		},
		{
			Level:   pslog.WarnLevel,
			Message: "retry",
			Fields: []Field{
				{Key: "service", Kind: FieldString, String: "api"},
				{Key: "ok", Kind: FieldBool, Bool: false},
				{Key: "attempt", Kind: FieldUint64, Uint64: ^uint64(0)},
			},
		},
	}

	data := NewLuaData(entries)
	if data == nil {
		t.Fatal("NewLuaData returned nil")
	}
	defer data.Close()

	logger, err := NewLuaLogger(nil, false)
	if err != nil {
		t.Fatalf("NewLuaLogger: %v", err)
	}
	defer logger.Close()

	prepared, err := PrepareLuaData(logger, data)
	if err != nil {
		t.Fatalf("PrepareLuaData: %v", err)
	}
	defer prepared.Close()

	rawResult, err := RunLua(logger, data, 32)
	if err != nil {
		t.Fatalf("RunLua: %v", err)
	}
	preparedResult, err := RunLuaPrepared(logger, prepared, 32)
	if err != nil {
		t.Fatalf("RunLuaPrepared: %v", err)
	}

	if rawResult.BytesWritten != preparedResult.BytesWritten {
		t.Fatalf("expected prepared bytes %d to match raw bytes %d", preparedResult.BytesWritten, rawResult.BytesWritten)
	}
	if rawResult.Writes != preparedResult.Writes {
		t.Fatalf("expected prepared writes %d to match raw writes %d", preparedResult.Writes, rawResult.Writes)
	}
	if rawResult.Ops != preparedResult.Ops {
		t.Fatalf("expected prepared ops %d to match raw ops %d", preparedResult.Ops, rawResult.Ops)
	}
}

func TestLuaPreparedTableBenchmarkBridgeMatchesRawRun(t *testing.T) {
	if !HasLuaPSLog() {
		t.Skip("lua pslog benchmark support is not available")
	}

	entries := []Entry{
		{
			Level:   pslog.InfoLevel,
			Message: "ready",
			Fields: []Field{
				{Key: "service", Kind: FieldString, String: "api"},
				{Key: "ok", Kind: FieldBool, Bool: true},
				{Key: "attempt", Kind: FieldInt64, Int64: 3},
			},
		},
		{
			Level:   pslog.WarnLevel,
			Message: "retry",
			Fields: []Field{
				{Key: "service", Kind: FieldString, String: "api"},
				{Key: "ok", Kind: FieldBool, Bool: false},
				{Key: "attempt", Kind: FieldUint64, Uint64: ^uint64(0)},
			},
		},
	}

	data := NewLuaData(entries)
	if data == nil {
		t.Fatal("NewLuaData returned nil")
	}
	defer data.Close()

	logger, err := NewLuaLogger(nil, false)
	if err != nil {
		t.Fatalf("NewLuaLogger: %v", err)
	}
	defer logger.Close()

	prepared, err := PrepareLuaTableData(logger, data)
	if err != nil {
		t.Fatalf("PrepareLuaTableData: %v", err)
	}
	defer prepared.Close()

	rawResult, err := RunLua(logger, data, 32)
	if err != nil {
		t.Fatalf("RunLua: %v", err)
	}
	preparedResult, err := RunLuaPreparedTable(logger, prepared, 32)
	if err != nil {
		t.Fatalf("RunLuaPreparedTable: %v", err)
	}

	if rawResult.BytesWritten != preparedResult.BytesWritten {
		t.Fatalf("expected prepared table bytes %d to match raw bytes %d", preparedResult.BytesWritten, rawResult.BytesWritten)
	}
	if rawResult.Writes != preparedResult.Writes {
		t.Fatalf("expected prepared table writes %d to match raw writes %d", preparedResult.Writes, rawResult.Writes)
	}
	if rawResult.Ops != preparedResult.Ops {
		t.Fatalf("expected prepared table ops %d to match raw ops %d", preparedResult.Ops, rawResult.Ops)
	}
}

func BenchmarkLuaTableForm(b *testing.B) {
	if !HasLuaPSLog() {
		b.Skip("lua pslog benchmark support is not available")
	}

	b.Run("Production", func(b *testing.B) {
		requireBenchData(b)

		activeData := benchLua
		staticFields := []Field(nil)
		if len(benchStaticFields) > 0 {
			activeData = benchLuaDyn
			staticFields = benchStaticFields
		}

		logger, err := NewLuaLogger(staticFields, true)
		if err != nil {
			b.Fatal(err)
		}
		defer logger.Close()

		prepared, err := PrepareLuaTableData(logger, activeData)
		if err != nil {
			b.Fatal(err)
		}
		defer prepared.Close()

		result, err := RunLuaPreparedTable(logger, prepared, b.N)
		if err != nil {
			b.Fatal(err)
		}
		if result.BytesWritten == 0 {
			b.Fatal("lua table production wrote zero bytes")
		}
		if result.Ops == 0 {
			b.Fatal("lua table production reported zero ops")
		}
		b.ReportMetric(float64(result.Elapsed.Nanoseconds())/float64(result.Ops), "c_ns/op")
		b.ReportMetric(float64(result.BytesWritten)/float64(result.Ops), "bytes/op")
	})

	b.Run("Fixed", func(b *testing.B) {
		entries := makeFixedEntries(256)
		data := NewLuaData(entries)
		if data == nil {
			b.Fatal("NewLuaData returned nil")
		}
		defer data.Close()

		logger, err := NewLuaLogger(nil, true)
		if err != nil {
			b.Fatal(err)
		}
		defer logger.Close()

		prepared, err := PrepareLuaTableData(logger, data)
		if err != nil {
			b.Fatal(err)
		}
		defer prepared.Close()

		result, err := RunLuaPreparedTable(logger, prepared, b.N)
		if err != nil {
			b.Fatal(err)
		}
		if result.BytesWritten == 0 {
			b.Fatal("lua table fixed wrote zero bytes")
		}
		if result.Ops == 0 {
			b.Fatal("lua table fixed reported zero ops")
		}
		b.ReportMetric(float64(result.Elapsed.Nanoseconds())/float64(result.Ops), "c_ns/op")
		b.ReportMetric(float64(result.BytesWritten)/float64(result.Ops), "bytes/op")
	})
}
