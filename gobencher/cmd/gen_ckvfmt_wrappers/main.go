package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

const (
	outputPath  = "cpslog_kvfmt_generated.go"
	datasetPath = "production_data_generated.go"
)

type fieldKind uint8

const (
	fieldNull fieldKind = iota
	fieldString
	fieldBool
	fieldInt64
	fieldUint64
	fieldFloat64
)

type field struct {
	key    string
	kind   fieldKind
	string string
	bool   bool
	int64  int64
	uint64 uint64
	float  float64
}

type entry struct {
	level   string
	message string
	fields  []field
}

func main() {
	signatures, _, err := collectSignatures()
	if err != nil {
		fail(err)
	}
	content, err := render(signatures)
	if err != nil {
		fail(err)
	}
	if err := os.WriteFile(outputPath, content, 0644); err != nil {
		fail(err)
	}
}

func collectSignatures() ([]string, []entry, error) {
	signatureSet := map[string]struct{}{
		"i,s,b,f": {},
	}
	lines, err := loadEmbeddedDataset(datasetPath)
	if err != nil {
		return nil, nil, err
	}
	entries, err := parseEntries(lines)
	if err != nil {
		return nil, nil, err
	}
	staticKeys := productionStaticKeys(entries)
	for _, current := range entries {
		signatureSet[current.dynamicSignature(nil)] = struct{}{}
	}
	for _, current := range entries {
		signature := current.dynamicSignature(staticKeys)
		signatureSet[signature] = struct{}{}
	}

	signatures := make([]string, 0, len(signatureSet))
	for signature := range signatureSet {
		signatures = append(signatures, signature)
	}
	sort.Strings(signatures)
	return signatures, entries, nil
}

func loadEmbeddedDataset(path string) ([]string, error) {
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, path, nil, 0)
	if err != nil {
		return nil, err
	}
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.VAR {
			continue
		}
		for _, spec := range genDecl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok || len(valueSpec.Names) != 1 {
				continue
			}
			if valueSpec.Names[0].Name != "EmbeddedProductionDataset" {
				continue
			}
			if len(valueSpec.Values) != 1 {
				return nil, fmt.Errorf("unexpected production dataset form")
			}
			composite, ok := valueSpec.Values[0].(*ast.CompositeLit)
			if !ok {
				return nil, fmt.Errorf("production dataset is not a composite literal")
			}
			lines := make([]string, 0, len(composite.Elts))
			for _, elt := range composite.Elts {
				lit, ok := elt.(*ast.BasicLit)
				if !ok || lit.Kind != token.STRING {
					return nil, fmt.Errorf("unexpected dataset literal element")
				}
				value, err := strconv.Unquote(lit.Value)
				if err != nil {
					return nil, err
				}
				lines = append(lines, value)
			}
			return lines, nil
		}
	}
	return nil, fmt.Errorf("production dataset not found in %s", filepath.Clean(path))
}

func parseEntries(lines []string) ([]entry, error) {
	entries := make([]entry, 0, len(lines))
	for _, line := range lines {
		current, err := parseLine(line)
		if err != nil {
			continue
		}
		entries = append(entries, current)
	}
	if len(entries) == 0 {
		return nil, fmt.Errorf("no production entries parsed")
	}
	return entries, nil
}

func parseLine(line string) (entry, error) {
	decoder := json.NewDecoder(strings.NewReader(line))
	decoder.UseNumber()

	raw := make(map[string]any)
	if err := decoder.Decode(&raw); err != nil {
		return entry{}, err
	}

	level := "info"
	if lvl, ok := raw["lvl"].(string); ok {
		switch strings.ToLower(lvl) {
		case "trace":
			level = "trace"
		case "debug":
			level = "debug"
		case "info":
			level = "info"
		case "warn", "warning":
			level = "warn"
		case "error":
			level = "error"
		case "fatal":
			level = "fatal"
		case "panic":
			level = "panic"
		case "none", "nolevel":
			level = "nolevel"
		}
	}
	delete(raw, "lvl")

	message := ""
	if msg, ok := raw["msg"].(string); ok {
		message = msg
	}
	delete(raw, "msg")
	delete(raw, "ts")
	delete(raw, "ts_iso")
	delete(raw, "time")

	keys := make([]string, 0, len(raw))
	for key := range raw {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	fields := make([]field, 0, len(keys))
	for _, key := range keys {
		fields = append(fields, fieldFromValue(key, sanitizeJSONValue(raw[key])))
	}
	return entry{level: level, message: message, fields: fields}, nil
}

func sanitizeJSONValue(v any) any {
	switch value := v.(type) {
	case json.Number:
		text := value.String()
		if !strings.ContainsAny(text, ".eE") {
			if i, err := value.Int64(); err == nil {
				return i
			}
			if u, err := strconv.ParseUint(text, 10, 64); err == nil {
				return u
			}
		}
		if f, err := value.Float64(); err == nil {
			return f
		}
		return text
	case map[string]any:
		normalized := make(map[string]any, len(value))
		for key, inner := range value {
			normalized[key] = sanitizeJSONValue(inner)
		}
		return marshalCompactJSON(normalized)
	case []any:
		normalized := make([]any, len(value))
		for i, inner := range value {
			normalized[i] = sanitizeJSONValue(inner)
		}
		return marshalCompactJSON(normalized)
	default:
		return value
	}
}

func marshalCompactJSON(v any) string {
	data, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(data)
}

func fieldFromValue(key string, value any) field {
	out := field{key: key}
	switch typed := value.(type) {
	case nil:
		out.kind = fieldNull
	case string:
		out.kind = fieldString
		out.string = typed
	case bool:
		out.kind = fieldBool
		out.bool = typed
	case int:
		out.kind = fieldInt64
		out.int64 = int64(typed)
	case int64:
		out.kind = fieldInt64
		out.int64 = typed
	case uint64:
		out.kind = fieldUint64
		out.uint64 = typed
	case float64:
		out.kind = fieldFloat64
		out.float = typed
	case []byte:
		out.kind = fieldString
		out.string = string(typed)
	default:
		out.kind = fieldString
		out.string = marshalCompactJSON(typed)
	}
	return out
}

func productionStaticKeys(entries []entry) map[string]struct{} {
	if len(entries) == 0 {
		return nil
	}
	constants := make(map[string]field, len(entries[0].fields))
	for _, current := range entries[0].fields {
		constants[current.key] = current
	}
	for _, current := range entries[1:] {
		currentFields := make(map[string]field, len(current.fields))
		for _, item := range current.fields {
			currentFields[item.key] = item
		}
		for key, value := range constants {
			other, ok := currentFields[key]
			if !ok || !reflect.DeepEqual(other, value) {
				delete(constants, key)
			}
		}
		if len(constants) == 0 {
			return nil
		}
	}
	staticKeys := make(map[string]struct{}, len(constants))
	for key := range constants {
		staticKeys[key] = struct{}{}
	}
	return staticKeys
}

func (e entry) dynamicSignature(staticKeys map[string]struct{}) string {
	var out strings.Builder
	first := true
	for _, current := range e.fields {
		if _, ok := staticKeys[current.key]; ok {
			continue
		}
		if !first {
			out.WriteByte(',')
		}
		first = false
		out.WriteString(current.signatureCode())
	}
	return out.String()
}

func (f field) signatureCode() string {
	switch f.kind {
	case fieldString:
		return "s"
	case fieldBool:
		return "b"
	case fieldInt64:
		return "i"
	case fieldUint64:
		return "u"
	case fieldFloat64:
		return "f"
	default:
		return "?"
	}
}

func render(signatures []string) ([]byte, error) {
	var out bytes.Buffer

	out.WriteString("package benchmark\n\n")
	out.WriteString("/*\n")
	out.WriteString("#include <stddef.h>\n")
	out.WriteString("#include <time.h>\n\n")
	out.WriteString("#include \"pslog.h\"\n\n")
	out.WriteString("typedef struct gobench_sink {\n")
	out.WriteString("\tunsigned long long bytes;\n")
	out.WriteString("\tunsigned long long writes;\n")
	out.WriteString("\tchar *capture;\n")
	out.WriteString("\tsize_t capture_cap;\n")
	out.WriteString("\tsize_t capture_len;\n")
	out.WriteString("} gobench_sink;\n\n")
	out.WriteString("typedef struct gobench_result {\n")
	out.WriteString("\tunsigned long long elapsed_ns;\n")
	out.WriteString("\tunsigned long long bytes;\n")
	out.WriteString("\tunsigned long long writes;\n")
	out.WriteString("\tunsigned long long ops;\n")
	out.WriteString("} gobench_result;\n\n")
	out.WriteString("typedef struct gobench_logger {\n")
	out.WriteString("\tgobench_sink *sink;\n")
	out.WriteString("\tpslog_logger *logger;\n")
	out.WriteString("\tint owns_sink;\n")
	out.WriteString("} gobench_logger;\n\n")
	out.WriteString("typedef struct gobench_kvfmt_arg {\n")
	out.WriteString("\tint kind;\n")
	out.WriteString("\tconst char *string_value;\n")
	out.WriteString("\tint bool_value;\n")
	out.WriteString("\tlong signed_value;\n")
	out.WriteString("\tunsigned long unsigned_value;\n")
	out.WriteString("\tdouble double_value;\n")
	out.WriteString("} gobench_kvfmt_arg;\n\n")
	out.WriteString("typedef struct gobench_kvfmt_entry {\n")
	out.WriteString("\tpslog_level level;\n")
	out.WriteString("\tconst char *msg;\n")
	out.WriteString("\tconst char *kvfmt;\n")
	out.WriteString("\tint signature_id;\n")
	out.WriteString("\tvoid (*call)(pslog_logger *logger, const struct gobench_kvfmt_entry *entry);\n")
	out.WriteString("\tconst gobench_kvfmt_arg *args;\n")
	out.WriteString("\tsize_t arg_count;\n")
	out.WriteString("} gobench_kvfmt_entry;\n\n")
	out.WriteString(renderCHelpers())

	for index, signature := range signatures {
		out.WriteString(renderWrapper(index, signature))
	}
	out.WriteString("int gobench_kvfmt_entry_set_signature(gobench_kvfmt_entry *entry, int signature_id) {\n")
	out.WriteString("\tif (entry == NULL) {\n")
	out.WriteString("\t\treturn -1;\n")
	out.WriteString("\t}\n")
	out.WriteString("\tentry->signature_id = signature_id;\n")
	out.WriteString("\tswitch (entry->level) {\n")
	out.WriteString("\tcase PSLOG_LEVEL_TRACE:\n")
	out.WriteString("\t\tswitch (signature_id) {\n")
	for index := range signatures {
		out.WriteString(fmt.Sprintf("\t\tcase %d:\n", index))
		out.WriteString(fmt.Sprintf("\t\t\tentry->call = gobench_logger_tracef_sig_%d;\n", index))
		out.WriteString("\t\t\treturn 0;\n")
	}
	out.WriteString("\t\tdefault:\n")
	out.WriteString("\t\t\tentry->call = NULL;\n")
	out.WriteString("\t\t\treturn -1;\n")
	out.WriteString("\t\t}\n")
	out.WriteString("\tcase PSLOG_LEVEL_DEBUG:\n")
	out.WriteString("\t\tswitch (signature_id) {\n")
	for index := range signatures {
		out.WriteString(fmt.Sprintf("\t\tcase %d:\n", index))
		out.WriteString(fmt.Sprintf("\t\t\tentry->call = gobench_logger_debugf_sig_%d;\n", index))
		out.WriteString("\t\t\treturn 0;\n")
	}
	out.WriteString("\t\tdefault:\n")
	out.WriteString("\t\t\tentry->call = NULL;\n")
	out.WriteString("\t\t\treturn -1;\n")
	out.WriteString("\t\t}\n")
	out.WriteString("\tcase PSLOG_LEVEL_INFO:\n")
	out.WriteString("\t\tswitch (signature_id) {\n")
	for index := range signatures {
		out.WriteString(fmt.Sprintf("\t\tcase %d:\n", index))
		out.WriteString(fmt.Sprintf("\t\t\tentry->call = gobench_logger_infof_sig_%d;\n", index))
		out.WriteString("\t\t\treturn 0;\n")
	}
	out.WriteString("\t\tdefault:\n")
	out.WriteString("\t\t\tentry->call = NULL;\n")
	out.WriteString("\t\t\treturn -1;\n")
	out.WriteString("\t\t}\n")
	out.WriteString("\tcase PSLOG_LEVEL_WARN:\n")
	out.WriteString("\t\tswitch (signature_id) {\n")
	for index := range signatures {
		out.WriteString(fmt.Sprintf("\t\tcase %d:\n", index))
		out.WriteString(fmt.Sprintf("\t\t\tentry->call = gobench_logger_warnf_sig_%d;\n", index))
		out.WriteString("\t\t\treturn 0;\n")
	}
	out.WriteString("\t\tdefault:\n")
	out.WriteString("\t\t\tentry->call = NULL;\n")
	out.WriteString("\t\t\treturn -1;\n")
	out.WriteString("\t\t}\n")
	out.WriteString("\tcase PSLOG_LEVEL_ERROR:\n")
	out.WriteString("\t\tswitch (signature_id) {\n")
	for index := range signatures {
		out.WriteString(fmt.Sprintf("\t\tcase %d:\n", index))
		out.WriteString(fmt.Sprintf("\t\t\tentry->call = gobench_logger_errorf_sig_%d;\n", index))
		out.WriteString("\t\t\treturn 0;\n")
	}
	out.WriteString("\t\tdefault:\n")
	out.WriteString("\t\t\tentry->call = NULL;\n")
	out.WriteString("\t\t\treturn -1;\n")
	out.WriteString("\t\t}\n")
	out.WriteString("\tdefault:\n")
	out.WriteString("\t\tswitch (signature_id) {\n")
	for index := range signatures {
		out.WriteString(fmt.Sprintf("\t\tcase %d:\n", index))
		out.WriteString(fmt.Sprintf("\t\t\tentry->call = gobench_logger_logf_sig_%d;\n", index))
		out.WriteString("\t\t\treturn 0;\n")
	}
	out.WriteString("\t\tdefault:\n")
	out.WriteString("\t\t\tentry->call = NULL;\n")
	out.WriteString("\t\t\treturn -1;\n")
	out.WriteString("\t\t}\n")
	out.WriteString("\t}\n")
	out.WriteString("}\n\n")

	out.WriteString("int gobench_logger_log_kvfmt(void *logger_ptr,\n")
	out.WriteString("\t\t\t\t const gobench_kvfmt_entry *entry) {\n")
	out.WriteString("\tgobench_logger *wrapper;\n\n")
	out.WriteString("\twrapper = (gobench_logger *)logger_ptr;\n")
	out.WriteString("\tif (wrapper == NULL || wrapper->logger == NULL || entry == NULL || entry->call == NULL) {\n")
	out.WriteString("\t\treturn -1;\n")
	out.WriteString("\t}\n")
	out.WriteString("\tentry->call(wrapper->logger, entry);\n")
	out.WriteString("\treturn 0;\n")
	out.WriteString("}\n\n")

	out.WriteString("int gobench_logger_run_kvfmt(void *logger_ptr,\n")
	out.WriteString("\t\t\t\t const gobench_kvfmt_entry *entries,\n")
	out.WriteString("\t\t\t\t size_t entry_count,\n")
	out.WriteString("\t\t\t\t const pslog_field *static_fields,\n")
	out.WriteString("\t\t\t\t size_t static_count,\n")
	out.WriteString("\t\t\t\t size_t iterations,\n")
	out.WriteString("\t\t\t\t gobench_result *out) {\n")
	out.WriteString("\tunsigned long long start_ns;\n")
	out.WriteString("\tunsigned long long end_ns;\n")
	out.WriteString("\tgobench_logger *wrapper;\n")
	out.WriteString("\tpslog_logger *active_logger;\n")
	out.WriteString("\tpslog_logger *derived;\n")
	out.WriteString("\tsize_t i;\n")
	out.WriteString("\tsize_t index;\n\n")
	out.WriteString("\twrapper = (gobench_logger *)logger_ptr;\n")
	out.WriteString("\tif (wrapper == NULL || wrapper->logger == NULL || entries == NULL || entry_count == 0u || out == NULL) {\n")
	out.WriteString("\t\treturn -1;\n")
	out.WriteString("\t}\n\n")
	out.WriteString("\tactive_logger = wrapper->logger;\n")
	out.WriteString("\tderived = NULL;\n")
	out.WriteString("\tif (static_count > 0u) {\n")
	out.WriteString("\t\tderived = wrapper->logger->with(wrapper->logger, static_fields, static_count);\n")
	out.WriteString("\t\tif (derived == NULL) {\n")
	out.WriteString("\t\t\treturn -1;\n")
	out.WriteString("\t\t}\n")
	out.WriteString("\t\tactive_logger = derived;\n")
	out.WriteString("\t}\n\n")
	out.WriteString("\tgobench_kvfmt_sink_reset(wrapper->sink);\n")
	out.WriteString("\tindex = 0u;\n")
	out.WriteString("\tstart_ns = gobench_kvfmt_monotonic_ns();\n")
	out.WriteString("\tfor (i = 0u; i < iterations; ++i) {\n")
	out.WriteString("\t\tconst gobench_kvfmt_entry *entry = entries + index;\n")
	out.WriteString("\t\tif (entry->call == NULL) {\n")
	out.WriteString("\t\t\tif (derived != NULL) {\n")
	out.WriteString("\t\t\t\tderived->destroy(derived);\n")
	out.WriteString("\t\t\t}\n")
	out.WriteString("\t\t\treturn -1;\n")
	out.WriteString("\t\t}\n")
	out.WriteString("\t\tentry->call(active_logger, entry);\n")
	out.WriteString("\t\t++index;\n")
	out.WriteString("\t\tif (index == entry_count) {\n")
	out.WriteString("\t\t\tindex = 0u;\n")
	out.WriteString("\t\t}\n")
	out.WriteString("\t}\n")
	out.WriteString("\tend_ns = gobench_kvfmt_monotonic_ns();\n\n")
	out.WriteString("\tout->elapsed_ns = end_ns - start_ns;\n")
	out.WriteString("\tout->bytes = wrapper->sink != NULL ? wrapper->sink->bytes : 0ull;\n")
	out.WriteString("\tout->writes = wrapper->sink != NULL ? wrapper->sink->writes : 0ull;\n")
	out.WriteString("\tout->ops = (unsigned long long)iterations;\n\n")
	out.WriteString("\tif (derived != NULL) {\n")
	out.WriteString("\t\tderived->destroy(derived);\n")
	out.WriteString("\t}\n")
	out.WriteString("\treturn 0;\n")
	out.WriteString("}\n")
	out.WriteString("*/\n")
	out.WriteString("import \"C\"\n\n")
	out.WriteString("func init() {\n")
	for index, signature := range signatures {
		out.WriteString(fmt.Sprintf("\tckvfmtSignatureIDs[%q] = %d\n", signature, index))
	}
	out.WriteString("}\n")

	return format.Source(out.Bytes())
}

func renderCHelpers() string {
	return `static unsigned long long gobench_kvfmt_monotonic_ns(void) {
#if defined(CLOCK_MONOTONIC)
	struct timespec ts;
	if (clock_gettime(CLOCK_MONOTONIC, &ts) == 0) {
		return ((unsigned long long)ts.tv_sec * 1000000000ull) + (unsigned long long)ts.tv_nsec;
	}
#endif
	return ((unsigned long long)clock() * 1000000000ull) / (unsigned long long)CLOCKS_PER_SEC;
}

static void gobench_kvfmt_sink_reset(gobench_sink *sink) {
	if (sink == NULL) {
		return;
	}
	sink->bytes = 0ull;
	sink->writes = 0ull;
	sink->capture_len = 0u;
	if (sink->capture != NULL && sink->capture_cap > 0u) {
		sink->capture[0] = '\0';
	}
}

`
}

func renderWrapper(index int, signature string) string {
	parts := []string{}
	if signature != "" {
		parts = strings.Split(signature, ",")
	}

	var out strings.Builder
	renderWrapperFn(&out, "tracef", index, parts)
	renderWrapperFn(&out, "debugf", index, parts)
	renderWrapperFn(&out, "infof", index, parts)
	renderWrapperFn(&out, "warnf", index, parts)
	renderWrapperFn(&out, "errorf", index, parts)
	renderWrapperFn(&out, "logf", index, parts)
	return out.String()
}

func renderWrapperFn(out *strings.Builder, levelFn string, index int, parts []string) {
	out.WriteString(fmt.Sprintf("static void gobench_logger_%s_sig_%d(pslog_logger *logger,\n", levelFn, index))
	out.WriteString("\t\t\t\t      const gobench_kvfmt_entry *entry) {\n")
	out.WriteString("\tconst char *msg = entry->msg;\n")
	out.WriteString("\tconst char *kvfmt = entry->kvfmt;\n")
	out.WriteString("\tconst gobench_kvfmt_arg *args = entry->args;\n")
	out.WriteString("\tlogger->")
	out.WriteString(levelFn)
	out.WriteString("(logger")
	if levelFn == "logf" {
		out.WriteString(", entry->level")
	}
	out.WriteString(", msg, kvfmt")
	for i, part := range parts {
		out.WriteString(", ")
		out.WriteString(wrapperArgExpr(part, i))
	}
	out.WriteString(");\n")
	out.WriteString("}\n\n")
}


func wrapperArgExpr(code string, index int) string {
	switch code {
	case "s":
		return fmt.Sprintf("args[%d].string_value", index)
	case "b":
		return fmt.Sprintf("args[%d].bool_value", index)
	case "i":
		return fmt.Sprintf("args[%d].signed_value", index)
	case "u":
		return fmt.Sprintf("args[%d].unsigned_value", index)
	case "f":
		return fmt.Sprintf("args[%d].double_value", index)
	default:
		panic("unsupported signature code: " + code)
	}
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
