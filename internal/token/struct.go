package token

import (
	"fmt"
	"strings"
)

type Token struct {
	Type   Type
	Lexeme string
	Object interface{}
	Line   int
}

func NewToken(t Type, lexeme string, object interface{}, line int) Token {
	return Token{
		Type:   t,
		Lexeme: lexeme,
		Object: object,
		Line:   line,
	}
}
func NewNumberToken(numStr string, line int) Token {
	origToken := numStr
	if strings.IndexByte(numStr, '.') < 0 {
		numStr += ".0"
	}
	return NewToken(NUMBER, origToken, numStr, line)

}

func (t Token) String() string {
	object := t.Object
	if object == nil {
		object = "null"
	}
	return fmt.Sprintf("%s %v %v", t.Type, t.Lexeme, object)
}
