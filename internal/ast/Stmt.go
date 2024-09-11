package ast

type Stmt interface {
	Accept(visitor StmtVisitor[any]) (any, error)
}

type Expression struct {
	Expression_ Expr
}

func NewStmtExpression(Expression_ Expr) Stmt { return &Expression{Expression_: Expression_} }

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

type StmtVisitor[T any] interface {
	VisitStmtExpression(stmt *Expression) (T, error)
	VisitStmtPrint(stmt *Print) (T, error)
}
