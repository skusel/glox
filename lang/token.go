package lang

import "fmt"

type TokenType int

/******************************************************************************
 * A token represents individual lexemes in the Lox source code.
 *****************************************************************************/

const (
	// single char tokens
	tokenTypeLeftParen TokenType = iota
	tokenTypeRightParen
	tokenTypeLeftBrace
	tokenTypeRightBrace
	tokenTypeComma
	tokenTypeDot
	tokenTypeMinus
	tokenTypePlus
	tokenTypeSemicolon
	tokenTypeSlash
	tokenTypeStar
	// comparison operator tokens
	tokenTypeBang
	tokenTypeBangEqual
	tokenTypeEqual
	tokenTypeEqualEqual
	tokenTypeGreater
	tokenTypeGreaterEqual
	tokenTypeLess
	tokenTypeLessEqual
	// literals
	tokenTypeIdentifier
	tokenTypeString
	tokenTypeNumber
	// keywords
	tokenTypeAnd
	tokenTypeClass
	tokenTypeElse
	tokenTypeFalse
	tokenTypeFun
	tokenTypeFor
	tokenTypeIf
	tokenTypeNil
	tokenTypeOr
	tokenTypePrint
	tokenTypeReturn
	tokenTypeSuper
	tokenTypeThis
	tokenTypeTrue
	tokenTypeVar
	tokenTypeWhile
	// end of file
	tokenTypeEndOfFile
)

type Token struct {
	tokenType TokenType
	lexeme    string
	literal   any
	line      int
}

func (t Token) ToString() string {
	return fmt.Sprintf("%d %s %s", t.tokenType, t.lexeme, t.literal)
}
