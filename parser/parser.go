package parser

import (
	"errors"

	"github.com/skusel/glox/ast"
	"github.com/skusel/glox/langerr"
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
	errorHandler *langerr.Handler
}

func NewParser(tokens []ast.Token, errorHandler *langerr.Handler) *Parser {
	return &Parser{tokens: tokens, current: 0, errorHandler: errorHandler}
}

func (p *Parser) Parse() ast.Expr {
	expr, err := p.expression()
	if err == nil {
		return expr
	} else {
		return nil
	}
}

func (p *Parser) expression() (ast.Expr, error) {
	return p.equality()
}

func (p *Parser) equality() (ast.Expr, error) {
	expr, err := p.comparison()
	for p.match(ast.BangEqual, ast.EqualEqual) && err == nil {
		operator := p.previous()
		right, lastErr := p.comparison()
		expr = ast.BinaryExpr{Left: expr, Operator: operator, Right: right}
		err = lastErr
	}
	return expr, err
}

func (p *Parser) comparison() (ast.Expr, error) {
	expr, err := p.term()
	for p.match(ast.Greater, ast.GreaterEqual, ast.Less, ast.LessEqual) && err == nil {
		operator := p.previous()
		right, lastErr := p.term()
		expr = ast.BinaryExpr{Left: expr, Operator: operator, Right: right}
		err = lastErr
	}
	return expr, err
}

func (p *Parser) term() (ast.Expr, error) {
	expr, err := p.factor()
	for p.match(ast.Minus, ast.Plus) && err == nil {
		operator := p.previous()
		right, lastErr := p.factor()
		expr = ast.BinaryExpr{Left: expr, Operator: operator, Right: right}
		err = lastErr
	}
	return expr, err
}

func (p *Parser) factor() (ast.Expr, error) {
	expr, err := p.unary()
	for p.match(ast.Slash, ast.Star) && err == nil {
		operator := p.previous()
		right, lastErr := p.unary()
		expr = ast.BinaryExpr{Left: expr, Operator: operator, Right: right}
		err = lastErr
	}
	return expr, err
}

func (p *Parser) unary() (ast.Expr, error) {
	if p.match(ast.Bang, ast.Minus) {
		operator := p.previous()
		right, err := p.primary()
		return ast.UnaryExpr{Operator: operator, Right: right}, err
	}
	return p.primary()
}

func (p *Parser) primary() (ast.Expr, error) {
	if p.match(ast.False) {
		return ast.LiteralExpr{Value: false}, nil
	} else if p.match(ast.True) {
		return ast.LiteralExpr{Value: true}, nil
	} else if p.match(ast.Nil) {
		return ast.LiteralExpr{Value: nil}, nil
	} else if p.match(ast.Number, ast.String) {
		return ast.LiteralExpr{Value: p.previous().Literal}, nil
	} else if p.match(ast.LeftParen) {
		expr, innerErr := p.expression()
		if innerErr != nil {
			return expr, innerErr
		} else {
			_, outterErr := p.consume(ast.RightParen, "Expect ')' after expression.")
			return ast.GroupingExpr{Expression: expr}, outterErr
		}
	}
	return nil, p.createError("Expect expression.")
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

func (p *Parser) consume(t ast.TokenType, msg string) (ast.Token, error) {
	if p.check(t) {
		return p.advance(), nil
	}
	return p.peek(), p.createError(msg)
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

func (p *Parser) createError(msg string) error {
	currentToken := p.peek()
	p.errorHandler.Report(currentToken.Line, currentToken.Lexeme, msg)
	return errors.New(msg)
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
