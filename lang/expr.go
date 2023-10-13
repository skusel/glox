package lang

type Expr interface {
	accept(exprVisitor exprVisitor) any
}

type exprVisitor interface {
	visitBinaryExpr(b BinaryExpr) any
	visitGroupingExpr(g GroupingExpr) any
	visitLiteralExpr(l LiteralExpr) any
	visitUnaryExpr(u UnaryExpr) any
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

type UnaryExpr struct {
	operator Token
	right    Expr
}

func (u UnaryExpr) accept(visitor exprVisitor) any {
	return visitor.visitUnaryExpr(u)
}
