package lang

import (
	"errors"
	"os"
)

/******************************************************************************
 * The parser defines an abstract syntax tree (AST) given a sequence of tokens.
 * The AST is designed so that it is easy for the interpreter to consume.
 *
 * The parser implemented here is a recursive descent parser with a single
 * token of lookahead.
 *
 * Some popular parser generator tools found out in the wild are ANTLR and
 * Bison.
 *
 * Backus-Naur Form (BNF) of Parser Grammar
 * ========================================
 * program     -> statement* EOF ;
 * declaration -> classDecl
 *              | funDecl
 *              | varDecl
 *              | statement ;
 * statement   -> exprStmt
 *              | forStmt
 *              | ifStmt
 *              | printStmt
 *              | returnStmt
 *              | whileStmt
 *              | block ;
 * exprStmt    -> expression ";" ;
 * forStmt     -> "for" "(" ( varDecl | exprStmt | ";" )
 *                expression? ";"
 *                expression? ")" statement ;
 * classDecl   -> "class" IDENTIFIER ( "<" IDENTIFIER )? "{" function* "}" ;
 * funDecl     -> "fun" function ;
 * function    -> IDENTIFIER "(" parameters? ")" block ;
 * parameters  -> IDENTIFIER ( "," IDENTIFIER )* ;
 * ifStmt      -> "if" "(" expression ")" statement ( "else" statement )? ;
 * printStmt   -> "print" expression ";" ;
 * returnStmt  -> "return" expression? ";" ;
 * whileStmt   -> "while" "(" expression ")" statement ;
 * block       -> "{" + declaration* + "}" ;
 * varDecl     -> "var" IDENTIFIER ( "=" expression )? ";" ;
 * expression  -> assignment ;
 * assignment  -> ( call "." )? IDENTIFIER "=" assignment | logic_or ;
 * logic_or    -> logic_and ( "or" logic_and )* ;
 * logic_and   -> equality ( "and" equality )* ;
 * equality    -> comparison ( ("!=" | "==") comparision)* ;
 * comparison  -> term ( ( ">" | ">=" | "<" | "<=") term )* ;
 * term        -> factor ( ( "-" | "+" ) factor )* ;
 * factor      -> unary ( ( "/" | "*") unary )* ;
 * unary       -> ( "!" | "-" ) unary | call ;
 * call        -> primary ( "(" arguments? ")" | "." IDENTIFIER )* ;
 * arguments   -> expression ( "," expression )* ;
 * primary     -> "true" | "false" | "nil"
 *              | NUMBER | STRING
 *			    | "(" expression ")"
 *              | IDENTIFIER | "super" . IDENTIFIER ;
 *****************************************************************************/

type Parser struct {
	tokens       []Token
	current      int
	nextExprId   int
	errorHandler *ErrorHandler
}

func NewParser(tokens []Token, errorHandler *ErrorHandler) *Parser {
	return &Parser{tokens: tokens, current: 0, errorHandler: errorHandler}
}

func (p *Parser) Parse() []Stmt {
	statements := make([]Stmt, 0, 0)
	for !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}
	return statements
}

func (p *Parser) declaration() (stmt Stmt) {
	defer func() {
		/**********************************************************************
		 * Recover from a static error if one occurred. ErrorHandler "panics"
		 * when a static error requires synchronization. This causes the call
		 * stack to unwind and enter this defered function so that the parser
		 * can synchronize and continue reporting more errors before exiting
		 * with the static error exit code (65).
		 *********************************************************************/
		err := recover()
		if err != nil {
			staticError, isStaticError := err.(staticError)
			if isStaticError {
				os.Stderr.WriteString(staticError.msg)
				p.synchronize()
				stmt = nil
			} else {
				// this is not a panic thrown by us - pass it on
				panic(err)
			}
		}
	}()

	if p.match(tokenTypeClass) {
		stmt = p.classDeclaration()
	} else if p.match(tokenTypeFun) {
		stmt = p.function("function")
	} else if p.match(tokenTypeVar) {
		stmt = p.varDeclaration()
	} else {
		stmt = p.statement()
	}
	return stmt
}

func (p *Parser) classDeclaration() Stmt {
	name := p.consume(tokenTypeIdentifier, "Expect class name.")
	var superclass VariableExpr
	if p.match(tokenTypeLess) {
		p.consume(tokenTypeIdentifier, "Expect superclass name.")
		superclass = VariableExpr{id: p.getNextExprId(), name: p.previous()}
	}
	p.consume(tokenTypeLeftBrace, "Expect '{' before class body.")
	methods := make([]FunctionStmt, 0, 0)
	for !p.check(tokenTypeRightBrace) && !p.isAtEnd() {
		methods = append(methods, p.function("method"))
	}
	p.consume(tokenTypeRightBrace, "Expect '}' after class body.")
	return ClassStmt{name: name, superclass: superclass, methods: methods}
}

func (p *Parser) function(kind string) FunctionStmt {
	name := p.consume(tokenTypeIdentifier, "Expect "+kind+" name.")
	p.consume(tokenTypeLeftParen, "Expect '(' after "+kind+" name.")
	params := make([]Token, 0, 0)
	if !p.check(tokenTypeRightParen) {
		params = append(params, p.consume(tokenTypeIdentifier, "Expect parameter name."))
		for p.match(tokenTypeComma) {
			if len(params) >= 255 {
				p.createError(p.peek(), "Can't have more than 255 parameters.", false) // don't need to sync
			}
			params = append(params, p.consume(tokenTypeIdentifier, "Expect parameter name."))
		}
	}
	p.consume(tokenTypeRightParen, "Expect ')' after parameters.")
	// blockStatement expects '{' has already been matched
	p.consume(tokenTypeLeftBrace, "Expect '{' before "+kind+" body.")
	body := p.blockStatement()
	return FunctionStmt{name: name, params: params, body: body}
}

func (p *Parser) varDeclaration() Stmt {
	name := p.consume(tokenTypeIdentifier, "Expect variable name.")
	var initializer Expr
	if p.match(tokenTypeEqual) {
		initializer = p.expression()
	} else {
		initializer = nil
	}
	p.consume(tokenTypeSemicolon, "Expect ';' after variable declaration.")
	return VarStmt{name: name, initializer: initializer}
}

func (p *Parser) statement() Stmt {
	if p.match(tokenTypeFor) {
		return p.forStatement()
	} else if p.match(tokenTypeIf) {
		return p.ifStatement()
	} else if p.match(tokenTypePrint) {
		return p.printStatement()
	} else if p.match(tokenTypeReturn) {
		return p.returnStatement()
	} else if p.match(tokenTypeWhile) {
		return p.whileStatment()
	} else if p.match(tokenTypeLeftBrace) {
		return BlockStmt{statements: p.blockStatement()}
	} else {
		return p.expressionStatment()
	}
}

func (p *Parser) expressionStatment() Stmt {
	expr := p.expression()
	p.consume(tokenTypeSemicolon, "Expect ';' after expression.")
	return ExprStmt{expr: expr}
}

func (p *Parser) forStatement() Stmt {
	// desugar for statements into while statements
	p.consume(tokenTypeLeftParen, "Expect '(' after 'for'.")
	var initializer Stmt
	if p.match(tokenTypeSemicolon) {
		initializer = nil
	} else if p.match(tokenTypeVar) {
		initializer = p.varDeclaration()
	} else {
		initializer = p.expressionStatment()
	}
	var condition Expr
	if !p.check(tokenTypeSemicolon) {
		condition = p.expression()
	}
	p.consume(tokenTypeSemicolon, "Expect ';' after loop condition.")
	var increment Expr
	if !p.check(tokenTypeSemicolon) {
		increment = p.expression()
	}
	p.consume(tokenTypeRightParen, "Expect ')' after for clauses.")
	body := p.statement()
	if increment != nil {
		statements := []Stmt{body, ExprStmt{expr: increment}}
		body = BlockStmt{statements: statements}
	}
	if condition == nil {
		condition = LiteralExpr{id: p.getNextExprId(), value: true}
	}
	body = WhileStmt{condition: condition, body: body}
	if initializer != nil {
		statements := []Stmt{initializer, body}
		body = BlockStmt{statements: statements}
	}
	return body
}

func (p *Parser) ifStatement() Stmt {
	p.consume(tokenTypeLeftParen, "Expect '(' after 'if'.")
	condition := p.expression()
	p.consume(tokenTypeRightParen, "Expect ')' after if condition")
	thenBranch := p.statement()
	var elseBranch Stmt
	if p.match(tokenTypeElse) {
		elseBranch = p.statement()
	}
	return IfStmt{condition: condition, thenBranch: thenBranch, elseBranch: elseBranch}
}

func (p *Parser) printStatement() Stmt {
	value := p.expression()
	p.consume(tokenTypeSemicolon, "Expect ';' after value.")
	return PrintStmt{expr: value}
}

func (p *Parser) returnStatement() Stmt {
	keyword := p.previous()
	var value Expr
	if !p.check(tokenTypeSemicolon) {
		value = p.expression()
	}
	p.consume(tokenTypeSemicolon, "Expect ';' after return value.")
	return ReturnStmt{keyword: keyword, value: value}
}

func (p *Parser) whileStatment() Stmt {
	p.consume(tokenTypeLeftParen, "Expect '(' after 'while'.")
	condition := p.expression()
	p.consume(tokenTypeRightParen, "Expect ')' after while condition")
	body := p.statement()
	return WhileStmt{condition: condition, body: body}
}

func (p *Parser) blockStatement() []Stmt {
	statements := make([]Stmt, 0, 0)
	for !p.check(tokenTypeRightBrace) && !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}
	p.consume(tokenTypeRightBrace, "Expect '}' after block.")
	return statements
}

func (p *Parser) expression() Expr {
	return p.assignment()
}

func (p *Parser) assignment() Expr {
	expr := p.or()
	if p.match(tokenTypeEqual) {
		equals := p.previous()
		value := p.assignment()

		variableExpr, isVariableExpr := expr.(VariableExpr)
		if isVariableExpr {
			return AssignExpr{id: p.getNextExprId(), name: variableExpr.name, value: value}
		}
		getExpr, isGetExpr := expr.(GetExpr)
		if isGetExpr {
			return SetExpr{id: p.getNextExprId(), object: getExpr.object, name: getExpr.name, value: value}
		}
		p.createError(equals, "Invalid assignment target.", false) // don't need to sync
	}
	return expr
}

func (p *Parser) or() Expr {
	expr := p.and()
	for p.match(tokenTypeOr) {
		operator := p.previous()
		right := p.and()
		expr = LogicalExpr{id: p.getNextExprId(), left: expr, operator: operator, right: right}
	}
	return expr
}

func (p *Parser) and() Expr {
	expr := p.equality()
	for p.match(tokenTypeAnd) {
		operator := p.previous()
		right := p.equality()
		expr = LogicalExpr{id: p.getNextExprId(), left: expr, operator: operator, right: right}
	}
	return expr
}

func (p *Parser) equality() Expr {
	expr := p.comparison()
	for p.match(tokenTypeBangEqual, tokenTypeEqualEqual) {
		operator := p.previous()
		right := p.comparison()
		expr = BinaryExpr{id: p.getNextExprId(), left: expr, operator: operator, right: right}
	}
	return expr
}

func (p *Parser) comparison() Expr {
	expr := p.term()
	for p.match(tokenTypeGreater, tokenTypeGreaterEqual, tokenTypeLess, tokenTypeLessEqual) {
		operator := p.previous()
		right := p.term()
		expr = BinaryExpr{id: p.getNextExprId(), left: expr, operator: operator, right: right}
	}
	return expr
}

func (p *Parser) term() Expr {
	expr := p.factor()
	for p.match(tokenTypeMinus, tokenTypePlus) {
		operator := p.previous()
		right := p.factor()
		expr = BinaryExpr{id: p.getNextExprId(), left: expr, operator: operator, right: right}
	}
	return expr
}

func (p *Parser) factor() Expr {
	expr := p.unary()
	for p.match(tokenTypeSlash, tokenTypeStar, tokenTypeMod) {
		operator := p.previous()
		right := p.unary()
		expr = BinaryExpr{id: p.getNextExprId(), left: expr, operator: operator, right: right}
	}
	return expr
}

func (p *Parser) unary() Expr {
	if p.match(tokenTypeBang, tokenTypeMinus) {
		operator := p.previous()
		right := p.primary()
		return UnaryExpr{id: p.getNextExprId(), operator: operator, right: right}
	}
	return p.call()
}

func (p *Parser) call() Expr {
	expr := p.primary()

	for {
		if p.match(tokenTypeLeftParen) {
			expr = p.finishCall(expr)
		} else if p.match(tokenTypeDot) {
			name := p.consume(tokenTypeIdentifier, "Expect property name after '.'.")
			expr = GetExpr{id: p.getNextExprId(), object: expr, name: name}
		} else {
			break
		}
	}

	return expr
}

func (p *Parser) finishCall(callee Expr) Expr {
	args := make([]Expr, 0, 0)
	if !p.check(tokenTypeRightParen) {
		args = append(args, p.expression())
		for p.match(tokenTypeComma) {
			if len(args) >= 255 {
				p.createError(p.peek(), "Can't have more than 255 arguments.", false) // don't need to sync
			}
			args = append(args, p.expression())
		}
	}
	paren := p.consume(tokenTypeRightParen, "Expect ')' after arguments.")
	return CallExpr{id: p.getNextExprId(), callee: callee, paren: paren, args: args}
}

func (p *Parser) primary() Expr {
	if p.match(tokenTypeFalse) {
		return LiteralExpr{id: p.getNextExprId(), value: false}
	} else if p.match(tokenTypeTrue) {
		return LiteralExpr{id: p.getNextExprId(), value: true}
	} else if p.match(tokenTypeNil) {
		return LiteralExpr{id: p.getNextExprId(), value: nil}
	} else if p.match(tokenTypeNumber, tokenTypeString) {
		return LiteralExpr{id: p.getNextExprId(), value: p.previous().literal}
	} else if p.match(tokenTypeSuper) {
		keyword := p.previous()
		p.consume(tokenTypeDot, "Expect '.' after 'super'.")
		method := p.consume(tokenTypeIdentifier, "Expect superclass method name.")
		return SuperExpr{id: p.getNextExprId(), keyword: keyword, method: method}
	} else if p.match(tokenTypeThis) {
		return ThisExpr{id: p.getNextExprId(), keyword: p.previous()}
	} else if p.match(tokenTypeIdentifier) {
		return VariableExpr{id: p.getNextExprId(), name: p.previous()}
	} else if p.match(tokenTypeLeftParen) {
		expr := p.expression()
		p.consume(tokenTypeRightParen, "Expect ')' after expression.")
		return GroupingExpr{id: p.getNextExprId(), expression: expr}
	}
	p.createError(p.peek(), "Expect expression.", true)
	return nil
}

func (p *Parser) match(tokenTypes ...TokenType) bool {
	for _, tokenType := range tokenTypes {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) consume(tokenType TokenType, msg string) Token {
	if p.check(tokenType) {
		return p.advance()
	}
	p.createError(p.peek(), msg, true)
	return p.peek()
}

func (p *Parser) check(tokenType TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().tokenType == tokenType
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().tokenType == tokenTypeEndOfFile
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

func (p *Parser) getNextExprId() int {
	p.nextExprId++
	return p.nextExprId
}

func (p *Parser) createError(token Token, msg string, synchronize bool) {
	p.errorHandler.reportStaticError(token.line, token.lexeme, errors.New(msg), synchronize)
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().tokenType == tokenTypeSemicolon {
			return
		}

		switch p.peek().tokenType {
		case tokenTypeClass:
			fallthrough
		case tokenTypeFor:
			fallthrough
		case tokenTypeFun:
			fallthrough
		case tokenTypeIf:
			fallthrough
		case tokenTypePrint:
			fallthrough
		case tokenTypeReturn:
			fallthrough
		case tokenTypeVar:
			fallthrough
		case tokenTypeWhile:
			return
		}

		p.advance()
	}
}
