package lang

import (
	"errors"
	"strconv"
	"unicode"
)

/******************************************************************************
 * The scanner takes in source code an transforms it into a list of tokens.
 * We are also able to catch a few static errors in the scanner. For example,
 * unterminated strings.
 *****************************************************************************/

type Scanner struct {
	source       string
	tokens       []Token
	start        int
	current      int
	line         int
	errorHandler *ErrorHandler
}

func NewScanner(source string, errorHandler *ErrorHandler) *Scanner {
	return &Scanner{source: source, start: 0, current: 0, line: 1, errorHandler: errorHandler}
}

func (s *Scanner) ScanTokens() []Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}
	s.tokens = append(s.tokens, Token{tokenType: tokenTypeEndOfFile, lexeme: "", literal: nil, line: s.line})
	return s.tokens
}

func (s *Scanner) addToken(t TokenType) {
	s.addGenericToken(t, nil)
}

func (s *Scanner) addStringToken() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		s.errorHandler.report(s.line, "", errors.New("Unterminated string."))
		return
	}

	s.advance() // The closing '"'

	// Trim the surrouding quotes
	value := s.source[s.start+1 : s.current-1]
	s.addGenericToken(tokenTypeString, value)
}

func (s *Scanner) addNumberToken() {
	for unicode.IsDigit(rune(s.peek())) {
		s.advance()
	}

	if s.peek() == '.' && unicode.IsDigit(rune(s.peekNext())) {
		s.advance()

		for unicode.IsDigit(rune(s.peek())) {
			s.advance()
		}
	}

	value, err := strconv.ParseFloat(s.source[s.start:s.current], 64)
	if err != nil {
		s.errorHandler.report(s.line, "", errors.New("Invalid number."))
	} else {
		s.addGenericToken(tokenTypeNumber, value)
	}
}

func (s *Scanner) addIdentifierToken() {
	for unicode.IsDigit(rune(s.peek())) || unicode.IsLetter(rune(s.peek())) || s.peek() == '_' {
		s.advance()
	}

	text := s.source[s.start:s.current]
	if text == "and" {
		s.addGenericToken(tokenTypeAnd, text)
	} else if text == "class" {
		s.addGenericToken(tokenTypeClass, text)
	} else if text == "else" {
		s.addGenericToken(tokenTypeElse, text)
	} else if text == "false" {
		s.addGenericToken(tokenTypeFalse, text)
	} else if text == "for" {
		s.addGenericToken(tokenTypeFor, text)
	} else if text == "fun" {
		s.addGenericToken(tokenTypeFun, text)
	} else if text == "if" {
		s.addGenericToken(tokenTypeIf, text)
	} else if text == "nil" {
		s.addGenericToken(tokenTypeNil, text)
	} else if text == "or" {
		s.addGenericToken(tokenTypeOr, text)
	} else if text == "print" {
		s.addGenericToken(tokenTypePrint, text)
	} else if text == "return" {
		s.addGenericToken(tokenTypeReturn, text)
	} else if text == "super" {
		s.addGenericToken(tokenTypeSuper, text)
	} else if text == "this" {
		s.addGenericToken(tokenTypeThis, text)
	} else if text == "true" {
		s.addGenericToken(tokenTypeTrue, text)
	} else if text == "var" {
		s.addGenericToken(tokenTypeVar, text)
	} else if text == "while" {
		s.addGenericToken(tokenTypeWhile, text)
	} else {
		s.addGenericToken(tokenTypeIdentifier, text)
	}
}

func (s *Scanner) addGenericToken(tokenType TokenType, literal any) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, Token{tokenType: tokenType, lexeme: text, literal: literal, line: s.line})
}

func (s *Scanner) scanToken() {
	c := s.advance()
	if c == ' ' || c == '\r' || c == '\t' {
		return
	}
	switch c {
	case '(':
		s.addToken(tokenTypeLeftParen)
	case ')':
		s.addToken(tokenTypeRightParen)
	case '{':
		s.addToken(tokenTypeLeftBrace)
	case '}':
		s.addToken(tokenTypeRightBrace)
	case ',':
		s.addToken(tokenTypeComma)
	case '.':
		s.addToken(tokenTypeDot)
	case '-':
		s.addToken(tokenTypeMinus)
	case '+':
		s.addToken(tokenTypePlus)
	case ';':
		s.addToken(tokenTypeSemicolon)
	case '*':
		s.addToken(tokenTypeStar)
	case '!':
		if s.match('=') {
			s.addToken(tokenTypeBangEqual)
		} else {
			s.addToken(tokenTypeBang)
		}
	case '=':
		if s.match('=') {
			s.addToken(tokenTypeEqualEqual)
		} else {
			s.addToken(tokenTypeEqual)
		}
	case '<':
		if s.match('=') {
			s.addToken(tokenTypeLessEqual)
		} else {
			s.addToken(tokenTypeLess)
		}
	case '>':
		if s.match('=') {
			s.addToken(tokenTypeGreaterEqual)
		} else {
			s.addToken(tokenTypeGreater)
		}
	case '/':
		if s.match('/') {
			// A comment goes until the end of the line
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(tokenTypeSlash)
		}
	case '\n':
		s.line++
	case '"':
		s.addStringToken()
	default:
		if unicode.IsDigit(rune(c)) {
			s.addNumberToken()
		} else if unicode.IsLetter(rune(c)) || c == '_' {
			s.addIdentifierToken()
		} else {
			s.errorHandler.report(s.line, "", errors.New("Unexpected character."))
		}
	}
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) advance() byte {
	nextC := s.source[s.current]
	s.current++
	return nextC
}

func (s *Scanner) match(expected byte) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != expected {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) peek() byte {
	// this scanner has a lookahead of 1
	if s.isAtEnd() {
		return 0
	} else {
		return s.source[s.current]
	}
}

func (s *Scanner) peekNext() byte {
	if s.current+1 >= len(s.source) {
		return 0
	} else {
		return s.source[s.current+1]
	}
}
