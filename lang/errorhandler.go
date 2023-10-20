package lang

import (
	"fmt"
	"os"
)

/******************************************************************************
 * Helper struct to assist with error reporting.
 *
 * Panics are used in a few spots in the interpreter implementation to unwind
 * the call stack. This unwinding is often the easiest solution given the
 * recursive nature of the parser, reolver, and interpreter implementations.
 *****************************************************************************/

type ErrorHandler struct {
	HadError        bool
	HadRuntimeError bool
}

type staticError struct {
	msg string
}

type runtimeError struct {
	msg string
}

func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{HadError: false, HadRuntimeError: false}
}

func (h *ErrorHandler) reportStaticError(line int, where string, err error, synchronize bool) {
	h.HadError = true
	var errorMsg string
	if len(where) > 0 {
		errorMsg = fmt.Sprintf("[line %d] Error %s: %s\n", line, where, err)
	} else {
		errorMsg = fmt.Sprintf("[line %d] Error: %s\n", line, err)
	}
	staticError := staticError{msg: errorMsg}
	if synchronize {
		// panic will unwind the call stack and we can "catch" the error with recover()
		panic(staticError)
	} else {
		// if we are not syncing, immediately report the error to stderr
		os.Stderr.WriteString(staticError.msg)
	}
}

func (h *ErrorHandler) reportRuntimeError(line int, err error) {
	h.HadRuntimeError = true
	runtimeError := runtimeError{msg: fmt.Sprintf("[line %d] %s\n", line, err)}
	// we always want to unwind the call stack and recover for runtime errors
	panic(runtimeError)
}
