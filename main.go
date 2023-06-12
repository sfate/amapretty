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

func Print(args ...interface{}) {
	timeNow := timeNow().Format(time.RFC3339)
	fmtTimeNow := fmt.Sprintf("\033[1;34m%s\033[0m", timeNow)
	fmtPrefix := fmt.Sprintf("\033[1;32m%s\033[0m", prefix)

	_, fileName, fileLine, ok := runtimeCaller(1)
	caller := ""
	if ok {
		caller = fmt.Sprintf("%s:%d", fileName, fileLine)
		caller = fmt.Sprintf("\033[1;36m%s\033[0m", caller)
	}

	s, _ := json.MarshalIndent(args, "", "\t")
	fmt.Fprintf(output, "[%s] %s %s -- %s\n", fmtPrefix, fmtTimeNow, caller, string(s))
}
