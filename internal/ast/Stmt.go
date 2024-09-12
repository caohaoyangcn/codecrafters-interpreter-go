package ast

import "github.com/codecrafters-io/interpreter-starter-go/internal/token"

type Stmt interface {
	Accept(visitor StmtVisitor[any]) (any, error)
}

type Block struct {
	Statements []Stmt
}

func NewStmtBlock(Statements []Stmt) Stmt { return &Block{Statements: Statements} }

func (b *Block) Accept(visitor StmtVisitor[any]) (any, error) {
	return visitor.VisitStmtBlock(b)
}

type Expression struct {
	Expression_  Expr
	HasSemicolon bool
}

func NewStmtExpression(Expression_ Expr, HasSemicolon bool) Stmt {
	return &Expression{Expression_: Expression_, HasSemicolon: HasSemicolon}
}

func (e *Expression) Accept(visitor StmtVisitor[any]) (any, error) {
	return visitor.VisitStmtExpression(e)
}

type Print struct {
	Expression_ Expr
}

func NewStmtPrint(Expression_ Expr) Stmt { return &Print{Expression_: Expression_} }

func (p *Print) Accept(visitor StmtVisitor[any]) (any, error) {
	return visitor.VisitStmtPrint(p)
}

type Var struct {
	Name        token.Token
	Initializer Expr
}

func NewStmtVar(Name token.Token, Initializer Expr) Stmt {
	return &Var{Name: Name, Initializer: Initializer}
}

func (v *Var) Accept(visitor StmtVisitor[any]) (any, error) {
	return visitor.VisitStmtVar(v)
}

type StmtVisitor[T any] interface {
	VisitStmtBlock(stmt *Block) (T, error)
	VisitStmtExpression(stmt *Expression) (T, error)
	VisitStmtPrint(stmt *Print) (T, error)
	VisitStmtVar(stmt *Var) (T, error)
}
