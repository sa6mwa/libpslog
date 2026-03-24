package benchmark

/*
#cgo CFLAGS: -I${SRCDIR}/../../bench
#cgo LDFLAGS: ${SRCDIR}/../../build/host/libpslog_bench_liblogger.so -Wl,-rpath,${SRCDIR}/../../build/host

#include <stdlib.h>
#include "bench_liblogger.h"
*/
import "C"

import (
	"fmt"
	"time"
	"unsafe"
)

type LibloggerData struct {
	entries    *C.pslog_bench_liblogger_entry_spec
	fields     *C.pslog_bench_liblogger_field_spec
	entryCount int
	fieldCount int
	strings    []*C.char
}

type LibloggerResult struct {
	Elapsed      time.Duration
	BytesWritten uint64
	Ops          uint64
}

func HasLiblogger() bool {
	return C.pslog_bench_liblogger_available() != 0
}

func libloggerComparisonNames() []string {
	if !HasLiblogger() {
		return nil
	}
	return []string{"jsonLiblogger"}
}

func cLibloggerLevel(level int) C.pslog_bench_liblogger_level {
	return C.pslog_bench_liblogger_level(level)
}

func cLibloggerKind(kind FieldKind) C.pslog_bench_liblogger_field_kind {
	switch kind {
	case FieldString:
		return C.PSLOG_BENCH_LIBLOGGER_FIELD_STRING
	case FieldBool:
		return C.PSLOG_BENCH_LIBLOGGER_FIELD_BOOL
	case FieldInt64:
		return C.PSLOG_BENCH_LIBLOGGER_FIELD_INT64
	case FieldUint64:
		return C.PSLOG_BENCH_LIBLOGGER_FIELD_UINT64
	case FieldFloat64:
		return C.PSLOG_BENCH_LIBLOGGER_FIELD_FLOAT64
	default:
		return C.PSLOG_BENCH_LIBLOGGER_FIELD_NULL
	}
}

func NewLibloggerData(entries []Entry) *LibloggerData {
	if !HasLiblogger() {
		return nil
	}
	if len(entries) == 0 {
		return &LibloggerData{}
	}

	totalFields := 0
	for _, entry := range entries {
		totalFields += len(entry.Fields)
	}

	data := &LibloggerData{
		entryCount: len(entries),
		fieldCount: totalFields,
	}
	data.entries = (*C.pslog_bench_liblogger_entry_spec)(C.calloc(C.size_t(len(entries)), C.size_t(C.sizeof_pslog_bench_liblogger_entry_spec)))
	if totalFields > 0 {
		data.fields = (*C.pslog_bench_liblogger_field_spec)(C.calloc(C.size_t(totalFields), C.size_t(C.sizeof_pslog_bench_liblogger_field_spec)))
	}
	if data.entries == nil || (totalFields > 0 && data.fields == nil) {
		data.Close()
		return nil
	}

	entrySlice := unsafe.Slice(data.entries, len(entries))
	fieldSlice := unsafe.Slice(data.fields, totalFields)
	fieldOffset := 0
	for i, entry := range entries {
		msg := C.CString(entry.Message)
		data.strings = append(data.strings, msg)
		entrySlice[i].level = cLibloggerLevel(int(entry.Level))
		entrySlice[i].message = msg
		entrySlice[i].field_count = C.size_t(len(entry.Fields))
		if len(entry.Fields) > 0 {
			entrySlice[i].fields = &fieldSlice[fieldOffset]
		}
		for _, field := range entry.Fields {
			dst := &fieldSlice[fieldOffset]
			key := C.CString(field.Key)
			data.strings = append(data.strings, key)
			dst.key = key
			dst.kind = cLibloggerKind(field.Kind)
			switch field.Kind {
			case FieldString:
				value := C.CString(field.String)
				data.strings = append(data.strings, value)
				dst.string_value = value
			case FieldBool:
				if field.Bool {
					dst.bool_value = 1
				} else {
					dst.bool_value = 0
				}
			case FieldInt64:
				dst.signed_value = C.long(field.Int64)
			case FieldUint64:
				dst.unsigned_value = C.ulong(field.Uint64)
			case FieldFloat64:
				dst.double_value = C.double(field.Float64)
			default:
			}
			fieldOffset++
		}
	}

	return data
}

func (d *LibloggerData) Close() {
	if d == nil {
		return
	}
	for _, text := range d.strings {
		C.free(unsafe.Pointer(text))
	}
	d.strings = nil
	if d.fields != nil {
		C.free(unsafe.Pointer(d.fields))
		d.fields = nil
	}
	if d.entries != nil {
		C.free(unsafe.Pointer(d.entries))
		d.entries = nil
	}
	d.entryCount = 0
	d.fieldCount = 0
}

func RunLiblogger(data *LibloggerData, iterations int) (LibloggerResult, error) {
	var result C.pslog_bench_liblogger_result

	if !HasLiblogger() {
		return LibloggerResult{}, fmt.Errorf("liblogger benchmark support is not enabled")
	}
	if data == nil || data.entries == nil || data.entryCount == 0 {
		return LibloggerResult{}, fmt.Errorf("liblogger benchmark data not initialized")
	}
	if rc := C.pslog_bench_liblogger_run(data.entries, C.size_t(data.entryCount), C.ulonglong(iterations), &result); rc != 0 {
		return LibloggerResult{}, fmt.Errorf("liblogger benchmark run failed")
	}
	return LibloggerResult{
		Elapsed:      time.Duration(result.elapsed_ns),
		BytesWritten: uint64(result.bytes),
		Ops:          uint64(result.ops),
	}, nil
}
