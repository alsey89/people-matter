package extractor

import (
	"runtime"
	"strings"
)

// GetCaller returns the name of the function that called it.
// Use skip to specify the level.
// 0 is the caller of GetCaller.
func GetCaller(skip int) string {
	pc, _, _, ok := runtime.Caller(skip)
	if !ok {
		return "unknown"
	}
	fn := runtime.FuncForPC(pc)
	return fn.Name()
}

/*
Retrieves function names up to the specified depth and concatenates them with a colon separator.
Depth specifies the number of functions to trace. If depth is 0, it returns "no callers".
example:

	func A() {
		B()
	}

	func B() {
		C()
	}

	func C() {
		log.Println(extractor.GetCallerTrace(3))
	}

Prints: A:B:C:
*/
func GetCallerTrace(depth int) string {
	if depth == 0 {
		return "no callers"
	}

	var callers []string
	for i := 1; i <= depth; i++ {
		caller := GetCaller(i)
		callers = append(callers, caller)
	}
	return strings.Join(callers, ":") + ":"
}
