package loxscanner

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"text/scanner"

	"github.com/codecrafters-io/interpreter-starter-go/internal/token"
)

type Scanner struct {
	sc             *scanner.Scanner
	tokens         []*token.Token
	errors         []error
	start, current scanner.Position
}

func NewScanner(src io.Reader) *Scanner {
	sc := &scanner.Scanner{}
	sc.Init(src)
	sc.Mode ^= scanner.ScanComments ^ scanner.SkipComments
	sc.Whitespace ^= scanner.GoWhitespace
	return &Scanner{
		sc:     sc,
		tokens: []*token.Token{},
	}
}

func (s *Scanner) scanToken() {
	s.start = s.sc.Position
	next := s.sc.Scan()
	s.current = s.sc.Pos()
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
			for peek := s.sc.Peek(); peek != '\n' &&
				peek != scanner.EOF; peek = s.sc.Peek() {
				s.sc.Next()
			}
		} else {
			s.addToken(token.SLASH)
		}
	case ' ':
	case '\r':
	case '\t':
	case '\n':
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

func (s *Scanner) Scan() []*token.Token {
	for s.sc.Peek() != scanner.EOF {
		s.scanToken()
	}
	s.addToken(token.EOF)
	return s.tokens
}
func (s *Scanner) scanString() {
	sb := &strings.Builder{}
	for s.sc.Peek() != '"' {
		sb.WriteRune(s.sc.Next())
	}
	if s.sc.Peek() == scanner.EOF {
		s.errors = append(s.errors, fmt.Errorf("[line %d] Error: Unterminated string.", s.getLine()))
		return
	}
	s.sc.Next()
	s.tokens = append(s.tokens, &token.Token{
		Type:   token.STRING,
		Lexeme: sb.String(),
		Line:   s.getLine(),
		Object: nil,
	})

}

func (s *Scanner) addToken(t token.Type) {
	newToken := token.NewToken(t, t.Repr(), nil, s.getLine())
	s.tokens = append(s.tokens, &newToken)
}
func (s *Scanner) addTokenLexeme(t token.Type, obj interface{}) {
	newToken := token.NewToken(t, t.Repr(), obj, s.getLine())
	s.tokens = append(s.tokens, &newToken)
}

func (s *Scanner) getLine() int {
	pos := s.sc.Pos()
	return pos.Line
}

func (s *Scanner) match(expected rune) bool {
	if s.sc.Peek() == expected {
		s.sc.Next()
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
	for isDigit(s.sc.Peek()) {
		sb.WriteRune(s.sc.Next())
	}
	if s.sc.Peek() == '.' && isDigit(s.sc.Peek()) {
		next = s.sc.Next()
		sb.WriteRune(next)
		for isDigit(s.sc.Peek()) {
			sb.WriteRune(s.sc.Next())
		}
	}
	num, _ := strconv.ParseFloat(sb.String(), 64)
	s.addTokenLexeme(token.NUMBER, num)
}
