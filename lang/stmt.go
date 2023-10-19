package lang

/******************************************************************************
 * Statement definitions. Statements are nodes of the AST.
 *****************************************************************************/

type Stmt interface {
	accept(stmtVisitor stmtVisitor) any
}

type stmtVisitor interface {
	visitBlockStmt(stmt BlockStmt) any
	visitClassStmt(stmt ClassStmt) any
	visitExprStmt(stmt ExprStmt) any
	visitFunctionStmt(stmt FunctionStmt) any
	visitIfStmt(stmt IfStmt) any
	visitPrintStmt(stmt PrintStmt) any
	visitReturnStmt(stmt ReturnStmt) any
	visitVarStmt(stmt VarStmt) any
	visitWhileStmt(stmt WhileStmt) any
}

type BlockStmt struct {
	statements []Stmt
}

func (stmt BlockStmt) accept(visitor stmtVisitor) any {
	return visitor.visitBlockStmt(stmt)
}

type ClassStmt struct {
	name       Token
	superclass VariableExpr
	methods    []FunctionStmt
}

func (stmt ClassStmt) accept(visitor stmtVisitor) any {
	return visitor.visitClassStmt(stmt)
}

type ExprStmt struct {
	expr Expr
}

func (stmt ExprStmt) accept(visitor stmtVisitor) any {
	return visitor.visitExprStmt(stmt)
}

type FunctionStmt struct {
	name   Token
	params []Token
	body   []Stmt
}

func (stmt FunctionStmt) accept(visitor stmtVisitor) any {
	return visitor.visitFunctionStmt(stmt)
}

type IfStmt struct {
	condition  Expr
	thenBranch Stmt
	elseBranch Stmt
}

func (stmt IfStmt) accept(visitor stmtVisitor) any {
	return visitor.visitIfStmt(stmt)
}

type PrintStmt struct {
	expr Expr
}

func (stmt PrintStmt) accept(visitor stmtVisitor) any {
	return visitor.visitPrintStmt(stmt)
}

type ReturnStmt struct {
	keyword Token
	value   Expr
}

func (stmt ReturnStmt) accept(visitor stmtVisitor) any {
	return visitor.visitReturnStmt(stmt)
}

type VarStmt struct {
	name        Token
	initializer Expr
}

func (stmt VarStmt) accept(visitor stmtVisitor) any {
	return visitor.visitVarStmt(stmt)
}

type WhileStmt struct {
	condition Expr
	body      Stmt
}

func (stmt WhileStmt) accept(visitor stmtVisitor) any {
	return visitor.visitWhileStmt(stmt)
}
