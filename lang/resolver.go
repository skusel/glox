package lang

import (
	"errors"
)

/******************************************************************************
 * The in this glox implementation the resolver performs a separate pass on
 * the AST before it is passed to the interpreter. It performs semantic
 * analysis. That includes, calculating how many hops away the declared
 * variable is in the environment chain (i.e. resolving variables), checking
 * for multiple variable decalrations, and making sure return statements are
 * in a function body.
 *****************************************************************************/

type FunctionType int

const (
	ftNone FunctionType = iota
	ftFunction
)

type Resolver struct {
	interpreter         *Interpreter
	scopes              []map[string]bool
	currentFunctionType FunctionType
	errorHandler        *ErrorHandler
}

func NewResolver(interpreter *Interpreter) *Resolver {
	return &Resolver{interpreter: interpreter, scopes: make([]map[string]bool, 0, 0),
		currentFunctionType: ftNone, errorHandler: interpreter.errorHandler}
}

func (r *Resolver) ResolveStatements(statements []Stmt) {
	for _, stmt := range statements {
		r.resolveStatement(stmt)
	}
}

func (r *Resolver) resolveStatement(stmt Stmt) {
	stmt.accept(r)
}

func (r *Resolver) resolveExpression(expr Expr) {
	expr.accept(r)
}

func (r *Resolver) resolveFunction(function FunctionStmt, functionType FunctionType) {
	enclosingFunctionType := r.currentFunctionType
	r.currentFunctionType = functionType
	r.beginScope()
	for _, param := range function.params {
		r.declare(param)
		r.define(param)
	}
	r.ResolveStatements(function.body)
	r.endScope()
	r.currentFunctionType = enclosingFunctionType
}

func (r *Resolver) beginScope() {
	r.scopes = append(r.scopes, make(map[string]bool))
}

func (r *Resolver) endScope() {
	r.scopes = r.scopes[:len(r.scopes)-1]
}

func (r *Resolver) declare(name Token) {
	if len(r.scopes) == 0 {
		return
	}
	scope := r.scopes[0]
	_, hasVar := scope[name.lexeme]
	if hasVar {
		err := errors.New("Already a variable with this name in this scope.")
		r.errorHandler.reportStaticError(name.line, name.lexeme, err, false)
	}
	scope[name.lexeme] = false
}

func (r *Resolver) define(name Token) {
	if len(r.scopes) == 0 {
		return
	}
	r.scopes[0][name.lexeme] = true
}

func (r *Resolver) resolveLocal(expr Expr, name Token) {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		_, hasVar := r.scopes[i][name.lexeme]
		if hasVar {
			r.interpreter.resolve(expr, len(r.scopes)-1-i)
			return
		}
	}
}

func (r *Resolver) visitBlockStmt(stmt BlockStmt) any {
	r.beginScope()
	r.ResolveStatements(stmt.statements)
	r.endScope()
	return nil
}

func (r *Resolver) visitExprStmt(stmt ExprStmt) any {
	r.resolveExpression(stmt.expr)
	return nil
}

func (r *Resolver) visitFunctionStmt(stmt FunctionStmt) any {
	// declare and define immediately to allow self recursion
	r.declare(stmt.name)
	r.define(stmt.name)
	r.resolveFunction(stmt, ftFunction)
	return nil
}

func (r *Resolver) visitIfStmt(stmt IfStmt) any {
	// don't consider condition - check both branches regardless
	r.resolveExpression(stmt.condition)
	r.resolveStatement(stmt.thenBranch)
	if stmt.elseBranch != nil {
		r.resolveStatement(stmt.elseBranch)
	}
	return nil
}

func (r *Resolver) visitPrintStmt(stmt PrintStmt) any {
	r.resolveExpression(stmt.expr)
	return nil
}

func (r *Resolver) visitReturnStmt(stmt ReturnStmt) any {
	if r.currentFunctionType == ftNone {
		err := errors.New("Can't return from top-level code.")
		r.errorHandler.reportStaticError(stmt.keyword.line, stmt.keyword.lexeme, err, false)
	}
	if stmt.value != nil {
		r.resolveExpression(stmt.value)
	}
	return nil
}

func (r *Resolver) visitVarStmt(stmt VarStmt) any {
	r.declare(stmt.name)
	if stmt.initializer != nil {
		r.resolveExpression(stmt.initializer)
	}
	r.define(stmt.name)
	return nil
}

func (r *Resolver) visitWhileStmt(stmt WhileStmt) any {
	r.resolveExpression(stmt.condition)
	r.resolveStatement(stmt.body)
	return nil
}

func (r *Resolver) visitAssignExpr(expr AssignExpr) any {
	r.resolveExpression(expr.value)
	r.resolveLocal(expr, expr.name)
	return nil
}

func (r *Resolver) visitBinaryExpr(expr BinaryExpr) any {
	r.resolveExpression(expr.left)
	r.resolveExpression(expr.right)
	return nil
}

func (r *Resolver) visitCallExpr(expr CallExpr) any {
	r.resolveExpression(expr.callee)
	for _, arg := range expr.args {
		r.resolveExpression(arg)
	}
	return nil
}

func (r *Resolver) visitGroupingExpr(expr GroupingExpr) any {
	r.resolveExpression(expr.expression)
	return nil
}

func (r *Resolver) visitLiteralExpr(expr LiteralExpr) any {
	return nil
}

func (r *Resolver) visitLogicalExpr(expr LogicalExpr) any {
	r.resolveExpression(expr.left)
	r.resolveExpression(expr.right)
	return nil
}

func (r *Resolver) visitUnaryExpr(expr UnaryExpr) any {
	r.resolveExpression(expr.right)
	return nil
}

func (r *Resolver) visitVariableExpr(expr VariableExpr) any {
	if len(r.scopes) != 0 {
		varDefined, hasVar := r.scopes[0][expr.name.lexeme]
		if hasVar && !varDefined {
			err := errors.New("Can't read local variable in its own initializer.")
			r.errorHandler.reportStaticError(expr.name.line, expr.name.lexeme, err, false)
		}
	}
	r.resolveLocal(expr, expr.name)
	return nil
}
