package ast

import "github.com/codecrafters-io/interpreter-starter-go/internal/token"

type Expr interface {
	Accept(visitor Visitor[any]) any
}

type Binary struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func NewExprBinary(Left Expr, Operator token.Token, Right Expr) Expr {
	return &Binary{Left: Left, Operator: Operator, Right: Right}
}

func (b *Binary) Accept(visitor Visitor[any]) any {
	return visitor.VisitExprBinary(b)
}

type Grouping struct {
	Expression Expr
}

func NewExprGrouping(Expression Expr) Expr { return &Grouping{Expression: Expression} }

func (g *Grouping) Accept(visitor Visitor[any]) any {
	return visitor.VisitExprGrouping(g)
}

type Literal struct {
	Value any
}

func NewExprLiteral(Value any) Expr { return &Literal{Value: Value} }

func (l *Literal) Accept(visitor Visitor[any]) any {
	return visitor.VisitExprLiteral(l)
}

type Unary struct {
	Operator token.Token
	Right    Expr
}

func NewExprUnary(Operator token.Token, Right Expr) Expr {
	return &Unary{Operator: Operator, Right: Right}
}

func (u *Unary) Accept(visitor Visitor[any]) any {
	return visitor.VisitExprUnary(u)
}

type Visitor[T any] interface {
	VisitExprBinary(expr *Binary) T
	VisitExprGrouping(expr *Grouping) T
	VisitExprLiteral(expr *Literal) T
	VisitExprUnary(expr *Unary) T
}
