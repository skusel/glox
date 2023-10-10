package ast

import "fmt"

type TokenType int

const (
	// single char tokens
	LeftParen TokenType = iota
	RightParen
	LeftBrace
	RightBrace
	Comma
	Dot
	Minus
	Plus
	Semicolon
	Slash
	Star
	// comparison operator tokens
	Bang
	BangEqual
	Equal
	EqualEqual
	Greater
	GreaterEqual
	Less
	LessEqual
	// literals
	Identifier
	String
	Number
	// keywords
	And
	Class
	Else
	False
	Fun
	For
	If
	Nil
	Or
	Print
	Return
	Super
	This
	True
	Var
	While
	// end of file
	EndOfFile
)

type Token struct {
	TokenType TokenType
	Lexeme    string
	Literal   any
	Line      int
}

func (t Token) ToString() string {
	return fmt.Sprintf("%d %s %s", t.TokenType, t.Lexeme, t.Literal)
}
