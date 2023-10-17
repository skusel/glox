package lang

/******************************************************************************
 * Expresssion definitions. Expressions are nodes of the AST.
 *
 * Expression IDs are populated by the parser. They are uniquely assigned
 * whenever any expression is created so that the resolver and interpreter are
 * able to recognize when they are referring to the same expression.
 *****************************************************************************/

type Expr interface {
	getId() int
	accept(exprVisitor exprVisitor) any
}

type exprVisitor interface {
	visitAssignExpr(a AssignExpr) any
	visitBinaryExpr(b BinaryExpr) any
	visitCallExpr(c CallExpr) any
	visitGroupingExpr(g GroupingExpr) any
	visitLiteralExpr(l LiteralExpr) any
	visitLogicalExpr(l LogicalExpr) any
	visitUnaryExpr(u UnaryExpr) any
	visitVariableExpr(v VariableExpr) any
}

type AssignExpr struct {
	id    int
	name  Token
	value Expr
}

func (a AssignExpr) getId() int {
	return a.id
}

func (a AssignExpr) accept(visitor exprVisitor) any {
	return visitor.visitAssignExpr(a)
}

type BinaryExpr struct {
	id       int
	left     Expr
	operator Token
	right    Expr
}

func (b BinaryExpr) getId() int {
	return b.id
}

func (b BinaryExpr) accept(visitor exprVisitor) any {
	return visitor.visitBinaryExpr(b)
}

type CallExpr struct {
	id     int
	callee Expr
	paren  Token
	args   []Expr
}

func (c CallExpr) getId() int {
	return c.id
}

func (c CallExpr) accept(visitor exprVisitor) any {
	return visitor.visitCallExpr(c)
}

type GroupingExpr struct {
	id         int
	expression Expr
}

func (g GroupingExpr) getId() int {
	return g.id
}

func (g GroupingExpr) accept(visitor exprVisitor) any {
	return visitor.visitGroupingExpr(g)
}

type LiteralExpr struct {
	id    int
	value any
}

func (l LiteralExpr) getId() int {
	return l.id
}

func (l LiteralExpr) accept(visitor exprVisitor) any {
	return visitor.visitLiteralExpr(l)
}

type LogicalExpr struct {
	id       int
	left     Expr
	operator Token
	right    Expr
}

func (l LogicalExpr) getId() int {
	return l.id
}

func (l LogicalExpr) accept(visitor exprVisitor) any {
	return visitor.visitLogicalExpr(l)
}

type UnaryExpr struct {
	id       int
	operator Token
	right    Expr
}

func (u UnaryExpr) getId() int {
	return u.id
}

func (u UnaryExpr) accept(visitor exprVisitor) any {
	return visitor.visitUnaryExpr(u)
}

type VariableExpr struct {
	id   int
	name Token
}

func (v VariableExpr) getId() int {
	return v.id
}

func (v VariableExpr) accept(visitor exprVisitor) any {
	return visitor.visitVariableExpr(v)
}
