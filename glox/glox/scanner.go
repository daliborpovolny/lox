package main

import (
	"fmt"
	"strconv"
)

type Scanner struct {
	source string
	tokens []Token

	start   int
	current int
	line    int
}

func NewScanner(source string) *Scanner {
	s := Scanner{
		source: source,
	}

	s.tokens = make([]Token, 0, 10)
	return &s
}

func (s *Scanner) scanTokens() []Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, Token{
		tokenType: EOF,
		lexeme:    "",
		object:    nil,
		line:      s.line,
	})

	return s.tokens
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(LEFT_PAREN, nil)
	case ')':
		s.addToken(RIGHT_PAREN, nil)
	case '{':
		s.addToken(LEFT_BRACE, nil)
	case '}':
		s.addToken(RIGHT_BRACE, nil)
	case ',':
		s.addToken(COMMA, nil)
	case '.':
		s.addToken(DOT, nil)
	case '-':
		s.addToken(MINUS, nil)
	case '+':
		s.addToken(PLUS, nil)
	case ';':
		s.addToken(SEMICOLON, nil)
	case '*':
		s.addToken(STAR, nil)
	case '!':
		isNE := s.match('=')
		if isNE {
			s.addToken(BANG_EQUAL, nil)
		} else {
			s.addToken(BANG, nil)
		}
	case '=':
		isEQEQ := s.match('=')
		if isEQEQ {
			s.addToken(EQUAL_EQUAL, nil)
		} else {
			s.addToken(EQUAL, nil)
		}
	case '<':
		isLE := s.match('=')
		if isLE {
			s.addToken(LESS_EQUAL, nil)
		} else {
			s.addToken(LESS, nil)
		}
	case '>':
		isGE := s.match('=')
		if isGE {
			s.addToken(GREATER_EQUAL, nil)
		} else {
			s.addToken(GREATER, nil)
		}
	case '/':
		if s.match('/') {
			for s.peek() != 0 && !s.isAtEnd() {
				s.advance()
			}
		} else if s.match('*') {
			for {
				next := s.peek()
				nextNext := s.peekNext()

				if next == '*' && nextNext == '/' {
					s.current += 2
					break
				}

				if next == '\n' {
					s.line++
				}

				if s.isAtEnd() {
					lox.error(s.line, "Nonterminated multiline comment")
				}

				s.advance()
			}
		} else {
			s.addToken(SLASH, nil)
		}
	case ' ':
	case '\r':
	case '\t':
	case '\n':
		s.line++
	case '"':
		s.string()
	default:
		if isDigit(c) {
			s.number()
		} else if isAlpha(c) {
			s.identifier()
		} else {
			lox.error(s.line, "Unexpected character")
		}
	}
}

func (s *Scanner) addToken(tokenType TokenType, object Object) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, Token{
		tokenType: tokenType,
		lexeme:    text,
		object:    object,
		line:      s.line,
	})
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

// move the current counter, return next character
func (s *Scanner) advance() rune {
	s.current++
	return rune(s.source[s.current-1])

}

// peek at the next character, 0 if end reached
func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return 0
	}
	return rune(s.source[s.current])
}

// peek at the next next character, 0 if end reached
func (s *Scanner) peekNext() rune {
	if (s.current + 1) >= len(s.source) {
		return 0
	}
	return rune(s.source[s.current+1])
}

// if the next character is the expected one -> consume it and return true, otherwise return false
func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}

	if rune(s.source[s.current]) != expected {
		return false
	}

	s.current++
	return true
}

// process a string
func (s *Scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {

		if s.peek() == '\n' {
			s.line++
		}

		s.advance()
	}

	if s.isAtEnd() {
		lox.error(s.line, "Unterminated string")
		return
	}

	s.advance()

	value := s.source[s.start+1 : s.current-1]
	s.addToken(STRING, value)
}

// process a number
func (s *Scanner) number() {

	for isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.advance()

		for isDigit(s.peek()) {
			s.advance()
		}
	}

	numberAsString := s.source[s.start:s.current]
	n, err := strconv.ParseFloat(numberAsString, 64)
	if err != nil {
		fmt.Println()
		panic(err)
	}

	s.addToken(NUMBER, n)

}

// process an identifier
func (s *Scanner) identifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := s.source[s.start:s.current]
	tokenType := keyword(text)

	s.addToken(tokenType, nil)
}

func isDigit(c rune) bool {
	if c >= '0' && c <= '9' {
		return true
	}
	return false
}

func isAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func isAlphaNumeric(c rune) bool {
	return isAlpha(c) || isDigit(c)
}
