package loxscanner

import (
	"bytes"
	"fmt"
	"strings"
	"text/scanner"

	"github.com/codecrafters-io/interpreter-starter-go/internal/token"
)

type Scanner struct {
	line          int
	content       []rune
	contentOffset int
	tokens        []*token.Token
	errors        []error
	startOffset   int
}

func (s *Scanner) Next() rune {
	idx := s.contentOffset
	if idx >= len(s.content) {
		return scanner.EOF
	}
	s.contentOffset++
	return s.content[idx]
}
func (s *Scanner) Peek() rune {
	idx := s.contentOffset
	if idx >= len(s.content) {
		return scanner.EOF
	}
	return s.content[idx]
}
func (s *Scanner) PeekNext() rune {
	idx := s.contentOffset
	if idx+1 >= len(s.content) {
		return scanner.EOF
	}
	return s.content[idx+1]
}

func NewScanner(src string) *Scanner {
	return &Scanner{
		line:    1,
		content: bytes.Runes([]byte(src)),
		tokens:  []*token.Token{},
	}
}

func (s *Scanner) scanToken() {
	s.startOffset = s.contentOffset
	next := s.Next()
	switch next {
	case '(':
		s.addToken(token.LEFT_PAREN)
	case ')':
		s.addToken(token.RIGHT_PAREN)
	case '{':
		s.addToken(token.LEFT_BRACE)
	case '}':
		s.addToken(token.RIGHT_BRACE)
	case ',':
		s.addToken(token.COMMA)
	case '.':
		s.addToken(token.DOT)
	case '-':
		s.addToken(token.MINUS)
	case '+':
		s.addToken(token.PLUS)
	case ';':
		s.addToken(token.SEMICOLON)
	case '*':
		s.addToken(token.STAR)
	case '!':
		if s.match('=') {
			s.addToken(token.BANG_EQUAL)
		} else {
			s.addToken(token.BANG)
		}
	case '=':
		if s.match('=') {
			s.addToken(token.EQUAL_EQUAL)
		} else {
			s.addToken(token.EQUAL)
		}
	case '<':
		if s.match('=') {
			s.addToken(token.LESS_EQUAL)
		} else {
			s.addToken(token.LESS)
		}
	case '>':
		if s.match('=') {
			s.addToken(token.GREATER_EQUAL)
		} else {
			s.addToken(token.GREATER)
		}
	case '/':
		if s.match('/') {
			for peek := s.Peek(); peek != '\n' &&
				peek != scanner.EOF; peek = s.Peek() {
				s.Next()
			}
		} else {
			s.addToken(token.SLASH)
		}
	case ' ':
	case '\r':
	case '\t':
	case '\n':
		s.line++
	case '"':
		s.scanString()
	case '?':
		s.addToken(token.QUESTION_MARK)
	case ':':
		s.addToken(token.COLON)
	default:
		if isDigit(next) {
			s.scanNumber(next)
			break
		} else if isAlpha(next) {
			s.scanIdentifier(next)
			break
		}
		s.errors = append(s.errors, fmt.Errorf("[line %d] Error: Unexpected character: %s",
			s.getLine(), string(next)))
	}

}

func isDigit(c rune) bool {
	return '0' <= c && c <= '9'
}
func isAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		c == '_'
}
func isAlphaNumeric(c rune) bool {
	return isAlpha(c) || isDigit(c)
}
func (s *Scanner) scanIdentifier(curr rune) {
	sb := &strings.Builder{}
	sb.WriteRune(curr)
	for isAlphaNumeric(s.Peek()) {
		sb.WriteRune(s.Next())
	}
	s.addIdentifier(sb.String())
}

func (s *Scanner) ScanAll() []*token.Token {
	for s.Peek() != scanner.EOF {
		s.scanToken()
	}
	s.addToken(token.EOF)
	return s.tokens
}
func (s *Scanner) scanString() {
	sb := &strings.Builder{}
	for s.Peek() != '"' && s.Peek() != scanner.EOF {
		sb.WriteRune(s.Next())
	}
	if s.Peek() == scanner.EOF {
		s.errors = append(s.errors, fmt.Errorf("[line %d] Error: Unterminated string.", s.getLine()))
		return
	}
	s.Next()
	s2 := sb.String()
	s.addTokenLexeme(token.STRING, s2)

}

func (s *Scanner) addToken(t token.Type) {
	newToken := token.NewToken(t, t.Repr(nil), nil, s.getLine())
	s.tokens = append(s.tokens, &newToken)
}
func (s *Scanner) addTokenLexeme(t token.Type, obj interface{}) {
	newToken := token.NewToken(t, t.Repr(obj), obj, s.getLine())
	s.tokens = append(s.tokens, &newToken)
}
func (s *Scanner) addNumberToken(numStr string) {
	newToken := token.NewNumberToken(numStr, s.getLine())
	s.tokens = append(s.tokens, &newToken)
}
func (s *Scanner) addIdentifier(i string) {
	if is, type_ := token.IsKeyword(i); is {
		s.addToken(type_)
	} else {
		newToken := token.NewToken(token.IDENTIFIER, token.IDENTIFIER.Repr(i), nil, s.getLine())
		s.tokens = append(s.tokens, &newToken)
	}
}

func (s *Scanner) getLine() int {
	pos := s.line
	return pos
}

func (s *Scanner) match(expected rune) bool {
	if s.Peek() == expected {
		s.Next()
		return true
	}
	return false
}

func (s *Scanner) Errors() []error {
	return s.errors
}

func (s *Scanner) scanNumber(firstDigit rune) {
	sb := &strings.Builder{}
	sb.WriteRune(firstDigit)
	var next rune
	for isDigit(s.Peek()) {
		sb.WriteRune(s.Next())
	}
	if s.Peek() == '.' && isDigit(s.PeekNext()) {
		next = s.Next()
		sb.WriteRune(next)
		for isDigit(s.Peek()) {
			sb.WriteRune(s.Next())
		}
	}
	s.addNumberToken(sb.String())
}
