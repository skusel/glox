package lang

import "fmt"

type AstPrinter struct{}

func (printer AstPrinter) Print(expr Expr) string {
	return expr.accept(printer).(string)
}

func (printer AstPrinter) visitBinaryExpr(expr BinaryExpr) any {
	return printer.parenthesize(expr.operator.lexeme, expr.left, expr.right)
}

func (printer AstPrinter) visitGroupingExpr(expr GroupingExpr) any {
	return printer.parenthesize("group", expr.expression)
}

func (printer AstPrinter) visitLiteralExpr(expr LiteralExpr) any {
	if expr.value == nil {
		return "nil"
	}
	return fmt.Sprint(expr.value)
}

func (printer AstPrinter) visitUnaryExpr(expr UnaryExpr) any {
	return printer.parenthesize(expr.operator.lexeme, expr.right)
}

func (printer AstPrinter) parenthesize(name string, exprs ...Expr) string {
	prettyString := "(" + name
	for _, expr := range exprs {
		prettyString += " "
		prettyString += expr.accept(printer).(string)
	}
	prettyString += ")"
	return prettyString
}
