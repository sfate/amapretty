package amapretty

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"time"
)

var timeNow = time.Now
var runtimeCaller = runtime.Caller
var output = os.Stdout

const (
	prefix = "amapretty"
)

// Print writes values to stdout as indented JSON with a timestamp and caller reference.
func Print(args ...interface{}) {
	printWithCallerSkip(1, args...)
}

// Printf formats according to a format specifier and writes the result through Print.
func Printf(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	printWithCallerSkip(2, s)
}

func printWithCallerSkip(skip int, args ...interface{}) {
	timeNow := timeNow().Format(time.RFC3339)
	fmtTimeNow := fmt.Sprintf("\033[1;34m%s\033[0m", timeNow)
	fmtPrefix := fmt.Sprintf("\033[1;32m%s\033[0m", prefix)

	_, fileName, fileLine, ok := runtimeCaller(skip)
	caller := ""
	if ok {
		caller = fmt.Sprintf("%s:%d", fileName, fileLine)
		caller = fmt.Sprintf("\033[1;36m%s\033[0m", caller)
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
	_, _ = fmt.Fprintf(output, "[%s] %s %s -- %s\n", fmtPrefix, fmtTimeNow, caller, string(s))
}
