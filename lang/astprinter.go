package lang

import "fmt"

/******************************************************************************
 * Helper struct to display the AST and expression operation precendence in
 * the earlier stages of development.
 *****************************************************************************/

type AstPrinter struct{}

func (printer AstPrinter) Print(expr Expr) string {
	return expr.accept(printer).(string)
}

func (printer AstPrinter) visitAssignExpr(expr AssignExpr) any {
	panic("AstPrinter is not able to print assignment expressions at this time.")
}

func (printer AstPrinter) visitBinaryExpr(expr BinaryExpr) any {
	return printer.parenthesize(expr.operator.lexeme, expr.left, expr.right)
}

func (printer AstPrinter) visitCallExpr(expr CallExpr) any {
	panic("AstPrinter is not able to print call expressions at this time.")
}

func (printer AstPrinter) visitGetExpr(expr GetExpr) any {
	panic("AstPrinter is not able to print get expressions at this time.")
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

func (printer AstPrinter) visitLogicalExpr(expr LogicalExpr) any {
	return printer.parenthesize(expr.operator.lexeme, expr.left, expr.right)
}

func (printer AstPrinter) visitSetExpr(expr SetExpr) any {
	panic("AstPrinter is not able to print set expressions at this time.")
}

func (printer AstPrinter) visitThisExpr(expr ThisExpr) any {
	panic("AstPrinter is not able to print this expressions at this time.")
}

func (printer AstPrinter) visitUnaryExpr(expr UnaryExpr) any {
	return printer.parenthesize(expr.operator.lexeme, expr.right)
}

func (printer AstPrinter) visitVariableExpr(expr VariableExpr) any {
	panic("AstPrinter is not able to print variable expressions at this time.")
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
