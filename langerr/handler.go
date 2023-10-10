package langerr

import (
	"fmt"
)

type Handler struct {
	HadError bool
}

func (h *Handler) Report(line int, where string, message string) {
	if len(where) > 0 {
		fmt.Printf("[line %d] Error %s: %s\n", line, where, message)
	} else {
		fmt.Printf("[line %d] Error%s: %s\n", line, where, message)
	}
	h.HadError = true
}
