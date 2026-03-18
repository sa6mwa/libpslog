package benchmark

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	pslog "pkt.systems/pslog"
)

const DefaultProductionLimit = 2048

type FieldKind uint8

const (
	FieldNull FieldKind = iota
	FieldString
	FieldBool
	FieldInt64
	FieldUint64
	FieldFloat64
)

type Field struct {
	Key     string
	Kind    FieldKind
	String  string
	Bool    bool
	Int64   int64
	Uint64  uint64
	Float64 float64
}

type Entry struct {
	Level   pslog.Level
	Message string
	Fields  []Field
	Keyvals []any
}

type KVFmtArg struct {
	Kind    FieldKind
	String  string
	Bool    bool
	Int64   int64
	Uint64  uint64
	Float64 float64
}

type KVFmtEntry struct {
	Level   pslog.Level
	Message string
	Format  string
	Args    []KVFmtArg
}

func (e Entry) LogGo(logger pslog.Logger) {
	logger.Log(e.Level, e.Message, e.Keyvals...)
}

func EntryToKVFmt(entry Entry) (KVFmtEntry, error) {
	if len(entry.Fields) == 0 {
		return KVFmtEntry{
			Level:   entry.Level,
			Message: entry.Message,
		}, nil
	}

	var format strings.Builder
	args := make([]KVFmtArg, len(entry.Fields))
	for i, field := range entry.Fields {
		verb, arg, err := fieldToKVFmtArg(field)
		if err != nil {
			return KVFmtEntry{}, err
		}
		if i > 0 {
			format.WriteByte(' ')
		}
		format.WriteString(field.Key)
		format.WriteByte('=')
		format.WriteString(verb)
		args[i] = arg
	}

	return KVFmtEntry{
		Level:   entry.Level,
		Message: entry.Message,
		Format:  format.String(),
		Args:    args,
	}, nil
}

func EntriesToKVFmt(entries []Entry) ([]KVFmtEntry, error) {
	out := make([]KVFmtEntry, len(entries))
	for i, entry := range entries {
		converted, err := EntryToKVFmt(entry)
		if err != nil {
			return nil, err
		}
		out[i] = converted
	}
	return out, nil
}

func (e KVFmtEntry) Signature() string {
	if len(e.Args) == 0 {
		return ""
	}
	var out strings.Builder
	for i, arg := range e.Args {
		if i > 0 {
			out.WriteByte(',')
		}
		out.WriteString(arg.signatureCode())
	}
	return out.String()
}

func LoadProductionEntries(limit int) ([]Entry, error) {
	if len(EmbeddedProductionDataset) == 0 {
		return nil, errors.New("embedded production dataset empty")
	}
	if limit <= 0 || limit > len(EmbeddedProductionDataset) {
		limit = len(EmbeddedProductionDataset)
	}

	entries := make([]Entry, 0, limit)
	for i := 0; i < limit; i++ {
		entry, err := ParseProductionLine(EmbeddedProductionDataset[i])
		if err != nil {
			continue
		}
		if entry.Level == pslog.Disabled {
			continue
		}
		entries = append(entries, entry)
	}
	if len(entries) == 0 {
		return nil, errors.New("no production log entries parsed")
	}
	return entries, nil
}

func ParseProductionLine(line string) (Entry, error) {
	decoder := json.NewDecoder(strings.NewReader(line))
	decoder.UseNumber()

	raw := make(map[string]any)
	if err := decoder.Decode(&raw); err != nil {
		return Entry{}, err
	}

	level := pslog.InfoLevel
	if lvl, ok := raw["lvl"].(string); ok {
		if parsed, ok := pslog.ParseLevel(lvl); ok {
			level = parsed
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

	fields := make([]Field, 0, len(keys))
	keyvals := make([]any, 0, len(keys)*2)
	for _, key := range keys {
		field := fieldFromValue(key, sanitizeJSONValue(raw[key]))
		fields = append(fields, field)
		keyvals = append(keyvals, key, field.Value())
	}

	return Entry{
		Level:   level,
		Message: message,
		Fields:  fields,
		Keyvals: keyvals,
	}, nil
}

func ProductionStaticFields(entries []Entry) ([]Field, map[string]struct{}) {
	if len(entries) == 0 {
		return nil, nil
	}

	constants := make(map[string]Field, len(entries[0].Fields))
	for _, field := range entries[0].Fields {
		constants[field.Key] = field
	}

	for _, entry := range entries[1:] {
		current := make(map[string]Field, len(entry.Fields))
		for _, field := range entry.Fields {
			current[field.Key] = field
		}
		for key, value := range constants {
			other, ok := current[key]
			if !ok || !reflect.DeepEqual(other, value) {
				delete(constants, key)
			}
		}
		if len(constants) == 0 {
			return nil, nil
		}
	}

	keys := make([]string, 0, len(constants))
	for key := range constants {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	fields := make([]Field, 0, len(keys))
	staticKeys := make(map[string]struct{}, len(keys))
	for _, key := range keys {
		staticKeys[key] = struct{}{}
		fields = append(fields, constants[key])
	}
	return fields, staticKeys
}

func ProductionEntriesWithout(entries []Entry, staticKeys map[string]struct{}) []Entry {
	if len(staticKeys) == 0 {
		return entries
	}

	filtered := make([]Entry, len(entries))
	for i, entry := range entries {
		nextFields := make([]Field, 0, len(entry.Fields))
		nextKeyvals := make([]any, 0, len(entry.Keyvals))
		for _, field := range entry.Fields {
			if _, ok := staticKeys[field.Key]; ok {
				continue
			}
			nextFields = append(nextFields, field)
			nextKeyvals = append(nextKeyvals, field.Key, field.Value())
		}
		filtered[i] = Entry{
			Level:   entry.Level,
			Message: entry.Message,
			Fields:  nextFields,
			Keyvals: nextKeyvals,
		}
	}
	return filtered
}

func FieldsToKeyvals(fields []Field) []any {
	if len(fields) == 0 {
		return nil
	}
	out := make([]any, 0, len(fields)*2)
	for _, field := range fields {
		out = append(out, field.Key, field.Value())
	}
	return out
}

func fieldToKVFmtArg(field Field) (string, KVFmtArg, error) {
	arg := KVFmtArg{Kind: field.Kind}
	switch field.Kind {
	case FieldString:
		arg.String = field.String
		return "%s", arg, nil
	case FieldBool:
		arg.Bool = field.Bool
		return "%b", arg, nil
	case FieldInt64:
		arg.Int64 = field.Int64
		return "%ld", arg, nil
	case FieldUint64:
		arg.Uint64 = field.Uint64
		return "%lu", arg, nil
	case FieldFloat64:
		arg.Float64 = field.Float64
		return "%f", arg, nil
	default:
		return "", KVFmtArg{}, fmt.Errorf("kvfmt does not support field %q kind %d", field.Key, field.Kind)
	}
}

func (a KVFmtArg) signatureCode() string {
	switch a.Kind {
	case FieldString:
		return "s"
	case FieldBool:
		return "b"
	case FieldInt64:
		return "i"
	case FieldUint64:
		return "u"
	case FieldFloat64:
		return "f"
	default:
		return "?"
	}
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

func fieldFromValue(key string, value any) Field {
	field := Field{Key: key}
	switch typed := value.(type) {
	case nil:
		field.Kind = FieldNull
	case string:
		field.Kind = FieldString
		field.String = typed
	case bool:
		field.Kind = FieldBool
		field.Bool = typed
	case int:
		field.Kind = FieldInt64
		field.Int64 = int64(typed)
	case int64:
		field.Kind = FieldInt64
		field.Int64 = typed
	case uint64:
		field.Kind = FieldUint64
		field.Uint64 = typed
	case float64:
		field.Kind = FieldFloat64
		field.Float64 = typed
	case []byte:
		field.Kind = FieldString
		field.String = string(typed)
	default:
		field.Kind = FieldString
		field.String = marshalCompactJSON(typed)
	}
	return field
}

func (f Field) Value() any {
	switch f.Kind {
	case FieldNull:
		return nil
	case FieldString:
		return f.String
	case FieldBool:
		return f.Bool
	case FieldInt64:
		return f.Int64
	case FieldUint64:
		return f.Uint64
	case FieldFloat64:
		return f.Float64
	default:
		return nil
	}
}

func CompactJSON(line string) string {
	var out bytes.Buffer
	if err := json.Compact(&out, []byte(line)); err != nil {
		return line
	}
	return out.String()
}
