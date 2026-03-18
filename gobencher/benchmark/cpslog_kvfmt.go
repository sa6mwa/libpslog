package benchmark

//go:generate go run ../cmd/gen_ckvfmt_wrappers

/*
#cgo CFLAGS: -O3 -DNDEBUG -I${SRCDIR}/../../include -I${SRCDIR}/../../build/release/generated/include
#cgo LDFLAGS: ${SRCDIR}/../../build/release/libpslog.a

#include <stddef.h>
#include <stdlib.h>

#include "pslog.h"

typedef struct gobench_result {
	unsigned long long elapsed_ns;
	unsigned long long bytes;
	unsigned long long writes;
	unsigned long long ops;
} gobench_result;

typedef struct gobench_kvfmt_arg {
	int kind;
	const char *string_value;
	int bool_value;
	long signed_value;
	unsigned long unsigned_value;
	double double_value;
} gobench_kvfmt_arg;

typedef struct gobench_kvfmt_entry {
	pslog_level level;
	const char *msg;
	const char *kvfmt;
	int signature_id;
	void (*call)(pslog_logger *logger, const struct gobench_kvfmt_entry *entry);
	const gobench_kvfmt_arg *args;
	size_t arg_count;
} gobench_kvfmt_entry;

int gobench_kvfmt_entry_set_signature(gobench_kvfmt_entry *entry, int signature_id);
int gobench_logger_log_kvfmt(void *logger_ptr, const gobench_kvfmt_entry *entry);
int gobench_logger_run_kvfmt(void *logger_ptr,
                             const gobench_kvfmt_entry *entries,
                             size_t entry_count,
                             const pslog_field *static_fields,
                             size_t static_count,
                             size_t iterations,
                             gobench_result *out);
*/
import "C"

import (
	"fmt"
	"time"
	"unsafe"

	pslog "pkt.systems/pslog"
)

const (
	cKVFmtArgString = 1
	cKVFmtArgBool   = 2
	cKVFmtArgSigned = 3
	cKVFmtArgUint   = 4
	cKVFmtArgDouble = 5
)

type CKVFmtEntry struct {
	ptr *C.gobench_kvfmt_entry
}

type CKVFmtData struct {
	Entries     []CKVFmtEntry
	entriesPtr  *C.gobench_kvfmt_entry
	argsPtr     *C.gobench_kvfmt_arg
	pool        map[string]*C.char
	allocations []unsafe.Pointer
}

var ckvfmtSignatureIDs = map[string]int{}

func NewCKVFmtData(entries []Entry) (*CKVFmtData, error) {
	kvfmtEntries, err := EntriesToKVFmt(entries)
	if err != nil {
		return nil, err
	}
	return NewCKVFmtDataFromEntries(kvfmtEntries)
}

func NewCKVFmtDataFromEntries(entries []KVFmtEntry) (*CKVFmtData, error) {
	data := &CKVFmtData{
		Entries: make([]CKVFmtEntry, len(entries)),
		pool:    make(map[string]*C.char),
	}
	if len(entries) == 0 {
		return data, nil
	}

	totalArgs := 0
	for _, entry := range entries {
		if _, ok := ckvfmtSignatureID(entry.Signature()); !ok {
			data.Close()
			return nil, fmt.Errorf("unsupported kvfmt signature %q", entry.Signature())
		}
		totalArgs += len(entry.Args)
	}

	data.entriesPtr = (*C.gobench_kvfmt_entry)(data.alloc(C.size_t(len(entries)) * C.size_t(C.sizeof_gobench_kvfmt_entry)))

	var argBase unsafe.Pointer
	if totalArgs > 0 {
		data.argsPtr = (*C.gobench_kvfmt_arg)(data.alloc(C.size_t(totalArgs) * C.size_t(C.sizeof_gobench_kvfmt_arg)))
		argBase = unsafe.Pointer(data.argsPtr)
	}

	offset := 0
	for i, entry := range entries {
		var argsPtr *C.gobench_kvfmt_arg
		signatureID, _ := ckvfmtSignatureID(entry.Signature())
		if len(entry.Args) > 0 {
			argsPtr = (*C.gobench_kvfmt_arg)(unsafe.Add(argBase, uintptr(offset)*unsafe.Sizeof(*data.argsPtr)))
			data.fillArgs(argsPtr, entry.Args)
			offset += len(entry.Args)
		}

		data.Entries[i].ptr = (*C.gobench_kvfmt_entry)(unsafe.Add(unsafe.Pointer(data.entriesPtr), uintptr(i)*unsafe.Sizeof(*data.entriesPtr)))
		data.Entries[i].ptr.level = cLevelFromGo(entry.Level)
		data.Entries[i].ptr.msg = data.cstring(entry.Message)
		data.Entries[i].ptr.kvfmt = data.cstring(entry.Format)
		if C.gobench_kvfmt_entry_set_signature(data.Entries[i].ptr, C.int(signatureID)) != 0 {
			data.Close()
			return nil, fmt.Errorf("set kvfmt signature %q", entry.Signature())
		}
		data.Entries[i].ptr.args = argsPtr
		data.Entries[i].ptr.arg_count = C.size_t(len(entry.Args))
	}

	return data, nil
}

func (d *CKVFmtData) Close() {
	if d == nil {
		return
	}
	for _, value := range d.pool {
		C.free(unsafe.Pointer(value))
	}
	for _, ptr := range d.allocations {
		C.free(ptr)
	}
	d.pool = nil
	d.allocations = nil
	d.Entries = nil
	d.entriesPtr = nil
	d.argsPtr = nil
}

func (l *CLogger) RunKVFmt(data *CKVFmtData, staticFields CPreparedFields, iterations int) (CRunResult, error) {
	var result C.gobench_result

	if l == nil || l.ptr == nil {
		return CRunResult{}, fmt.Errorf("nil C logger")
	}
	if data == nil || data.entriesPtr == nil || len(data.Entries) == 0 {
		return CRunResult{}, fmt.Errorf("empty kvfmt data")
	}
	if iterations < 0 {
		iterations = 0
	}

	if C.gobench_logger_run_kvfmt(
		unsafe.Pointer(l.ptr),
		data.entriesPtr,
		C.size_t(len(data.Entries)),
		(*C.pslog_field)(unsafe.Pointer(staticFields.ptr)),
		staticFields.count,
		C.size_t(iterations),
		&result,
	) != 0 {
		return CRunResult{}, fmt.Errorf("run kvfmt benchmark")
	}

	return CRunResult{
		Elapsed:      time.Duration(result.elapsed_ns),
		BytesWritten: uint64(result.bytes),
		Writes:       uint64(result.writes),
		Ops:          uint64(result.ops),
	}, nil
}

func (l *CLogger) LogKVFmt(entry CKVFmtEntry) error {
	if l == nil || l.ptr == nil || entry.ptr == nil {
		return fmt.Errorf("nil C logger")
	}
	if C.gobench_logger_log_kvfmt(unsafe.Pointer(l.ptr), entry.ptr) != 0 {
		return fmt.Errorf("log kvfmt entry")
	}
	return nil
}

func (d *CKVFmtData) alloc(size C.size_t) unsafe.Pointer {
	var ptr unsafe.Pointer
	if size == 0 {
		return nil
	}
	ptr = C.malloc(size)
	d.allocations = append(d.allocations, ptr)
	return ptr
}

func (d *CKVFmtData) cstring(text string) *C.char {
	if existing, ok := d.pool[text]; ok {
		return existing
	}
	value := C.CString(text)
	d.pool[text] = value
	return value
}

func (d *CKVFmtData) fillArgs(dst *C.gobench_kvfmt_arg, args []KVFmtArg) {
	for i, arg := range args {
		dstArg := (*C.gobench_kvfmt_arg)(unsafe.Add(unsafe.Pointer(dst), uintptr(i)*unsafe.Sizeof(*dst)))
		dstArg.string_value = nil
		dstArg.bool_value = 0
		dstArg.signed_value = 0
		dstArg.unsigned_value = 0
		dstArg.double_value = 0
		switch arg.Kind {
		case FieldString:
			dstArg.kind = C.int(cKVFmtArgString)
			dstArg.string_value = d.cstring(arg.String)
		case FieldBool:
			dstArg.kind = C.int(cKVFmtArgBool)
			if arg.Bool {
				dstArg.bool_value = 1
			}
		case FieldInt64:
			dstArg.kind = C.int(cKVFmtArgSigned)
			dstArg.signed_value = C.long(arg.Int64)
		case FieldUint64:
			dstArg.kind = C.int(cKVFmtArgUint)
			dstArg.unsigned_value = C.ulong(arg.Uint64)
		case FieldFloat64:
			dstArg.kind = C.int(cKVFmtArgDouble)
			dstArg.double_value = C.double(arg.Float64)
		default:
			dstArg.kind = 0
		}
	}
}

func ckvfmtSignatureID(signature string) (int, bool) {
	id, ok := ckvfmtSignatureIDs[signature]
	return id, ok
}

func cLevelFromGo(level pslog.Level) C.pslog_level {
	switch level {
	case pslog.TraceLevel:
		return C.PSLOG_LEVEL_TRACE
	case pslog.DebugLevel:
		return C.PSLOG_LEVEL_DEBUG
	case pslog.InfoLevel:
		return C.PSLOG_LEVEL_INFO
	case pslog.WarnLevel:
		return C.PSLOG_LEVEL_WARN
	case pslog.ErrorLevel:
		return C.PSLOG_LEVEL_ERROR
	case pslog.FatalLevel:
		return C.PSLOG_LEVEL_FATAL
	case pslog.PanicLevel:
		return C.PSLOG_LEVEL_PANIC
	case pslog.NoLevel:
		return C.PSLOG_LEVEL_NOLEVEL
	default:
		return C.PSLOG_LEVEL_INFO
	}
}
