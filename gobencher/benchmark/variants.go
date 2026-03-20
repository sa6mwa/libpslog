package benchmark

import (
	"io"

	pslog "pkt.systems/pslog"
)

type Variant struct {
	Name  string
	Mode  pslog.Mode
	Color bool
}

func Variants() []Variant {
	return []Variant{
		{Name: "json", Mode: pslog.ModeStructured, Color: false},
		{Name: "jsoncolor", Mode: pslog.ModeStructured, Color: true},
		{Name: "console", Mode: pslog.ModeConsole, Color: false},
		{Name: "consolecolor", Mode: pslog.ModeConsole, Color: true},
	}
}

func ComparisonNames() []string {
	names := []string{
		"jsonGo",
		"jsonC",
		"coljsonGo",
		"coljsonC",
		"consoleGo",
		"consoleC",
		"colconGo",
		"colconC",
	}
	return append(names, libloggerComparisonNames()...)
}

func ComparisonNamesWithQuill() []string {
	names := ComparisonNames()
	return append(names, quillComparisonNames()...)
}

func newGoLoggerWithTimestamp(w io.Writer, variant Variant, timestamps bool) pslog.Logger {
	opts := pslog.Options{
		Mode:       variant.Mode,
		MinLevel:   pslog.TraceLevel,
		ForceColor: variant.Color,
		NoColor:    !variant.Color,
	}
	if !timestamps {
		opts.DisableTimestamp = true
	}
	return pslog.NewWithOptions(nil, w, opts)
}

func NewGoLogger(w io.Writer, variant Variant) pslog.Logger {
	return newGoLoggerWithTimestamp(w, variant, true)
}
