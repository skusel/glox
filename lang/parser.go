package lang

import (
	"errors"
)

/******************************************************************************
 * BNF Grammar
 * ===========
 * expression -> equality ;
 * equality   -> comparison ( ("!=" | "==") comparision)* ;
 * comparison -> term ( ( ">" | ">=" | "<" | "<=") term )* ;
 * term       -> factor ( ( "-" | "+" ) factor )* ;
 * factor     -> unary ( ( "/" | "*") unary )* ;
 * unary      -> ( "!" | "-" ) unary
 *             | primary ;
 * primary    -> NUMBER | STRING | "true" | "false" | "nil"
 *			   | "(" expression ")" ;
 *****************************************************************************/

type Parser struct {
	tokens       []Token
	current      int
	errorHandler *ErrorHandler
}

func NewParser(tokens []Token, errorHandler *ErrorHandler) *Parser {
	return &Parser{tokens: tokens, current: 0, errorHandler: errorHandler}
}

func (p *Parser) Parse() Expr {
	expr := p.expression()
	if !p.errorHandler.HadError {
		return expr
	} else {
		return nil
	}
}

func (p *Parser) expression() Expr {
	return p.equality()
}

func (p *Parser) equality() Expr {
	expr := p.comparison()
	for p.match(tokenTypeBangEqual, tokenTypeEqualEqual) && !p.errorHandler.HadError {
		operator := p.previous()
		right := p.comparison()
		expr = BinaryExpr{left: expr, operator: operator, right: right}
	}
	return expr
}

func (p *Parser) comparison() Expr {
	expr := p.term()
	for p.match(tokenTypeGreater, tokenTypeGreaterEqual, tokenTypeLess, tokenTypeLessEqual) && !p.errorHandler.HadError {
		operator := p.previous()
		right := p.term()
		expr = BinaryExpr{left: expr, operator: operator, right: right}
	}
	return expr
}

func (p *Parser) term() Expr {
	expr := p.factor()
	for p.match(tokenTypeMinus, tokenTypePlus) && !p.errorHandler.HadError {
		operator := p.previous()
		right := p.factor()
		expr = BinaryExpr{left: expr, operator: operator, right: right}
	}
	return expr
}

func (p *Parser) factor() Expr {
	expr := p.unary()
	for p.match(tokenTypeSlash, tokenTypeStar) && !p.errorHandler.HadError {
		operator := p.previous()
		right := p.unary()
		expr = BinaryExpr{left: expr, operator: operator, right: right}
	}
	return expr
}

func (p *Parser) unary() Expr {
	if p.match(tokenTypeBang, tokenTypeMinus) {
		operator := p.previous()
		right := p.primary()
		return UnaryExpr{operator: operator, right: right}
	}
	return p.primary()
}

func (p *Parser) primary() Expr {
	if p.match(tokenTypeFalse) {
		return LiteralExpr{value: false}
	} else if p.match(tokenTypeTrue) {
		return LiteralExpr{value: true}
	} else if p.match(tokenTypeNil) {
		return LiteralExpr{value: nil}
	} else if p.match(tokenTypeNumber, tokenTypeString) {
		return LiteralExpr{value: p.previous().literal}
	} else if p.match(tokenTypeLeftParen) {
		expr := p.expression()
		if p.errorHandler.HadError {
			return expr
		} else {
			p.consume(tokenTypeRightParen, "Expect ')' after expression.")
			return GroupingExpr{expression: expr}
		}
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
