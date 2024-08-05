package loxscanner

import (
	"bytes"
	"fmt"
	"strconv"
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
		break
	case ')':
		s.addToken(token.RIGHT_PAREN)
		break
	case '{':
		s.addToken(token.LEFT_BRACE)
		break
	case '}':
		s.addToken(token.RIGHT_BRACE)
		break
	case ',':
		s.addToken(token.COMMA)
		break
	case '.':
		s.addToken(token.DOT)
		break
	case '-':
		s.addToken(token.MINUS)
		break
	case '+':
		s.addToken(token.PLUS)
		break
	case ';':
		s.addToken(token.SEMICOLON)
		break
	case '*':
		s.addToken(token.STAR)
		break
	case '!':
		if s.match('=') {
			s.addToken(token.BANG_EQUAL)
		} else {
			s.addToken(token.BANG)
		}
		break
	case '=':
		if s.match('=') {
			s.addToken(token.EQUAL_EQUAL)
		} else {
			s.addToken(token.EQUAL)
		}
		break
	case '<':
		if s.match('=') {
			s.addToken(token.LESS_EQUAL)
		} else {
			s.addToken(token.LESS)
		}
		break
	case '>':
		if s.match('=') {
			s.addToken(token.GREATER_EQUAL)
		} else {
			s.addToken(token.GREATER)
		}
		break
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
	default:
		if isDigit(next) {
			s.scanNumber(next)
		}
		s.errors = append(s.errors, fmt.Errorf("[line %d] Error: Unexpected character: %s",
			s.getLine(), string(next)))
		break
	}

}

func isDigit(c rune) bool {
	return '0' <= c && c <= '9'
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
	if s.Peek() == '.' && isDigit(s.Peek()) {
		next = s.Next()
		sb.WriteRune(next)
		for isDigit(s.Peek()) {
			sb.WriteRune(s.Next())
		}
	}
	num, _ := strconv.ParseFloat(sb.String(), 64)
	s.addTokenLexeme(token.NUMBER, num)
}
