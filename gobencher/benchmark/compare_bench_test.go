package benchmark

import (
	"fmt"
	"sync"
	"testing"

	pslog "pkt.systems/pslog"
)

var (
	benchDataOnce      sync.Once
	benchEntries       []Entry
	benchDynamic       []Entry
	benchKVFmtAll      *CKVFmtData
	benchKVFmtDyn      *CKVFmtData
	benchStaticFields  []Field
	benchPrepared      *CPreparedData
	benchPreparedDyn   *CPreparedData
	benchPreparedWith  CPreparedFields
	benchRaw           *CRawData
	benchRawDyn        *CRawData
	benchRawWith       CRawFields
	benchLua           *LuaData
	benchLuaDyn        *LuaData
	benchLiblogger     *LibloggerData
	benchQuill         *QuillData
	benchDataLoadError error
)

func loadBenchData() {
	benchEntries, benchDataLoadError = LoadProductionEntries(DefaultProductionLimit)
	if benchDataLoadError != nil {
		return
	}
	benchStaticFields, staticKeys := ProductionStaticFields(benchEntries)
	benchDynamic = ProductionEntriesWithout(benchEntries, staticKeys)
	benchPrepared = NewCPreparedData(benchEntries)
	benchPreparedDyn = NewCPreparedData(benchDynamic)
	benchPreparedWith = benchPreparedDyn.PrepareFields(benchStaticFields)
	benchKVFmtAll, benchDataLoadError = NewCKVFmtData(benchEntries)
	if benchDataLoadError != nil {
		return
	}
	benchKVFmtDyn, benchDataLoadError = NewCKVFmtData(benchDynamic)
	if benchDataLoadError != nil {
		return
	}
	benchRaw = NewCRawData(benchEntries)
	benchRawDyn = NewCRawData(benchDynamic)
	benchRawWith = benchRawDyn.PrepareFields(benchStaticFields)
	if HasLuaPSLog() {
		benchLua = NewLuaData(benchEntries)
		benchLuaDyn = NewLuaData(benchDynamic)
	}
	if HasLiblogger() {
		benchLiblogger = NewLibloggerData(benchEntries)
	}
	if HasQuill() {
		benchQuill = NewQuillData(benchEntries)
		if benchQuill != nil && benchQuill.initErr != nil {
			benchDataLoadError = benchQuill.initErr
			return
		}
	}
}

func requireBenchData(b *testing.B) {
	benchDataOnce.Do(loadBenchData)
	if benchDataLoadError != nil {
		b.Fatal(benchDataLoadError)
	}
	if len(benchEntries) == 0 {
		b.Fatal("no production entries loaded")
	}
}

func BenchmarkProductionCompare(b *testing.B) {
	requireBenchData(b)

	for _, variant := range Variants() {
		variant := variant

		b.Run(fmt.Sprintf("%sGo", variant.Name), func(b *testing.B) {
			sink := NewSink()
			logger := NewGoLogger(sink, variant)
			if len(benchStaticFields) > 0 {
				logger = logger.With(FieldsToKeyvals(benchStaticFields)...)
			}
			activeEntries := benchDynamic
			if len(benchStaticFields) == 0 {
				activeEntries = benchEntries
			}

			sink.Reset()
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				activeEntries[i%len(activeEntries)].LogGo(logger)
			}
			if sink.BytesWritten() == 0 {
				b.Fatalf("%sGo wrote zero bytes", variant.Name)
			}
			b.ReportMetric(float64(sink.BytesWritten())/float64(b.N), "bytes/op")
		})

		b.Run(fmt.Sprintf("%sCffi", variant.Name), func(b *testing.B) {
			logger, err := NewCBenchLogger(variant)
			if err != nil {
				b.Fatal(err)
			}
			defer logger.Close()

			activeEntries := benchDynamic
			activeLogger := logger
			if len(benchStaticFields) > 0 {
				withLogger, err := logger.With(benchStaticFields)
				if err != nil {
					b.Fatal(err)
				}
				defer withLogger.Close()
				activeLogger = withLogger
			} else {
				activeEntries = benchEntries
			}

			activeLogger.Reset()
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				activeLogger.Log(activeEntries[i%len(activeEntries)])
			}
			if activeLogger.BytesWritten() == 0 {
				b.Fatalf("%sCffi wrote zero bytes", variant.Name)
			}
			b.ReportMetric(float64(activeLogger.BytesWritten())/float64(b.N), "bytes/op")
		})

		b.Run(fmt.Sprintf("%sCprepared", variant.Name), func(b *testing.B) {
			logger, err := NewCBenchLogger(variant)
			if err != nil {
				b.Fatal(err)
			}
			defer logger.Close()

			activeEntries := benchPreparedDyn.Entries
			activeLogger := logger
			if len(benchStaticFields) > 0 {
				withLogger, err := logger.WithPrepared(benchPreparedWith)
				if err != nil {
					b.Fatal(err)
				}
				defer withLogger.Close()
				activeLogger = withLogger
			} else {
				activeEntries = benchPrepared.Entries
			}

			activeLogger.Reset()
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				activeLogger.LogPrepared(activeEntries[i%len(activeEntries)])
			}
			if activeLogger.BytesWritten() == 0 {
				b.Fatalf("%sCffi wrote zero bytes", variant.Name)
			}
			b.ReportMetric(float64(activeLogger.BytesWritten())/float64(b.N), "bytes/op")
		})

		b.Run(fmt.Sprintf("%sC", variant.Name), func(b *testing.B) {
			logger, err := NewCBenchLogger(variant)
			if err != nil {
				b.Fatal(err)
			}
			defer logger.Close()

			activeEntries := benchPreparedDyn
			activeLogger := logger
			if len(benchStaticFields) > 0 {
				withLogger, err := logger.WithPrepared(benchPreparedWith)
				if err != nil {
					b.Fatal(err)
				}
				defer withLogger.Close()
				activeLogger = withLogger
			} else {
				activeEntries = benchPrepared
			}

			b.ReportAllocs()
			result, runErr := activeLogger.RunPrepared(activeEntries, b.N)
			if runErr != nil {
				b.Fatal(runErr)
			}
			if result.BytesWritten == 0 {
				b.Fatalf("%sC wrote zero bytes", variant.Name)
			}
			if result.Ops == 0 {
				b.Fatalf("%sC reported zero ops", variant.Name)
			}
			b.ReportMetric(float64(result.Elapsed.Nanoseconds())/float64(result.Ops), "c_ns/op")
			b.ReportMetric(float64(result.BytesWritten)/float64(result.Ops), "bytes/op")
		})

		b.Run(fmt.Sprintf("%sCkvfmt", variant.Name), func(b *testing.B) {
			logger, err := NewCBenchLogger(variant)
			if err != nil {
				b.Fatal(err)
			}
			defer logger.Close()

			activeData := benchKVFmtAll
			activeStatic := CPreparedFields{}
			if len(benchStaticFields) > 0 {
				activeData = benchKVFmtDyn
				activeStatic = benchPreparedWith
			}

			b.ReportAllocs()
			result, runErr := logger.RunKVFmt(activeData, activeStatic, b.N)
			if runErr != nil {
				b.Fatal(runErr)
			}
			if result.BytesWritten == 0 {
				b.Fatalf("%sCkvfmt wrote zero bytes", variant.Name)
			}
			if result.Ops == 0 {
				b.Fatalf("%sCkvfmt reported zero ops", variant.Name)
			}
			b.ReportMetric(float64(result.Elapsed.Nanoseconds())/float64(result.Ops), "c_ns/op")
			b.ReportMetric(float64(result.BytesWritten)/float64(result.Ops), "bytes/op")
		})

		b.Run(fmt.Sprintf("%sCraw", variant.Name), func(b *testing.B) {
			logger, err := NewCBenchLogger(variant)
			if err != nil {
				b.Fatal(err)
			}
			defer logger.Close()

			activeData := benchRaw
			activeStatic := CRawFields{}
			if len(benchStaticFields) > 0 {
				activeData = benchRawDyn
				activeStatic = benchRawWith
			}

			result, runErr := logger.RunRaw(activeData, activeStatic, b.N)
			if runErr != nil {
				b.Fatal(runErr)
			}
			if result.BytesWritten == 0 {
				b.Fatalf("%sCraw wrote zero bytes", variant.Name)
			}
			if result.Ops == 0 {
				b.Fatalf("%sCraw reported zero ops", variant.Name)
			}
			b.ReportMetric(float64(result.Elapsed.Nanoseconds())/float64(result.Ops), "c_ns/op")
			b.ReportMetric(float64(result.BytesWritten)/float64(result.Ops), "bytes/op")
		})
	}

	if HasLuaPSLog() {
		b.Run("jsonLua", func(b *testing.B) {
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

			prepared, err := PrepareLuaData(logger, activeData)
			if err != nil {
				b.Fatal(err)
			}
			defer prepared.Close()

			result, err := RunLuaPrepared(logger, prepared, b.N)
			if err != nil {
				b.Fatal(err)
			}
			if result.BytesWritten == 0 {
				b.Fatal("jsonLua wrote zero bytes")
			}
			if result.Ops == 0 {
				b.Fatal("jsonLua reported zero ops")
			}
			b.ReportMetric(float64(result.Elapsed.Nanoseconds())/float64(result.Ops), "c_ns/op")
			b.ReportMetric(float64(result.BytesWritten)/float64(result.Ops), "bytes/op")
		})
	}

	if HasLiblogger() {
		b.Run("jsonLiblogger", func(b *testing.B) {
			result, err := RunLiblogger(benchLiblogger, b.N)
			if err != nil {
				b.Fatal(err)
			}
			if result.BytesWritten == 0 {
				b.Fatal("jsonLiblogger wrote zero bytes")
			}
			if result.Ops == 0 {
				b.Fatal("jsonLiblogger reported zero ops")
			}
			b.ReportMetric(float64(result.Elapsed.Nanoseconds())/float64(result.Ops), "c_ns/op")
			b.ReportMetric(float64(result.BytesWritten)/float64(result.Ops), "bytes/op")
		})
	}
	if HasQuill() {
		b.Run("jsonQuill", func(b *testing.B) {
			result, err := RunQuill(benchQuill, b.N)
			if err != nil {
				b.Fatal(err)
			}
			if result.BytesWritten == 0 {
				b.Fatal("jsonQuill wrote zero bytes")
			}
			if result.Ops == 0 {
				b.Fatal("jsonQuill reported zero ops")
			}
			b.ReportMetric(float64(result.Elapsed.Nanoseconds())/float64(result.Ops), "c_ns/op")
			b.ReportMetric(float64(result.BytesWritten)/float64(result.Ops), "bytes/op")
		})
	}
}

func BenchmarkFixedCompare(b *testing.B) {
	entries := makeFixedEntries(256)
	prepared := NewCPreparedData(entries)
	kvfmt, err := NewCKVFmtData(entries)
	raw := NewCRawData(entries)
	luaData := NewLuaData(entries)
	libloggerData := NewLibloggerData(entries)
	quillData := NewQuillData(entries)
	defer prepared.Close()
	defer kvfmt.Close()
	defer raw.Close()
	if luaData != nil {
		defer luaData.Close()
	}
	if libloggerData != nil {
		defer libloggerData.Close()
	}
	if quillData != nil {
		defer quillData.Close()
	}
	if err != nil {
		b.Fatal(err)
	}

	for _, variant := range Variants() {
		variant := variant

		b.Run(fmt.Sprintf("%sGo", variant.Name), func(b *testing.B) {
			sink := NewSink()
			logger := NewGoLogger(sink, variant)
			sink.Reset()
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				entries[i%len(entries)].LogGo(logger)
			}
			if sink.BytesWritten() == 0 {
				b.Fatalf("%sGo wrote zero bytes", variant.Name)
			}
			b.ReportMetric(float64(sink.BytesWritten())/float64(b.N), "bytes/op")
		})

		b.Run(fmt.Sprintf("%sCffi", variant.Name), func(b *testing.B) {
			logger, err := NewCBenchLogger(variant)
			if err != nil {
				b.Fatal(err)
			}
			defer logger.Close()

			logger.Reset()
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				logger.Log(entries[i%len(entries)])
			}
			if logger.BytesWritten() == 0 {
				b.Fatalf("%sCffi wrote zero bytes", variant.Name)
			}
			b.ReportMetric(float64(logger.BytesWritten())/float64(b.N), "bytes/op")
		})

		b.Run(fmt.Sprintf("%sCprepared", variant.Name), func(b *testing.B) {
			logger, err := NewCBenchLogger(variant)
			if err != nil {
				b.Fatal(err)
			}
			defer logger.Close()

			logger.Reset()
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				logger.LogPrepared(prepared.Entries[i%len(prepared.Entries)])
			}
			if logger.BytesWritten() == 0 {
				b.Fatalf("%sCffi wrote zero bytes", variant.Name)
			}
			b.ReportMetric(float64(logger.BytesWritten())/float64(b.N), "bytes/op")
		})

		b.Run(fmt.Sprintf("%sC", variant.Name), func(b *testing.B) {
			logger, err := NewCBenchLogger(variant)
			if err != nil {
				b.Fatal(err)
			}
			defer logger.Close()

			b.ReportAllocs()
			result, runErr := logger.RunPrepared(prepared, b.N)
			if runErr != nil {
				b.Fatal(runErr)
			}
			if result.BytesWritten == 0 {
				b.Fatalf("%sC wrote zero bytes", variant.Name)
			}
			if result.Ops == 0 {
				b.Fatalf("%sC reported zero ops", variant.Name)
			}
			b.ReportMetric(float64(result.Elapsed.Nanoseconds())/float64(result.Ops), "c_ns/op")
			b.ReportMetric(float64(result.BytesWritten)/float64(result.Ops), "bytes/op")
		})

		b.Run(fmt.Sprintf("%sCkvfmt", variant.Name), func(b *testing.B) {
			logger, err := NewCBenchLogger(variant)
			if err != nil {
				b.Fatal(err)
			}
			defer logger.Close()

			b.ReportAllocs()
			result, runErr := logger.RunKVFmt(kvfmt, CPreparedFields{}, b.N)
			if runErr != nil {
				b.Fatal(runErr)
			}
			if result.BytesWritten == 0 {
				b.Fatalf("%sCkvfmt wrote zero bytes", variant.Name)
			}
			if result.Ops == 0 {
				b.Fatalf("%sCkvfmt reported zero ops", variant.Name)
			}
			b.ReportMetric(float64(result.Elapsed.Nanoseconds())/float64(result.Ops), "c_ns/op")
			b.ReportMetric(float64(result.BytesWritten)/float64(result.Ops), "bytes/op")
		})

		b.Run(fmt.Sprintf("%sCraw", variant.Name), func(b *testing.B) {
			logger, err := NewCBenchLogger(variant)
			if err != nil {
				b.Fatal(err)
			}
			defer logger.Close()

			result, runErr := logger.RunRaw(raw, CRawFields{}, b.N)
			if runErr != nil {
				b.Fatal(runErr)
			}
			if result.BytesWritten == 0 {
				b.Fatalf("%sCraw wrote zero bytes", variant.Name)
			}
			if result.Ops == 0 {
				b.Fatalf("%sCraw reported zero ops", variant.Name)
			}
			b.ReportMetric(float64(result.Elapsed.Nanoseconds())/float64(result.Ops), "c_ns/op")
			b.ReportMetric(float64(result.BytesWritten)/float64(result.Ops), "bytes/op")
		})
	}

	if HasLuaPSLog() {
		b.Run("jsonLua", func(b *testing.B) {
			logger, err := NewLuaLogger(nil, true)
			if err != nil {
				b.Fatal(err)
			}
			defer logger.Close()

			prepared, err := PrepareLuaData(logger, luaData)
			if err != nil {
				b.Fatal(err)
			}
			defer prepared.Close()

			result, err := RunLuaPrepared(logger, prepared, b.N)
			if err != nil {
				b.Fatal(err)
			}
			if result.BytesWritten == 0 {
				b.Fatal("jsonLua wrote zero bytes")
			}
			if result.Ops == 0 {
				b.Fatal("jsonLua reported zero ops")
			}
			b.ReportMetric(float64(result.Elapsed.Nanoseconds())/float64(result.Ops), "c_ns/op")
			b.ReportMetric(float64(result.BytesWritten)/float64(result.Ops), "bytes/op")
		})
	}

	if HasLiblogger() {
		b.Run("jsonLiblogger", func(b *testing.B) {
			result, err := RunLiblogger(libloggerData, b.N)
			if err != nil {
				b.Fatal(err)
			}
			if result.BytesWritten == 0 {
				b.Fatal("jsonLiblogger wrote zero bytes")
			}
			if result.Ops == 0 {
				b.Fatal("jsonLiblogger reported zero ops")
			}
			b.ReportMetric(float64(result.Elapsed.Nanoseconds())/float64(result.Ops), "c_ns/op")
			b.ReportMetric(float64(result.BytesWritten)/float64(result.Ops), "bytes/op")
		})
	}
	if HasQuill() {
		b.Run("jsonQuill", func(b *testing.B) {
			result, err := RunQuill(quillData, b.N)
			if err != nil {
				b.Fatal(err)
			}
			if result.BytesWritten == 0 {
				b.Fatal("jsonQuill wrote zero bytes")
			}
			if result.Ops == 0 {
				b.Fatal("jsonQuill reported zero ops")
			}
			b.ReportMetric(float64(result.Elapsed.Nanoseconds())/float64(result.Ops), "c_ns/op")
			b.ReportMetric(float64(result.BytesWritten)/float64(result.Ops), "bytes/op")
		})
	}
}

func makeFixedEntries(count int) []Entry {
	if count <= 0 {
		count = 1
	}
	entries := make([]Entry, count)
	for i := 0; i < count; i++ {
		entries[i] = Entry{
			Level:   pslog.InfoLevel,
			Message: "bench",
			Fields: []Field{
				{Key: "request_id", Kind: FieldInt64, Int64: int64(i)},
				{Key: "path", Kind: FieldString, String: "/v1/messages"},
				{Key: "cached", Kind: FieldBool, Bool: i%2 == 0},
				{Key: "ms", Kind: FieldFloat64, Float64: 1.25 + float64(i%7)},
			},
			Keyvals: []any{
				"request_id", int64(i),
				"path", "/v1/messages",
				"cached", i%2 == 0,
				"ms", 1.25 + float64(i%7),
			},
		}
	}
	return entries
}
