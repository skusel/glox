package lang

/******************************************************************************
 * Statement definitions. Statements are nodes of the AST.
 *****************************************************************************/

type Stmt interface {
	accept(stmtVisitor stmtVisitor) any
}

type stmtVisitor interface {
	visitBlockStmt(stmt BlockStmt) any
	visitExprStmt(stmt ExprStmt) any
	visitPrintStmt(stmt PrintStmt) any
	visitVarStmt(stmt VarStmt) any
}

type BlockStmt struct {
	statements []Stmt
}

func (stmt BlockStmt) accept(visitor stmtVisitor) any {
	return visitor.visitBlockStmt(stmt)
}

type ExprStmt struct {
	expr Expr
}

func (stmt ExprStmt) accept(visitor stmtVisitor) any {
	return visitor.visitExprStmt(stmt)
}

type PrintStmt struct {
	expr Expr
}

func (stmt PrintStmt) accept(visitor stmtVisitor) any {
	return visitor.visitPrintStmt(stmt)
}

type VarStmt struct {
	name        Token
	initializer Expr
}

func (stmt VarStmt) accept(visitor stmtVisitor) any {
	return visitor.visitVarStmt(stmt)
}
