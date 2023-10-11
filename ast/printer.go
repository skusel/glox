package ast

import "fmt"

type Printer struct{}

func (printer Printer) Print(expr Expr) string {
	return expr.Accept(printer).(string)
}

func (printer Printer) VisitBinaryExpr(expr BinaryExpr) any {
	return printer.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (printer Printer) VisitGroupingExpr(expr GroupingExpr) any {
	return printer.parenthesize("group", expr.Expression)
}

func (printer Printer) VisitLiteralExpr(expr LiteralExpr) any {
	if expr.Value == nil {
		return "nil"
	}
	return fmt.Sprint(expr.Value)
}

func (printer Printer) VisitUnaryExpr(expr UnaryExpr) any {
	return printer.parenthesize(expr.Operator.Lexeme, expr.Right)
}

func (printer Printer) parenthesize(name string, exprs ...Expr) string {
	prettyString := "(" + name
	for _, expr := range exprs {
		prettyString += " "
		prettyString += expr.Accept(printer).(string)
	}
	prettyString += ")"
	return prettyString
}
