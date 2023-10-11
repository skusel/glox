package langerr

import (
	"fmt"
)

type Handler struct {
	HadError        bool
	HadRuntimeError bool
}

func NewHandler() *Handler {
	return &Handler{HadError: false, HadRuntimeError: false}
}

func (h *Handler) Report(line int, where string, err error) {
	if len(where) > 0 {
		fmt.Printf("[line %d] Error %s: %s\n", line, where, err)
	} else {
		fmt.Printf("[line %d] Error%s: %s\n", line, where, err)
	}
	h.HadError = true
}

func (h *Handler) ReportRuntime(line int, err error) {
	fmt.Printf("%s\n[line %d]\n", err, line)
	h.HadRuntimeError = true
}
