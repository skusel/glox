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

func (env *environment) ancestor(distance int) *environment {
	ancestorEnv := env
	for i := 0; i < distance; i++ {
		ancestorEnv = ancestorEnv.enclosing
	}
	return ancestorEnv
}

func (env *environment) getAt(distance int, name Token) any {
	value, found := env.ancestor(distance).values[name.lexeme]
	if found {
		return value
	} else {
		env.errorHandler.reportRuntimeError(name.line, errors.New("Undefined variable '"+name.lexeme+"'."))
		return nil
	}
}

func (env *environment) getThisValue() any {
	// if this is called, we already checked that we are in a method
	return env.values["this"]
}

func (env *environment) getSubClassThisValue(distance int) any {
	// if this is called, we already checked that we are in a super class
	return env.ancestor(distance - 1).values["this"]
}

func (env *environment) get(name Token) any {
	value, found := env.values[name.lexeme]
	if found {
		return value
	} else if env.enclosing != nil {
		return env.enclosing.get(name)
	} else {
		env.errorHandler.reportRuntimeError(name.line, errors.New("Undefined variable '"+name.lexeme+"'."))
		return nil
	}
}

func (env *environment) assignAt(distance int, name Token, value any) {
	env.ancestor(distance).values[name.lexeme] = value
}

func (env *environment) assign(name Token, value any) {
	_, found := env.values[name.lexeme]
	if found {
		env.values[name.lexeme] = value
	} else if env.enclosing != nil {
		env.enclosing.assign(name, value)
	} else {
		env.errorHandler.reportRuntimeError(name.line, errors.New("Undefined variable '"+name.lexeme+"'."))
	}
}
