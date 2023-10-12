package parser

import (
	"errors"

	"github.com/skusel/glox/ast"
	"github.com/skusel/glox/gloxerror"
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
	tokens       []ast.Token
	current      int
	errorHandler *gloxerror.Handler
}

func NewParser(tokens []ast.Token, errorHandler *gloxerror.Handler) *Parser {
	return &Parser{tokens: tokens, current: 0, errorHandler: errorHandler}
}

func (p *Parser) Parse() ast.Expr {
	expr := p.expression()
	if !p.errorHandler.HadError {
		return expr
	} else {
		return nil
	}
}

func (p *Parser) expression() ast.Expr {
	return p.equality()
}

func (p *Parser) equality() ast.Expr {
	expr := p.comparison()
	for p.match(ast.BangEqual, ast.EqualEqual) && !p.errorHandler.HadError {
		operator := p.previous()
		right := p.comparison()
		expr = ast.BinaryExpr{Left: expr, Operator: operator, Right: right}
	}
	return expr
}

func (p *Parser) comparison() ast.Expr {
	expr := p.term()
	for p.match(ast.Greater, ast.GreaterEqual, ast.Less, ast.LessEqual) && !p.errorHandler.HadError {
		operator := p.previous()
		right := p.term()
		expr = ast.BinaryExpr{Left: expr, Operator: operator, Right: right}
	}
	return expr
}

func (p *Parser) term() ast.Expr {
	expr := p.factor()
	for p.match(ast.Minus, ast.Plus) && !p.errorHandler.HadError {
		operator := p.previous()
		right := p.factor()
		expr = ast.BinaryExpr{Left: expr, Operator: operator, Right: right}
	}
	return expr
}

func (p *Parser) factor() ast.Expr {
	expr := p.unary()
	for p.match(ast.Slash, ast.Star) && !p.errorHandler.HadError {
		operator := p.previous()
		right := p.unary()
		expr = ast.BinaryExpr{Left: expr, Operator: operator, Right: right}
	}
	return expr
}

func (p *Parser) unary() ast.Expr {
	if p.match(ast.Bang, ast.Minus) {
		operator := p.previous()
		right := p.primary()
		return ast.UnaryExpr{Operator: operator, Right: right}
	}
	return p.primary()
}

func (p *Parser) primary() ast.Expr {
	if p.match(ast.False) {
		return ast.LiteralExpr{Value: false}
	} else if p.match(ast.True) {
		return ast.LiteralExpr{Value: true}
	} else if p.match(ast.Nil) {
		return ast.LiteralExpr{Value: nil}
	} else if p.match(ast.Number, ast.String) {
		return ast.LiteralExpr{Value: p.previous().Literal}
	} else if p.match(ast.LeftParen) {
		expr := p.expression()
		if p.errorHandler.HadError {
			return expr
		} else {
			p.consume(ast.RightParen, "Expect ')' after expression.")
			return ast.GroupingExpr{Expression: expr}
		}
	}
	p.createError("Expect expression.")
	return nil
}

func (p *Parser) match(types ...ast.TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) consume(t ast.TokenType, msg string) ast.Token {
	if p.check(t) {
		return p.advance()
	}
	p.createError(msg)
	return p.peek()
}

func (p *Parser) check(t ast.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().TokenType == t
}

func (p *Parser) advance() ast.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().TokenType == ast.EndOfFile
}

func (p *Parser) peek() ast.Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() ast.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) createError(msg string) {
	currentToken := p.peek()
	p.errorHandler.Report(currentToken.Line, currentToken.Lexeme, errors.New(msg))
}

func (p *Parser) synchronize() {
	for !p.isAtEnd() {
		if p.previous().TokenType == ast.Semicolon {
			return
		}

		switch p.peek().TokenType {
		case ast.Class:
			fallthrough
		case ast.For:
			fallthrough
		case ast.Fun:
			fallthrough
		case ast.If:
			fallthrough
		case ast.Print:
			fallthrough
		case ast.Return:
			fallthrough
		case ast.Var:
			fallthrough
		case ast.While:
			return
		}

		p.advance()
	}
}
