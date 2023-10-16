package lang

/******************************************************************************
 * function implements the callable interface. It is used to represent
 * function calls in the interpreter's runtime.
 *****************************************************************************/

type returnContent struct {
	value any
}

type function struct {
	declaration FunctionStmt
	closure     *environment
}

func (fun function) arity() int {
	return len(fun.declaration.params)
}

func (fun function) call(interpreter *Interpreter, args []any) (value any) {
	defer func() {
		/**********************************************************************
		 * This is a hacky way of unwinding the call stack that is created
		 * within executeBlock when a return statement is hit.
		 *********************************************************************/
		err := recover()
		if err != nil {
			returnContent, isReturnContent := err.(returnContent)
			if isReturnContent {
				// update the return value to be the called functions return value
				value = returnContent.value
			} else {
				// this is not a panic thrown by us, pass it on
				panic(err)
			}
		}
	}()

	funEnv := newChildEnvironment(fun.closure)
	for i, param := range fun.declaration.params {
		funEnv.define(param.lexeme, args[i])
	}
	interpreter.executeBlock(fun.declaration.body, funEnv)
	return nil
}

func (fun function) toString() string {
	return "<fun " + fun.declaration.name.lexeme + ">"
}
