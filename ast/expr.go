package ast

type Expr interface {
	Accept(exprVisitor ExprVisitor) any
}

type ExprVisitor interface {
	VisitBinaryExpr(b BinaryExpr) any
	VisitGroupingExpr(g GroupingExpr) any
	VisitLiteralExpr(l LiteralExpr) any
	VisitUnaryExpr(u UnaryExpr) any
}

type BinaryExpr struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (b BinaryExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitBinaryExpr(b)
}

type GroupingExpr struct {
	Expression Expr
}

func (g GroupingExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitGroupingExpr(g)
}

type LiteralExpr struct {
	Value any
}

func (l LiteralExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitLiteralExpr(l)
}

type UnaryExpr struct {
	Operator Token
	Right    Expr
}

func (u UnaryExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitUnaryExpr(u)
}
