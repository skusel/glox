package lang

/******************************************************************************
 * Expresssion definitions. Expressions are nodes of the AST.
 *****************************************************************************/

type Expr interface {
	accept(exprVisitor exprVisitor) any
}

type exprVisitor interface {
	visitAssignExpr(a AssignExpr) any
	visitBinaryExpr(b BinaryExpr) any
	visitGroupingExpr(g GroupingExpr) any
	visitLiteralExpr(l LiteralExpr) any
	visitLogicalExpr(l LogicalExpr) any
	visitUnaryExpr(u UnaryExpr) any
	visitVariableExpr(v VariableExpr) any
}

type AssignExpr struct {
	name  Token
	value Expr
}

func (a AssignExpr) accept(visitor exprVisitor) any {
	return visitor.visitAssignExpr(a)
}

type BinaryExpr struct {
	left     Expr
	operator Token
	right    Expr
}

func (b BinaryExpr) accept(visitor exprVisitor) any {
	return visitor.visitBinaryExpr(b)
}

type GroupingExpr struct {
	expression Expr
}

func (g GroupingExpr) accept(visitor exprVisitor) any {
	return visitor.visitGroupingExpr(g)
}

type LiteralExpr struct {
	value any
}

func (l LiteralExpr) accept(visitor exprVisitor) any {
	return visitor.visitLiteralExpr(l)
}

type LogicalExpr struct {
	left     Expr
	operator Token
	right    Expr
}

func (l LogicalExpr) accept(visitor exprVisitor) any {
	return visitor.visitLogicalExpr(l)
}

type UnaryExpr struct {
	operator Token
	right    Expr
}

func (u UnaryExpr) accept(visitor exprVisitor) any {
	return visitor.visitUnaryExpr(u)
}

type VariableExpr struct {
	name Token
}

func (v VariableExpr) accept(visitor exprVisitor) any {
	return visitor.visitVariableExpr(v)
}
