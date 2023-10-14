package lang

import (
	"fmt"
)

/******************************************************************************
 * Helper struct to assist with error reporting.
 *****************************************************************************/

type ErrorHandler struct {
	needToSynchronize bool // used to report multiple static errors
	HadError          bool
	HadRuntimeError   bool
}

func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{needToSynchronize: false, HadError: false, HadRuntimeError: false}
}

func (h *ErrorHandler) report(line int, where string, err error) {
	if h.needToSynchronize {
		// currently out of sync - don't report another error
		return
	}
	if len(where) > 0 {
		fmt.Printf("[line %d] Error %s: %s\n", line, where, err)
	} else {
		fmt.Printf("[line %d] Error%s: %s\n", line, where, err)
	}
	h.HadError = true
	h.needToSynchronize = true
}

func (h *ErrorHandler) reportRuntime(line int, err error) {
	fmt.Printf("[line %d] %s\n", line, err)
	h.HadRuntimeError = true
}
