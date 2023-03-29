package amapretty

import (
	"encoding/json"
	"fmt"
	"runtime"
	"time"
)

var (
	black   = fmtColor("\033[1;30m%s\033[0m")
	red     = fmtColor("\033[1;31m%s\033[0m")
	green   = fmtColor("\033[1;32m%s\033[0m")
	yellow  = fmtColor("\033[1;33m%s\033[0m")
	purple  = fmtColor("\033[1;34m%s\033[0m")
	magenta = fmtColor("\033[1;35m%s\033[0m")
	teal    = fmtColor("\033[1;36m%s\033[0m")
	white   = fmtColor("\033[1;37m%s\033[0m")
)

var timeNow = time.Now
var runtimeCaller = runtime.Caller

func fmtColor(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}

func Print(args ...interface{}) {
	timeNow := timeNow().Format("01-02-2006 15:04:05")
	fmtTimeNow := purple(timeNow)

	prefix := fmt.Sprintf("[%s] %s -- ", "PrettyPrint", fmtTimeNow)
	fmtPrefix := green(prefix)

	_, fileName, fileLine, ok := runtimeCaller(1)

	caller := ""
	if ok {
		caller = fmt.Sprintf("%s:%d", fileName, fileLine)
		caller = teal(caller)
	}

	fmt.Printf("\n%s%s\n", fmtPrefix, caller)

	s, _ := json.MarshalIndent(args, "", "\t")
	fmt.Printf("%s%s\n", fmtPrefix, string(s))
}
