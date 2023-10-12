package interpreter

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/skusel/glox/ast"
	"github.com/skusel/glox/gloxerror"
)

type Interpreter struct {
	errorHandler *gloxerror.Handler
}

func NewInterpreter(errorHandler *gloxerror.Handler) *Interpreter {
	return &Interpreter{errorHandler: errorHandler}
}

func (interperter *Interpreter) Interpret(expr ast.Expr) {
	value := interperter.evaluate(expr)
	if interperter.errorHandler.HadRuntimeError {
		return
	} else {
		fmt.Println(stringify(value))
	}
}

func (interpreter *Interpreter) evaluate(expr ast.Expr) any {
	return expr.Accept(interpreter)
}

func (interpreter *Interpreter) VisitBinaryExpr(expr ast.BinaryExpr) any {
	left := interpreter.evaluate(expr.Left)
	right := interpreter.evaluate(expr.Right)

	switch expr.Operator.TokenType {
	case ast.Greater:
		valid, leftFloat, rightFloat := areValuesValidFloats(left, right)
		if !valid {
			err := errors.New("Operands must be numbers when using the '>' operator.")
			interpreter.errorHandler.ReportRuntime(expr.Operator.Line, err)
		}
		return leftFloat > rightFloat
	case ast.GreaterEqual:
		valid, leftFloat, rightFloat := areValuesValidFloats(left, right)
		if !valid {
			err := errors.New("Operands must be numbers when using the '>=' operator.")
			interpreter.errorHandler.ReportRuntime(expr.Operator.Line, err)
		}
		return leftFloat >= rightFloat
	case ast.Less:
		valid, leftFloat, rightFloat := areValuesValidFloats(left, right)
		if !valid {
			err := errors.New("Operands must be numbers when using the '<' operator.")
			interpreter.errorHandler.ReportRuntime(expr.Operator.Line, err)
		}
		return leftFloat < rightFloat
	case ast.LessEqual:
		valid, leftFloat, rightFloat := areValuesValidFloats(left, right)
		if !valid {
			err := errors.New("Operands must be numbers when using the '<=' operator.")
			interpreter.errorHandler.ReportRuntime(expr.Operator.Line, err)
		}
		return leftFloat <= rightFloat
	case ast.Minus:
		valid, leftFloat, rightFloat := areValuesValidFloats(left, right)
		if !valid {
			err := errors.New("Operands must be numbers when using the '-' operator.")
			interpreter.errorHandler.ReportRuntime(expr.Operator.Line, err)
		}
		return leftFloat - rightFloat
	case ast.Plus:
		validFloats, leftFloat, rightFloat := areValuesValidFloats(left, right)
		if validFloats {
			return leftFloat + rightFloat
		}
		validStrings, leftString, rightString := areValuesValidStrings(left, right)
		if validStrings {
			return leftString + rightString
		}
		err := errors.New("Operands must be numbers or strings and be the same type when using the '+' operator.")
		interpreter.errorHandler.ReportRuntime(expr.Operator.Line, err)
	case ast.Slash:
		valid, leftFloat, rightFloat := areValuesValidFloats(left, right)
		if !valid {
			err := errors.New("Operands must be numbers when using the '/' operator.")
			interpreter.errorHandler.ReportRuntime(expr.Operator.Line, err)
		}
		return leftFloat / rightFloat
	case ast.Star:
		valid, leftFloat, rightFloat := areValuesValidFloats(left, right)
		if !valid {
			err := errors.New("Operands must be numbers when using the '*' operator.")
			interpreter.errorHandler.ReportRuntime(expr.Operator.Line, err)
		}
		return leftFloat * rightFloat
	case ast.EqualEqual:
		return reflect.DeepEqual(left, right)
	case ast.BangEqual:
		return !reflect.DeepEqual(left, right)
	}

	// unreachable
	return nil
}

func (interpreter *Interpreter) VisitGroupingExpr(expr ast.GroupingExpr) any {
	return interpreter.evaluate(expr.Expression)
}

func (interperter *Interpreter) VisitLiteralExpr(expr ast.LiteralExpr) any {
	return expr.Value
}

func (interpreter *Interpreter) VisitUnaryExpr(expr ast.UnaryExpr) any {
	right := interpreter.evaluate(expr.Right)
	switch expr.Operator.TokenType {
	case ast.Bang:
		return !isTruthy(right)
	case ast.Minus:
		rightFloat, rightFloatValid := right.(float64)
		if !rightFloatValid {
			err := errors.New("Operand must be a number.")
			interpreter.errorHandler.ReportRuntime(expr.Operator.Line, err)
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
