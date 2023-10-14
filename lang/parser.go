package lang

import (
	"errors"
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
 * declaration -> varDecl
 *              | statement ;
 * statement   -> exprStmt
 *              | printStmt
 *              | block ;
 * exprStmt    -> expression ";" ;
 * block       -> "{" + declaration* + "}" ;
 * printStmt   -> "print" expression ";" ;
 * varDecl     -> "var" IDENTIFIER ( "=" expression )? ";" ;
 * expression  -> assignment ;
 * assignment  -> IDENTIFIER "=" assignment
 *              | equality ;
 * equality    -> comparison ( ("!=" | "==") comparision)* ;
 * comparison  -> term ( ( ">" | ">=" | "<" | "<=") term )* ;
 * term        -> factor ( ( "-" | "+" ) factor )* ;
 * factor      -> unary ( ( "/" | "*") unary )* ;
 * unary       -> ( "!" | "-" ) unary
 *              | primary ;
 * primary     -> "true" | "false" | "nil"
 *              | NUMBER | STRING
 *			    | "(" expression ")"
 *              | IDENTIFIER ;
 *****************************************************************************/

type Parser struct {
	tokens       []Token
	current      int
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

func (p *Parser) declaration() Stmt {
	var stmt Stmt
	if p.match(tokenTypeVar) {
		stmt = p.varDeclaration()
	} else {
		stmt = p.statement()
	}
	if p.errorHandler.needToSynchronize {
		p.synchronize()
		p.errorHandler.needToSynchronize = false
		return nil
	} else {
		return stmt
	}
}

func (p *Parser) varDeclaration() Stmt {
	name := p.consume(tokenTypeIdentifier, "Expect variable name.")
	if p.errorHandler.needToSynchronize {
		return nil
	}
	var initializer Expr
	if p.match(tokenTypeEqual) {
		initializer = p.expression()
		if p.errorHandler.needToSynchronize {
			return nil
		}
	} else {
		initializer = nil
	}
	p.consume(tokenTypeSemicolon, "Expect ';' after variable declaration.")
	if p.errorHandler.needToSynchronize {
		return nil
	}
	return VarStmt{name: name, initializer: initializer}
}

func (p *Parser) statement() Stmt {
	if p.match(tokenTypePrint) {
		return p.printStatement()
	} else if p.match(tokenTypeLeftBrace) {
		return BlockStmt{statements: p.blockStatement()}
	}
	return p.expressionStatment()
}

func (p *Parser) expressionStatment() Stmt {
	expr := p.expression()
	if p.errorHandler.needToSynchronize {
		return nil
	}
	p.consume(tokenTypeSemicolon, "Expect ';' after expression.")
	if p.errorHandler.needToSynchronize {
		return nil
	}
	return ExprStmt{expr: expr}
}

func (p *Parser) blockStatement() []Stmt {
	statements := make([]Stmt, 0, 0)
	for !p.check(tokenTypeRightBrace) && !p.isAtEnd() && !p.errorHandler.needToSynchronize {
		statements = append(statements, p.declaration())
	}
	p.consume(tokenTypeRightBrace, "Expect '}' after block.")
	if p.errorHandler.needToSynchronize {
		return nil
	}
	return statements
}

func (p *Parser) printStatement() Stmt {
	value := p.expression()
	if p.errorHandler.needToSynchronize {
		return nil
	}
	p.consume(tokenTypeSemicolon, "Expect ';' after value.")
	if p.errorHandler.needToSynchronize {
		return nil
	}
	return PrintStmt{expr: value}
}

func (p *Parser) expression() Expr {
	if p.errorHandler.needToSynchronize {
		return nil
	}
	return p.assignment()
}

func (p *Parser) assignment() Expr {
	expr := p.equality()
	if p.match(tokenTypeEqual) && !p.errorHandler.needToSynchronize {
		equals := p.previous()
		value := p.assignment()

		variableExpr, isVariableExpr := expr.(VariableExpr)
		if isVariableExpr {
			name := variableExpr.name
			return AssignExpr{name: name, value: value}
		}
		p.errorHandler.report(equals.line, "", errors.New("Invalid assignment target."))
	}
	return expr
}

func (p *Parser) equality() Expr {
	expr := p.comparison()
	for p.match(tokenTypeBangEqual, tokenTypeEqualEqual) && !p.errorHandler.needToSynchronize {
		operator := p.previous()
		right := p.comparison()
		expr = BinaryExpr{left: expr, operator: operator, right: right}
	}
	return expr
}

func (p *Parser) comparison() Expr {
	expr := p.term()
	for p.match(tokenTypeGreater, tokenTypeGreaterEqual, tokenTypeLess, tokenTypeLessEqual) && !p.errorHandler.needToSynchronize {
		operator := p.previous()
		right := p.term()
		expr = BinaryExpr{left: expr, operator: operator, right: right}
	}
	return expr
}

func (p *Parser) term() Expr {
	expr := p.factor()
	for p.match(tokenTypeMinus, tokenTypePlus) && !p.errorHandler.needToSynchronize {
		operator := p.previous()
		right := p.factor()
		expr = BinaryExpr{left: expr, operator: operator, right: right}
	}
	return expr
}

func (p *Parser) factor() Expr {
	expr := p.unary()
	for p.match(tokenTypeSlash, tokenTypeStar) && !p.errorHandler.needToSynchronize {
		operator := p.previous()
		right := p.unary()
		expr = BinaryExpr{left: expr, operator: operator, right: right}
	}
	return expr
}

func (p *Parser) unary() Expr {
	if p.match(tokenTypeBang, tokenTypeMinus) && !p.errorHandler.needToSynchronize {
		operator := p.previous()
		right := p.primary()
		return UnaryExpr{operator: operator, right: right}
	}
	return p.primary()
}

func (p *Parser) primary() Expr {
	if p.errorHandler.needToSynchronize {
		return nil
	}
	if p.match(tokenTypeFalse) {
		return LiteralExpr{value: false}
	} else if p.match(tokenTypeTrue) {
		return LiteralExpr{value: true}
	} else if p.match(tokenTypeNil) {
		return LiteralExpr{value: nil}
	} else if p.match(tokenTypeNumber, tokenTypeString) {
		return LiteralExpr{value: p.previous().literal}
	} else if p.match(tokenTypeIdentifier) {
		return VariableExpr{name: p.previous()}
	} else if p.match(tokenTypeLeftParen) {
		expr := p.expression()
		p.consume(tokenTypeRightParen, "Expect ')' after expression.")
		return GroupingExpr{expression: expr}
	}
	p.createError("Expect expression.")
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
	p.createError(msg)
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

func (p *Parser) createError(msg string) {
	currentToken := p.peek()
	p.errorHandler.report(currentToken.line, currentToken.lexeme, errors.New(msg))
}

func (p *Parser) synchronize() {
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
