package main

import (
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/internal/ast"
	"github.com/codecrafters-io/interpreter-starter-go/internal/token"
	"github.com/codecrafters-io/interpreter-starter-go/internal/visitor"
)

func main() {
	v := &visitor.AstPrinter{}
	expr := &ast.Binary{
		Left: &ast.Literal{
			Value: 2,
		},
		Operator: token.NewToken(token.PLUS, "+", nil, 1),
		Right: &ast.Literal{
			Value: 1,
		},
	}
	fmt.Println(v.Print(expr))
}
