package visitor

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/codecrafters-io/interpreter-starter-go/internal/ast"
)

type AstPrinter struct {
}

func (a *AstPrinter) VisitStmtBlock(stmt *ast.Block) (any, error) {
	//TODO implement me
	panic("implement me")
}

func (a *AstPrinter) VisitStmtVar(stmt *ast.Var) (any, error) {
	//TODO implement me
	panic("implement me")
}

func (a *AstPrinter) VisitExprVariable(expr *ast.Variable) (any, error) {
	//TODO implement me
	panic("implement me")
}

func (a *AstPrinter) VisitExprAssign(expr *ast.Assign) (any, error) {
	//TODO implement me
	panic("implement me")
}

func (a *AstPrinter) VisitStmtExpression(stmt *ast.Expression) (any, error) {
	if stmt.HasSemicolon {
		return fmt.Sprintf("%s;", a.PrintExpr(stmt.Expression_)), nil
	}
	return fmt.Sprintf("%s", a.PrintExpr(stmt.Expression_)), nil
}

func (a *AstPrinter) VisitStmtPrint(stmt *ast.Print) (any, error) {
	return fmt.Sprintf("print %s", a.PrintExpr(stmt.Expression_)), nil
}

func (a *AstPrinter) VisitExprTernary(expr *ast.Ternary) (any, error) {
	return a.parenthesize(expr.Question.Lexeme+" "+expr.Colon.Lexeme, expr.Test, expr.Left, expr.Right), nil
}

func (a *AstPrinter) VisitExprBinary(expr *ast.Binary) (any, error) {
	return a.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right), nil
}

func (a *AstPrinter) VisitExprGrouping(expr *ast.Grouping) (any, error) {
	return a.parenthesize("group", expr.Expression), nil
}

func (a *AstPrinter) VisitExprLiteral(expr *ast.Literal) (any, error) {
	if expr.Value == nil {
		return "nil", nil
	}
	return ParserPrinter(expr.Value), nil
}

func (a *AstPrinter) VisitExprUnary(expr *ast.Unary) (any, error) {
	return a.parenthesize(expr.Operator.Lexeme, expr.Right), nil
}

var (
	_ ast.ExprVisitor[any] = &AstPrinter{}
	_ ast.StmtVisitor[any] = &AstPrinter{}
)

func (a *AstPrinter) PrintExpr(expr ast.Expr) string {
	accept, _ := expr.Accept(a)
	return accept.(string)
}
func (a *AstPrinter) PrintStmt(expr ast.Stmt) string {
	accept, _ := expr.Accept(a)
	return accept.(string)
}

func (a *AstPrinter) parenthesize(name string, exprs ...ast.Expr) string {
	sb := &strings.Builder{}
	sb.WriteRune('(')
	sb.WriteString(name)
	for _, expr := range exprs {
		sb.WriteRune(' ')
		sb.WriteString(a.PrintExpr(expr))
	}
	sb.WriteRune(')')
	return sb.String()
}

func ParserPrinter(obj any) string {
	if obj == nil {
		return "nil"
	}
	switch obj := obj.(type) {
	case float64:
		val := strconv.FormatFloat(obj, 'f', -1, 64)
		if strings.Contains(val, ".") {
			return val
		}
		return val + ".0"
	default:
		return fmt.Sprintf("%v", obj)
	}
}
