package lang

/******************************************************************************
 * Any function that is callable will need to implement this interface.
 *****************************************************************************/

type callable interface {
	arity() int
	call(interpreter *Interpreter, args []any) any
	toString() string
}
