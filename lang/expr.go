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
	visitGetExpr(g GetExpr) any
	visitGroupingExpr(g GroupingExpr) any
	visitLiteralExpr(l LiteralExpr) any
	visitLogicalExpr(l LogicalExpr) any
	visitSetExpr(s SetExpr) any
	visitSuperExpr(s SuperExpr) any
	visitThisExpr(t ThisExpr) any
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

type GetExpr struct {
	id     int
	object Expr
	name   Token
}

func (g GetExpr) getId() int {
	return g.id
}

func (g GetExpr) accept(visitor exprVisitor) any {
	return visitor.visitGetExpr(g)
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

type SetExpr struct {
	id     int
	object Expr
	name   Token
	value  Expr
}

func (s SetExpr) getId() int {
	return s.id
}

func (s SetExpr) accept(visitor exprVisitor) any {
	return visitor.visitSetExpr(s)
}

type SuperExpr struct {
	id      int
	keyword Token
	method  Token
}

func (s SuperExpr) getId() int {
	return s.id
}

func (s SuperExpr) accept(visitor exprVisitor) any {
	return visitor.visitSuperExpr(s)
}

type ThisExpr struct {
	id      int
	keyword Token
}

func (t ThisExpr) getId() int {
	return t.id
}

func (t ThisExpr) accept(visitor exprVisitor) any {
	return visitor.visitThisExpr(t)
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
