package amapretty

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"time"
)

var timeNow = time.Now
var runtimeCaller = runtime.Caller
var output = os.Stdout
var lookupEnv = os.LookupEnv
var isTerminal = func(w io.Writer) bool {
	f, ok := w.(*os.File)
	if !ok {
		return false
	}

	stat, err := f.Stat()
	return err == nil && stat.Mode()&os.ModeCharDevice != 0
}

const (
	prefix = "amapretty"
)

// Print writes values to stdout as indented JSON with a timestamp and caller reference.
func Print(args ...interface{}) {
	_, _ = fprintWithCallerSkip(output, 3, args...)
}

// Printf formats according to a format specifier and writes the result through Print.
func Printf(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	_, _ = fprintWithCallerSkip(output, 3, s)
}

// Fprint writes values to w as indented JSON with a timestamp and caller reference.
func Fprint(w io.Writer, args ...interface{}) (int, error) {
	return fprintWithCallerSkip(w, 3, args...)
}

// Fprintf formats according to a format specifier and writes the result to w.
func Fprintf(w io.Writer, format string, args ...interface{}) (int, error) {
	s := fmt.Sprintf(format, args...)
	return fprintWithCallerSkip(w, 3, s)
}

// Sprint returns values formatted as indented JSON with a timestamp and caller reference.
func Sprint(args ...interface{}) string {
	return sprintWithCallerSkip(2, args...)
}

// Sprintf formats according to a format specifier and returns the formatted output.
func Sprintf(format string, args ...interface{}) string {
	s := fmt.Sprintf(format, args...)
	return sprintWithCallerSkip(2, s)
}

func fprintWithCallerSkip(w io.Writer, skip int, args ...interface{}) (int, error) {
	return fmt.Fprint(w, formatWithCallerSkip(skip, shouldUseColor(w), args...))
}

func sprintWithCallerSkip(skip int, args ...interface{}) string {
	return formatWithCallerSkip(skip, shouldUseColor(output), args...)
}

func formatWithCallerSkip(skip int, useColor bool, args ...interface{}) string {
	timeNow := timeNow().Format(time.RFC3339)
	fmtTimeNow := colorize(useColor, "1;34", timeNow)
	fmtPrefix := colorize(useColor, "1;32", prefix)

	_, fileName, fileLine, ok := runtimeCaller(skip)
	caller := ""
	if ok {
		caller = fmt.Sprintf("%s:%d", fileName, fileLine)
		caller = colorize(useColor, "1;36", caller)
	}

	s, err := json.MarshalIndent(args, "", "\t")
	if err != nil {
		s, _ = json.MarshalIndent(struct {
			Error string `json:"error"`
			Args  string `json:"args"`
		}{
			Error: err.Error(),
			Args:  fmt.Sprintf("%#v", args),
		}, "", "\t")
	}
	return fmt.Sprintf("[%s] %s %s -- %s\n", fmtPrefix, fmtTimeNow, caller, string(s))
}

func shouldUseColor(w io.Writer) bool {
	if _, ok := lookupEnv("NO_COLOR"); ok {
		return false
	}
	return isTerminal(w)
}

func colorize(enabled bool, code string, s string) string {
	if !enabled {
		return s
	}
	return fmt.Sprintf("\033[%sm%s\033[0m", code, s)
}
