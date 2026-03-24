package benchmark

/*
#cgo CFLAGS: -I${SRCDIR}/../../bench
#cgo LDFLAGS: ${SRCDIR}/../../build/host/libpslog_bench_quill.so -Wl,-rpath,${SRCDIR}/../../build/host

#include <stdlib.h>
#include "bench_quill.h"
*/
import "C"

import (
	"fmt"
	"math"
	"time"
	"unsafe"
)

type QuillData struct {
	entries    *C.pslog_bench_quill_entry_spec
	fields     *C.pslog_bench_quill_field_spec
	entryCount int
	fieldCount int
	strings    []*C.char
	initErr    error
}

type QuillResult struct {
	Elapsed      time.Duration
	BytesWritten uint64
	Ops          uint64
}

func HasQuill() bool {
	return C.pslog_bench_quill_available() != 0
}

func quillComparisonNames() []string {
	if !HasQuill() {
		return nil
	}
	return []string{"jsonQuill"}
}

func cQuillLevel(level int) C.pslog_bench_quill_level {
	return C.pslog_bench_quill_level(level)
}

func cQuillKind(kind FieldKind) C.pslog_bench_quill_field_kind {
	switch kind {
	case FieldString:
		return C.PSLOG_BENCH_QUILL_FIELD_STRING
	case FieldBool:
		return C.PSLOG_BENCH_QUILL_FIELD_BOOL
	case FieldInt64:
		return C.PSLOG_BENCH_QUILL_FIELD_INT64
	case FieldUint64:
		return C.PSLOG_BENCH_QUILL_FIELD_UINT64
	case FieldFloat64:
		return C.PSLOG_BENCH_QUILL_FIELD_FLOAT64
	default:
		return C.PSLOG_BENCH_QUILL_FIELD_NULL
	}
}

func NewQuillData(entries []Entry) *QuillData {
	if !HasQuill() {
		return nil
	}
	if err := validateQuillEntries(entries); err != nil {
		return &QuillData{initErr: err}
	}
	if len(entries) == 0 {
		return &QuillData{}
	}

	totalFields := 0
	for _, entry := range entries {
		totalFields += len(entry.Fields)
	}

	data := &QuillData{
		entryCount: len(entries),
		fieldCount: totalFields,
	}
	data.entries = (*C.pslog_bench_quill_entry_spec)(C.calloc(C.size_t(len(entries)), C.size_t(C.sizeof_pslog_bench_quill_entry_spec)))
	if totalFields > 0 {
		data.fields = (*C.pslog_bench_quill_field_spec)(C.calloc(C.size_t(totalFields), C.size_t(C.sizeof_pslog_bench_quill_field_spec)))
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

		entrySlice[i].level = cQuillLevel(int(entry.Level))
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
			dst.kind = cQuillKind(field.Kind)
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
				dst.signed_value = C.longlong(field.Int64)
			case FieldUint64:
				dst.unsigned_value = C.ulonglong(field.Uint64)
			case FieldFloat64:
				dst.double_value = C.double(field.Float64)
			default:
			}
			fieldOffset++
		}
	}

	return data
}

func validateQuillEntries(entries []Entry) error {
	for entryIndex, entry := range entries {
		for fieldIndex, field := range entry.Fields {
			if field.Kind == FieldFloat64 && !isFiniteFloat64(field.Float64) {
				return fmt.Errorf(
					"entry %d field %d (%q) has non-finite float64 value",
					entryIndex,
					fieldIndex,
					field.Key,
				)
			}
		}
	}
	return nil
}

func isFiniteFloat64(value float64) bool {
	return !math.IsNaN(value) && !math.IsInf(value, 0)
}

func (d *QuillData) Close() {
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

func RunQuill(data *QuillData, iterations int) (QuillResult, error) {
	var result C.pslog_bench_quill_result

	if !HasQuill() {
		return QuillResult{}, fmt.Errorf("quill benchmark support is not enabled")
	}
	if data != nil && data.initErr != nil {
		return QuillResult{}, fmt.Errorf("quill benchmark data invalid: %w", data.initErr)
	}
	if data == nil || data.entries == nil || data.entryCount == 0 {
		return QuillResult{}, fmt.Errorf("quill benchmark data not initialized")
	}
	if rc := C.pslog_bench_quill_run(data.entries, C.size_t(data.entryCount), C.ulonglong(iterations), &result); rc != 0 {
		return QuillResult{}, fmt.Errorf("quill benchmark run failed")
	}
	return QuillResult{
		Elapsed:      time.Duration(result.elapsed_ns),
		BytesWritten: uint64(result.bytes),
		Ops:          uint64(result.ops),
	}, nil
}

func RenderQuillEntry(entry Entry) (string, error) {
	if !HasQuill() {
		return "", fmt.Errorf("quill benchmark support is not enabled")
	}
	data := NewQuillData([]Entry{entry})
	if data == nil {
		return "", fmt.Errorf("quill benchmark data not initialized")
	}
	if data.initErr != nil {
		return "", fmt.Errorf("quill benchmark data invalid: %w", data.initErr)
	}
	if data.entries == nil || data.entryCount == 0 {
		return "", fmt.Errorf("quill benchmark data not initialized")
	}
	defer data.Close()

	buf := make([]byte, 8192)
	var written C.size_t
	if rc := C.pslog_bench_quill_render(data.entries, (*C.char)(unsafe.Pointer(&buf[0])), C.size_t(len(buf)), &written); rc != 0 {
		return "", fmt.Errorf("quill benchmark render failed")
	}
	return string(buf[:int(written)]), nil
}
