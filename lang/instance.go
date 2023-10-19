package lang

import "errors"

type instance struct {
	class        class
	fields       map[string]any
	errorHandler *ErrorHandler
}

func newInstance(class class, errorHandler *ErrorHandler) instance {
	return instance{class: class, fields: make(map[string]any), errorHandler: errorHandler}
}

func (inst instance) get(name Token) any {
	fieldValue, hasField := inst.fields[name.lexeme]
	if hasField {
		return fieldValue
	}
	method, hasMethod := inst.class.findMethod(name.lexeme).(function)
	if hasMethod {
		return method.bind(inst)
	}
	err := errors.New("Undefined property '" + name.lexeme + "'.")
	inst.errorHandler.reportRuntimeError(name.line, err)
	return nil
}

func (inst instance) set(name Token, value any) {
	inst.fields[name.lexeme] = value
}

func (inst instance) toString() string {
	return inst.class.name + " instance"
}
