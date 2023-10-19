package lang

/******************************************************************************
 * The class struct is used to represent classes in Lox. class implements the
 * callable interface (that's how classes are instantiated).
 *****************************************************************************/

type class struct {
	name         string
	superclass   *class
	methods      map[string]function
	errorHandler *ErrorHandler
}

func (c class) arity() int {
	initializer, hasInitializer := c.findMethod("init").(function)
	if hasInitializer {
		return initializer.arity()
	}
	return 0
}

func (c class) call(interpreter *Interpreter, args []any) any {
	inst := instance{class: c, fields: make(map[string]any), errorHandler: c.errorHandler}
	initializer, hasInitializer := c.findMethod("init").(function)
	if hasInitializer {
		initializer.bind(inst).call(interpreter, args)
	}
	return inst
}

func (c class) findMethod(name string) any {
	method, foundMethod := c.methods[name]
	if foundMethod {
		return method
	} else if c.superclass != nil {
		return c.superclass.findMethod(name)
	} else {
		return nil
	}
}

func (c class) toString() string {
	return c.name
}
