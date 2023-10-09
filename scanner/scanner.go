package scanner

import (
	"strconv"
	"unicode"

	"github.com/skusel/glox/error"
	"github.com/skusel/glox/gloxtoken"
)

type Scanner struct {
	source       string
	tokens       []gloxtoken.Token
	start        int
	current      int
	line         int
	errorHandler *error.Handler
}

func NewScanner(source string, errorHandler *error.Handler) *Scanner {
	return &Scanner{source: source, start: 0, current: 0, line: 1, errorHandler: errorHandler}
}

func (s *Scanner) ScanTokens() []gloxtoken.Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}
	s.tokens = append(s.tokens, gloxtoken.Token{TokenType: gloxtoken.EndOfFile, Lexeme: "", Literal: nil, Line: s.line})
	return s.tokens
}

func (s *Scanner) addToken(token gloxtoken.TokenType) {
	s.addGenericToken(token, nil)
}

func (s *Scanner) addStringToken() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		s.errorHandler.Report(s.line, "", "Unterminated string.")
		return
	}

	s.advance() // The closing '"'

	// Trim the surrouding quotes
	value := s.source[s.start+1 : s.current-1]
	s.addGenericToken(gloxtoken.String, value)
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
		s.errorHandler.Report(s.line, "", "Invalid number.")
	} else {
		s.addGenericToken(gloxtoken.Number, value)
	}
}

func (s *Scanner) addIdentifierToken() {
	for unicode.IsDigit(rune(s.peek())) || unicode.IsLetter(rune(s.peek())) || s.peek() == '_' {
		s.advance()
	}

	text := s.source[s.start:s.current]
	if text == "and" {
		s.addGenericToken(gloxtoken.And, text)
	} else if text == "class" {
		s.addGenericToken(gloxtoken.Class, text)
	} else if text == "else" {
		s.addGenericToken(gloxtoken.Else, text)
	} else if text == "false" {
		s.addGenericToken(gloxtoken.False, text)
	} else if text == "for" {
		s.addGenericToken(gloxtoken.For, text)
	} else if text == "fun" {
		s.addGenericToken(gloxtoken.Fun, text)
	} else if text == "if" {
		s.addGenericToken(gloxtoken.If, text)
	} else if text == "nil" {
		s.addGenericToken(gloxtoken.Nil, text)
	} else if text == "or" {
		s.addGenericToken(gloxtoken.Or, text)
	} else if text == "print" {
		s.addGenericToken(gloxtoken.Print, text)
	} else if text == "return" {
		s.addGenericToken(gloxtoken.Return, text)
	} else if text == "super" {
		s.addGenericToken(gloxtoken.Super, text)
	} else if text == "this" {
		s.addGenericToken(gloxtoken.This, text)
	} else if text == "true" {
		s.addGenericToken(gloxtoken.True, text)
	} else if text == "var" {
		s.addGenericToken(gloxtoken.Var, text)
	} else if text == "while" {
		s.addGenericToken(gloxtoken.While, text)
	} else {
		s.addGenericToken(gloxtoken.Identifier, text)
	}
}

func (s *Scanner) addGenericToken(token gloxtoken.TokenType, literal any) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, gloxtoken.Token{TokenType: token, Lexeme: text, Literal: literal, Line: s.line})
}

func (s *Scanner) scanToken() {
	c := s.advance()
	if c == ' ' || c == '\r' || c == '\t' {
		return
	}
	switch c {
	case '(':
		s.addToken(gloxtoken.LeftParen)
	case ')':
		s.addToken(gloxtoken.RightParen)
	case '{':
		s.addToken(gloxtoken.LeftBrace)
	case '}':
		s.addToken(gloxtoken.RightBrace)
	case ',':
		s.addToken(gloxtoken.Comma)
	case '.':
		s.addToken(gloxtoken.Dot)
	case '-':
		s.addToken(gloxtoken.Minus)
	case '+':
		s.addToken(gloxtoken.Plus)
	case ';':
		s.addToken(gloxtoken.Semicolon)
	case '*':
		s.addToken(gloxtoken.Star)
	case '!':
		if s.match('=') {
			s.addToken(gloxtoken.BangEqual)
		} else {
			s.addToken(gloxtoken.Bang)
		}
	case '=':
		if s.match('=') {
			s.addToken(gloxtoken.EqualEqual)
		} else {
			s.addToken(gloxtoken.Equal)
		}
	case '<':
		if s.match('=') {
			s.addToken(gloxtoken.LessEqual)
		} else {
			s.addToken(gloxtoken.Less)
		}
	case '>':
		if s.match('=') {
			s.addToken(gloxtoken.GreaterEqual)
		} else {
			s.addToken(gloxtoken.Greater)
		}
	case '/':
		if s.match('/') {
			// A comment goes until the end of the line
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(gloxtoken.Slash)
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
			s.errorHandler.Report(s.line, "", "Unexpected character.")
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
