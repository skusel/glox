package lang

import (
	"fmt"
)

type ErrorHandler struct {
	HadError        bool
	HadRuntimeError bool
}

func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{HadError: false, HadRuntimeError: false}
}

func (h *ErrorHandler) report(line int, where string, err error) {
	if len(where) > 0 {
		fmt.Printf("[line %d] Error %s: %s\n", line, where, err)
	} else {
		fmt.Printf("[line %d] Error%s: %s\n", line, where, err)
	}
	h.HadError = true
}

func (h *ErrorHandler) reportRuntime(line int, err error) {
	fmt.Printf("%s\n[line %d]\n", err, line)
	h.HadRuntimeError = true
}
