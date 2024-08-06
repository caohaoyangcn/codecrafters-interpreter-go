package token

import (
	"fmt"
	"strconv"
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
	num, _ := strconv.ParseFloat(numStr, 64)

	return NewToken(NUMBER, origToken, num, line)

}

func (t Token) String() string {
	object := t.Object
	objStr := ""
	if object == nil {
		objStr = "null"
		return fmt.Sprintf("%s %v %v", t.Type, t.Lexeme, objStr)
	}
	switch object.(type) {
	case float64:
		val := strconv.FormatFloat(object.(float64), 'f', -1, 64)
		if strings.Contains(val, ".") {
			objStr = val
		} else {
			objStr = val + ".0"
		}
	default:
		objStr = fmt.Sprintf("%v", object)
	}
	return fmt.Sprintf("%s %v %v", t.Type, t.Lexeme, objStr)
}
