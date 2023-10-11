package scanner

import (
	"errors"
	"strconv"
	"unicode"

	"github.com/skusel/glox/ast"
	"github.com/skusel/glox/langerr"
)

type Scanner struct {
	source       string
	tokens       []ast.Token
	start        int
	current      int
	line         int
	errorHandler *langerr.Handler
}

func NewScanner(source string, errorHandler *langerr.Handler) *Scanner {
	return &Scanner{source: source, start: 0, current: 0, line: 1, errorHandler: errorHandler}
}

func (s *Scanner) ScanTokens() []ast.Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}
	s.tokens = append(s.tokens, ast.Token{TokenType: ast.EndOfFile, Lexeme: "", Literal: nil, Line: s.line})
	return s.tokens
}

func (s *Scanner) addToken(t ast.TokenType) {
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
		s.errorHandler.Report(s.line, "", errors.New("Unterminated string."))
		return
	}

	s.advance() // The closing '"'

	// Trim the surrouding quotes
	value := s.source[s.start+1 : s.current-1]
	s.addGenericToken(ast.String, value)
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
		s.errorHandler.Report(s.line, "", errors.New("Invalid number."))
	} else {
		s.addGenericToken(ast.Number, value)
	}
}

func (s *Scanner) addIdentifierToken() {
	for unicode.IsDigit(rune(s.peek())) || unicode.IsLetter(rune(s.peek())) || s.peek() == '_' {
		s.advance()
	}

	text := s.source[s.start:s.current]
	if text == "and" {
		s.addGenericToken(ast.And, text)
	} else if text == "class" {
		s.addGenericToken(ast.Class, text)
	} else if text == "else" {
		s.addGenericToken(ast.Else, text)
	} else if text == "false" {
		s.addGenericToken(ast.False, text)
	} else if text == "for" {
		s.addGenericToken(ast.For, text)
	} else if text == "fun" {
		s.addGenericToken(ast.Fun, text)
	} else if text == "if" {
		s.addGenericToken(ast.If, text)
	} else if text == "nil" {
		s.addGenericToken(ast.Nil, text)
	} else if text == "or" {
		s.addGenericToken(ast.Or, text)
	} else if text == "print" {
		s.addGenericToken(ast.Print, text)
	} else if text == "return" {
		s.addGenericToken(ast.Return, text)
	} else if text == "super" {
		s.addGenericToken(ast.Super, text)
	} else if text == "this" {
		s.addGenericToken(ast.This, text)
	} else if text == "true" {
		s.addGenericToken(ast.True, text)
	} else if text == "var" {
		s.addGenericToken(ast.Var, text)
	} else if text == "while" {
		s.addGenericToken(ast.While, text)
	} else {
		s.addGenericToken(ast.Identifier, text)
	}
}

func (s *Scanner) addGenericToken(t ast.TokenType, literal any) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, ast.Token{TokenType: t, Lexeme: text, Literal: literal, Line: s.line})
}

func (s *Scanner) scanToken() {
	c := s.advance()
	if c == ' ' || c == '\r' || c == '\t' {
		return
	}
	switch c {
	case '(':
		s.addToken(ast.LeftParen)
	case ')':
		s.addToken(ast.RightParen)
	case '{':
		s.addToken(ast.LeftBrace)
	case '}':
		s.addToken(ast.RightBrace)
	case ',':
		s.addToken(ast.Comma)
	case '.':
		s.addToken(ast.Dot)
	case '-':
		s.addToken(ast.Minus)
	case '+':
		s.addToken(ast.Plus)
	case ';':
		s.addToken(ast.Semicolon)
	case '*':
		s.addToken(ast.Star)
	case '!':
		if s.match('=') {
			s.addToken(ast.BangEqual)
		} else {
			s.addToken(ast.Bang)
		}
	case '=':
		if s.match('=') {
			s.addToken(ast.EqualEqual)
		} else {
			s.addToken(ast.Equal)
		}
	case '<':
		if s.match('=') {
			s.addToken(ast.LessEqual)
		} else {
			s.addToken(ast.Less)
		}
	case '>':
		if s.match('=') {
			s.addToken(ast.GreaterEqual)
		} else {
			s.addToken(ast.Greater)
		}
	case '/':
		if s.match('/') {
			// A comment goes until the end of the line
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(ast.Slash)
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
			s.errorHandler.Report(s.line, "", errors.New("Unexpected character."))
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
