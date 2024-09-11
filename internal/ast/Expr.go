package ast

import "github.com/codecrafters-io/interpreter-starter-go/internal/token"

type Expr interface {
	Accept(visitor ExprVisitor[any]) (any, error)
}

type Binary struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func NewExprBinary(Left Expr, Operator token.Token, Right Expr) Expr {
	return &Binary{Left: Left, Operator: Operator, Right: Right}
}

func (b *Binary) Accept(visitor ExprVisitor[any]) (any, error) {
	return visitor.VisitExprBinary(b)
}

type Grouping struct {
	Expression Expr
}

func NewExprGrouping(Expression Expr) Expr { return &Grouping{Expression: Expression} }

func (g *Grouping) Accept(visitor ExprVisitor[any]) (any, error) {
	return visitor.VisitExprGrouping(g)
}

type Literal struct {
	Value any
}

func NewExprLiteral(Value any) Expr { return &Literal{Value: Value} }

func (l *Literal) Accept(visitor ExprVisitor[any]) (any, error) {
	return visitor.VisitExprLiteral(l)
}

type Unary struct {
	Operator token.Token
	Right    Expr
}

func NewExprUnary(Operator token.Token, Right Expr) Expr {
	return &Unary{Operator: Operator, Right: Right}
}

func (u *Unary) Accept(visitor ExprVisitor[any]) (any, error) {
	return visitor.VisitExprUnary(u)
}

type Ternary struct {
	Test     Expr
	Question token.Token
	Left     Expr
	Colon    token.Token
	Right    Expr
}

func NewExprTernary(Test Expr, Question token.Token, Left Expr, Colon token.Token, Right Expr) Expr {
	return &Ternary{Test: Test, Question: Question, Left: Left, Colon: Colon, Right: Right}
}

func (t *Ternary) Accept(visitor ExprVisitor[any]) (any, error) {
	return visitor.VisitExprTernary(t)
}

type ExprVisitor[T any] interface {
	VisitExprBinary(expr *Binary) (T, error)
	VisitExprGrouping(expr *Grouping) (T, error)
	VisitExprLiteral(expr *Literal) (T, error)
	VisitExprUnary(expr *Unary) (T, error)
	VisitExprTernary(expr *Ternary) (T, error)
}
