package lang

import "time"

/******************************************************************************
 * structs in this file should implement the callable interface. Each struct
 * represents a native function call. That is, a function all that is built
 * into the language.
 *****************************************************************************/

type clock struct{}

func (c clock) arity() int {
	return 0
}

func (c clock) call(interpreter *Interpreter, args []any) any {
	return time.Now()
}

func (c clock) toString() string {
	return "<native fun>"
}
