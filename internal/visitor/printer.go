package visitor

import (
	"fmt"
	"strings"

	ast "github.com/codecrafters-io/interpreter-starter-go/cmd/generator/out"
)

type AstPrinter struct {
}

func (a *AstPrinter) VisitExprBinary(expr *ast.Binary) any {
	return a.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (a *AstPrinter) VisitExprGrouping(expr *ast.Grouping) any {
	return a.parenthesize("group", expr.Expression)
}

func (a *AstPrinter) VisitExprLiteral(expr *ast.Literal) any {
	if expr.Value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", expr.Value)
}

func (a *AstPrinter) VisitExprUnary(expr *ast.Unary) any {
	return a.parenthesize(expr.Operator.Lexeme, expr.Right)
}

var (
	_ ast.Visitor[any] = &AstPrinter{}
)

func (a *AstPrinter) Print(expr ast.Expr) string {
	return expr.Accept(a).(string)
}

func (a *AstPrinter) parenthesize(name string, exprs ...ast.Expr) string {
	sb := &strings.Builder{}
	sb.WriteRune('(')
	sb.WriteString(name)
	for _, expr := range exprs {
		sb.WriteRune(' ')
		sb.WriteString(a.Print(expr))
	}
	sb.WriteRune(')')
	return sb.String()
}
