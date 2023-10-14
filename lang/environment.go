package lang

import "errors"

/******************************************************************************
 * The language's environment tracks and stores variables and their values.
 * The origins of the term "environment" in this context date back to the
 * authors of Lisp.
 *****************************************************************************/

type environment struct {
	enclosing    *environment
	values       map[string]any
	errorHandler *ErrorHandler
}

func newEnvironment(errorHandler *ErrorHandler) *environment {
	return &environment{enclosing: nil, values: make(map[string]any), errorHandler: errorHandler}
}

func newChildEnvironment(parentEnv *environment) *environment {
	return &environment{enclosing: parentEnv, values: make(map[string]any), errorHandler: parentEnv.errorHandler}
}

func (env *environment) define(name string, value any) {
	env.values[name] = value
}

func (env *environment) get(name Token) any {
	value, found := env.values[name.lexeme]
	if found {
		return value
	} else if env.enclosing != nil {
		return env.enclosing.get(name)
	} else {
		env.errorHandler.reportRuntime(name.line, errors.New("Undefined variable '"+name.lexeme+"'."))
		return nil
	}
}

func (env *environment) assign(name Token, value any) {
	_, found := env.values[name.lexeme]
	if found {
		env.values[name.lexeme] = value
	} else if env.enclosing != nil {
		env.enclosing.assign(name, value)
	} else {
		env.errorHandler.reportRuntime(name.line, errors.New("Undefined variable '"+name.lexeme+"'."))
	}
}
