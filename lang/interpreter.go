package lang

import (
	"errors"
	"fmt"
	"reflect"
)

type Interpreter struct {
	errorHandler *ErrorHandler
}

func NewInterpreter(errorHandler *ErrorHandler) *Interpreter {
	return &Interpreter{errorHandler: errorHandler}
}

func (interperter *Interpreter) Interpret(expr Expr) {
	value := interperter.evaluate(expr)
	if interperter.errorHandler.HadRuntimeError {
		return
	} else {
		fmt.Println(stringify(value))
	}
}

func (interpreter *Interpreter) evaluate(expr Expr) any {
	return expr.accept(interpreter)
}

func (interpreter *Interpreter) visitBinaryExpr(expr BinaryExpr) any {
	left := interpreter.evaluate(expr.left)
	right := interpreter.evaluate(expr.right)

	switch expr.operator.tokenType {
	case tokenTypeGreater:
		valid, leftFloat, rightFloat := areValuesValidFloats(left, right)
		if !valid {
			err := errors.New("Operands must be numbers when using the '>' operator.")
			interpreter.errorHandler.reportRuntime(expr.operator.line, err)
		}
		return leftFloat > rightFloat
	case tokenTypeGreaterEqual:
		valid, leftFloat, rightFloat := areValuesValidFloats(left, right)
		if !valid {
			err := errors.New("Operands must be numbers when using the '>=' operator.")
			interpreter.errorHandler.reportRuntime(expr.operator.line, err)
		}
		return leftFloat >= rightFloat
	case tokenTypeLess:
		valid, leftFloat, rightFloat := areValuesValidFloats(left, right)
		if !valid {
			err := errors.New("Operands must be numbers when using the '<' operator.")
			interpreter.errorHandler.reportRuntime(expr.operator.line, err)
		}
		return leftFloat < rightFloat
	case tokenTypeLessEqual:
		valid, leftFloat, rightFloat := areValuesValidFloats(left, right)
		if !valid {
			err := errors.New("Operands must be numbers when using the '<=' operator.")
			interpreter.errorHandler.reportRuntime(expr.operator.line, err)
		}
		return leftFloat <= rightFloat
	case tokenTypeMinus:
		valid, leftFloat, rightFloat := areValuesValidFloats(left, right)
		if !valid {
			err := errors.New("Operands must be numbers when using the '-' operator.")
			interpreter.errorHandler.reportRuntime(expr.operator.line, err)
		}
		return leftFloat - rightFloat
	case tokenTypePlus:
		validFloats, leftFloat, rightFloat := areValuesValidFloats(left, right)
		if validFloats {
			return leftFloat + rightFloat
		}
		validStrings, leftString, rightString := areValuesValidStrings(left, right)
		if validStrings {
			return leftString + rightString
		}
		err := errors.New("Operands must be numbers or strings and be the same type when using the '+' operator.")
		interpreter.errorHandler.reportRuntime(expr.operator.line, err)
	case tokenTypeSlash:
		valid, leftFloat, rightFloat := areValuesValidFloats(left, right)
		if !valid {
			err := errors.New("Operands must be numbers when using the '/' operator.")
			interpreter.errorHandler.reportRuntime(expr.operator.line, err)
		}
		return leftFloat / rightFloat
	case tokenTypeStar:
		valid, leftFloat, rightFloat := areValuesValidFloats(left, right)
		if !valid {
			err := errors.New("Operands must be numbers when using the '*' operator.")
			interpreter.errorHandler.reportRuntime(expr.operator.line, err)
		}
		return leftFloat * rightFloat
	case tokenTypeEqualEqual:
		return reflect.DeepEqual(left, right)
	case tokenTypeBangEqual:
		return !reflect.DeepEqual(left, right)
	}

	// unreachable
	return nil
}

func (interpreter *Interpreter) visitGroupingExpr(expr GroupingExpr) any {
	return interpreter.evaluate(expr.expression)
}

func (interperter *Interpreter) visitLiteralExpr(expr LiteralExpr) any {
	return expr.value
}

func (interpreter *Interpreter) visitUnaryExpr(expr UnaryExpr) any {
	right := interpreter.evaluate(expr.right)
	switch expr.operator.tokenType {
	case tokenTypeBang:
		return !isTruthy(right)
	case tokenTypeMinus:
		rightFloat, rightFloatValid := right.(float64)
		if !rightFloatValid {
			err := errors.New("Operand must be a number.")
			interpreter.errorHandler.reportRuntime(expr.operator.line, err)
		}
		return -1 * rightFloat
	}
	return nil
}

func areValuesValidFloats(left, right any) (bool, float64, float64) {
	leftFloat, leftFloatValid := left.(float64)
	rightFloat, rightFloatValid := right.(float64)
	return leftFloatValid && rightFloatValid, leftFloat, rightFloat
}

func areValuesValidStrings(left, right any) (bool, string, string) {
	leftString, leftStringValid := left.(string)
	rightString, rightStringValid := right.(string)
	return leftStringValid && rightStringValid, leftString, rightString
}

func isTruthy(value any) bool {
	if value == nil {
		return false
	}
	boolVal, isBool := value.(bool)
	if isBool {
		return boolVal
	}
	return false
}

func stringify(value any) string {
	if value == nil {
		return "nil"
	}
	return fmt.Sprint(value)
}
