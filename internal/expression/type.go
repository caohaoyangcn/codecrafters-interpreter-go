package expression

import (
	"github.com/codecrafters-io/interpreter-starter-go/internal/token"
)

type Expression interface{}

type Binary struct {
	Left     Expression
	Operator token.Token
	Right    Expression
}

type Visitor[T any] interface {
	Visit(expr Expression) T
}
