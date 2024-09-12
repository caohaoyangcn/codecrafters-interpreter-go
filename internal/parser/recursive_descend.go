package parser

import (
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/internal/ast"
	"github.com/codecrafters-io/interpreter-starter-go/internal/token"
)

type ParseMode int

const (
	// Default mode
	Default ParseMode = iota
	REPL
)

type Parser struct {
	tokens []*token.Token
	curr   int
	errors []error
	//mode   ParseMode
}

func NewParser(tokens []*token.Token) *Parser {
	return &Parser{tokens: tokens}
}

//	func (p *Parser) SetMode(mode ParseMode) {
//		p.mode = mode
//	}
func (p *Parser) Errors() []error {
	return p.errors
}

func (p *Parser) Parse() []ast.Stmt {
	var stmts []ast.Stmt
	for !p.atEnd() {
		stmt, err := p.Declaration()
		if err != nil {
			p.errors = append(p.errors, err)
			return nil
		}
		stmts = append(stmts, stmt)
	}
	return stmts
}

// Statement implements the statement rule
//
// statement      → exprStmt
//
//		| printStmt
//	    | block ;
func (p *Parser) Statement() (ast.Stmt, error) {
	if p.match(token.PRINT) {
		return p.PrintStatement()
	}
	if p.match(token.LEFT_BRACE) {
		return p.block()
	}
	expr, err := p.Expression()
	if err != nil {
		return nil, fmt.Errorf("statement: %w", err)
	}
	//p.consume(token.SEMICOLON, "Expect ';' after expression.")
	if _, err := p.consume(token.SEMICOLON, "Expect ';' after expression."); err != nil {
		return ast.NewStmtExpression(expr, false), nil
	}
	return ast.NewStmtExpression(expr, true), nil
}

func (p *Parser) block() (ast.Stmt, error) {
	var enclosingStatements []ast.Stmt
	for !p.check(token.RIGHT_BRACE) && !p.atEnd() {
		d, err := p.Declaration()
		if err != nil {
			p.errors = append(p.errors, err)
			continue
		}
		enclosingStatements = append(enclosingStatements, d)
	}
	_, err := p.consume(token.RIGHT_BRACE, "Expect '}' after block.")
	if err != nil {
		return nil, err
	}
	return ast.NewStmtBlock(enclosingStatements), nil
}

func (p *Parser) varDeclaration() (ast.Stmt, error) {
	tok, err := p.consume(token.IDENTIFIER, "Expect variable name.")
	if err != nil {
		return nil, err
	}
	var expr ast.Expr
	if p.match(token.EQUAL) {
		expr, err = p.Expression()
		if err != nil {
			return nil, fmt.Errorf("varDeclaration: %w", err)
		}
	}
	if _, err := p.consume(token.SEMICOLON, "Expect ';' at the end of var declaration."); err != nil {
		return nil, err
	}
	return ast.NewStmtVar(*tok, expr), nil
}

// Declaration implements the declaration rule
// declaration    → varDecl
//
//	| statement ;
func (p *Parser) Declaration() (ast.Stmt, error) {
	if p.match(token.VAR) {
		if stmt, err := p.varDeclaration(); err != nil {
			p.synchronize()
			p.errors = append(p.errors, err)
			return nil, nil
		} else {
			return stmt, nil
		}
	}
	return p.Statement()
}

// PrintStatement implements the print statement rule
//
//	printStmt      → "print" expression ";" ;
func (p *Parser) PrintStatement() (ast.Stmt, error) {
	value, err := p.Expression()
	if err != nil {
		return nil, fmt.Errorf("PrintStatement: %w", err)
	}
	if _, err := p.consume(token.SEMICOLON, "Expect ';' after expression."); err != nil {
		return nil, err
	}
	return ast.NewStmtPrint(value), nil
}

func (p *Parser) Expression() (ast.Expr, error) {
	return p.Assignment()
}

func (p *Parser) Assignment() (ast.Expr, error) {
	expr, err := p.Comma()
	if err != nil {
		return nil, fmt.Errorf("assignment: %w", err)
	}
	if p.match(token.EQUAL) {
		tok := p.previous()
		right, err := p.Assignment()
		if err != nil {
			return nil, fmt.Errorf("assignment: %w", err)
		}
		if val, ok := expr.(*ast.Variable); ok {
			tok := val.Name
			return ast.NewExprAssign(tok, right), nil
		}
		// TODO throw error or just record it?
		p.errors = append(p.errors, errorFunc(*tok, "Invalid assignment target."))
	}
	return expr, nil
}

// Comma implements the comma operator in C
// ref: https://en.wikipedia.org/wiki/Comma_operator
func (p *Parser) Comma() (ast.Expr, error) {
	if err := p.checkBinaryOperatorHasLeftOperand(token.COMMA); err != nil {
		return nil, err
	}
	expr, err := p.Ternary()
	if err != nil {
		return nil, fmt.Errorf("comma: %w", err)
	}
	for p.match(token.COMMA) {
		operator := p.previous()
		rightExpr, err := p.Ternary() // TODO
		if err != nil {
			return nil, fmt.Errorf("comma: %w", err)
		}
		expr = ast.NewExprBinary(expr, *operator, rightExpr)
	}
	return expr, nil
}
func (p *Parser) Ternary() (ast.Expr, error) {
	expr, err := p.Equality()
	if err != nil {
		return nil, fmt.Errorf("ternary: %w", err)
	}
	for p.match(token.QUESTION_MARK) {
		q := p.previous()
		leftExpr, err := p.Equality()
		if err != nil {
			return nil, fmt.Errorf("ternary: %w", err)
		}
		if _, err := p.consume(token.COLON, "Expect ':' after '?'."); err != nil {
			return nil, err
		}
		c := p.previous()
		rightExpr, err := p.Equality()
		if err != nil {
			return nil, fmt.Errorf("ternary: %w", err)
		}
		expr = ast.NewExprTernary(expr, *q, leftExpr, *c, rightExpr)
	}

	return expr, nil
}
func (p *Parser) Equality() (ast.Expr, error) {
	if err := p.checkBinaryOperatorHasLeftOperand(token.EQUAL_EQUAL, token.BANG_EQUAL); err != nil {
		return nil, err
	}
	expr, err := p.Comparison()
	if err != nil {
		return nil, fmt.Errorf("equality: %w", err)
	}
	for p.match(token.EQUAL_EQUAL, token.BANG_EQUAL) {
		operator := p.previous()
		rightExpr, err := p.Comparison()

		if err != nil {
			return nil, fmt.Errorf("equality: %w", err)
		}
		expr = ast.NewExprBinary(expr, *operator, rightExpr)
	}
	return expr, nil
}
func (p *Parser) Comparison() (ast.Expr, error) {
	if err := p.checkBinaryOperatorHasLeftOperand(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL); err != nil {
		return nil, err
	}
	expr, err := p.Term()
	if err != nil {
		return nil, fmt.Errorf("comparison: %w", err)
	}
	for p.match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
		operator := p.previous()
		rightExpr, err := p.Term()
		if err != nil {
			return nil, fmt.Errorf("comparison: %w", err)
		}
		expr = ast.NewExprBinary(expr, *operator, rightExpr)
	}
	return expr, nil
}
func (p *Parser) Term() (ast.Expr, error) {
	if err := p.checkBinaryOperatorHasLeftOperand(token.PLUS); err != nil {
		tok := p.peek()
		return nil, errorFunc(tok, fmt.Sprintf("%s: left operand required", tok.Lexeme))
	}
	expr, err := p.Factor()
	if err != nil {
		return nil, fmt.Errorf("term: %w", err)
	}
	for p.match(token.MINUS, token.PLUS) {
		operator := p.previous()
		rightExpr, err := p.Factor()
		if err != nil {
			return nil, fmt.Errorf("term: %w", err)
		}

		expr = ast.NewExprBinary(expr, *operator, rightExpr)
	}
	return expr, nil
}
func (p *Parser) Factor() (ast.Expr, error) {
	if err := p.checkBinaryOperatorHasLeftOperand(token.SLASH, token.STAR); err != nil {
		return nil, err
	}
	expr, err := p.Unary()
	if err != nil {
		return nil, fmt.Errorf("factor: %w", err)
	}
	for p.match(token.SLASH, token.STAR) {
		operator := p.previous()
		rightExpr, err := p.Unary()
		if err != nil {
			return nil, fmt.Errorf("factor: %w", err)
		}
		expr = ast.NewExprBinary(expr, *operator, rightExpr)
	}
	return expr, nil
}
func (p *Parser) Unary() (ast.Expr, error) {
	if p.match(token.BANG, token.MINUS) {
		operator := p.previous()
		rightExpr, err := p.Unary()
		if err != nil {
			return nil, fmt.Errorf("unary: %w", err)
		}
		return ast.NewExprUnary(*operator, rightExpr), nil
	}
	return p.Primary()
}
func (p *Parser) Primary() (ast.Expr, error) {
	if p.match(token.FALSE) {
		return ast.NewExprLiteral(false), nil
	}
	if p.match(token.TRUE) {
		return ast.NewExprLiteral(true), nil
	}
	if p.match(token.NIL, token.NUMBER, token.STRING) {
		return ast.NewExprLiteral(p.previous().Object), nil
	}
	if p.match(token.IDENTIFIER) {
		return ast.NewExprVariable(*p.previous()), nil
	}
	if p.match(token.LEFT_PAREN) {
		expr, err := p.Expression()
		if err != nil {
			return nil, fmt.Errorf("primary: %w", err)
		}
		if _, err := p.consume(token.RIGHT_PAREN, "Expect ')' after expression."); err != nil {
			return nil, fmt.Errorf("primary: %w", err)
		}
		return ast.NewExprGrouping(expr), nil
	}
	return nil, errorFunc(p.peek(), "primary: expect expression")
}

func (p *Parser) peekMatch(types ...token.Type) bool {
	for _, t := range types {
		if p.check(t) {
			return true
		}
	}
	return false
}
func (p *Parser) match(types ...token.Type) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}
func (p *Parser) peek() token.Token {
	return *p.tokens[p.curr]
}
func (p *Parser) advance() *token.Token {
	if !p.atEnd() {
		p.curr++
	}
	return p.previous()
}
func (p *Parser) atEnd() bool {
	return p.peek().Type == token.EOF
}
func (p *Parser) previous() *token.Token {
	return p.tokens[p.curr-1]
}
func (p *Parser) check(t token.Type) bool {
	if p.atEnd() {
		return false
	}
	tok := p.peek()
	return tok.Type == t
}
func (p *Parser) consume(t token.Type, msg string) (*token.Token, error) {
	if p.check(t) {
		return p.advance(), nil
	}
	return nil, errorFunc(p.peek(), msg)
}

func errorFunc(tok token.Token, msg string) error {
	if tok.Type == token.EOF {
		return fmt.Errorf("%d at end: %v", tok.Line, msg)
	}
	return fmt.Errorf("%d at '%v': %s", tok.Line, tok.Lexeme, msg)
}

func (p *Parser) synchronize() {
	p.advance()
	for !p.atEnd() {
		if p.previous().Type == token.SEMICOLON {
			return
		}
		switch p.peek().Type {
		case token.CLASS:
		case token.FUN:
		case token.VAR:
		case token.FOR:
		case token.IF:
		case token.WHILE:
		case token.PRINT:
		case token.RETURN:
			return
		}

		// discard tokens until we find a statement
		p.advance()
	}
}

func (p *Parser) checkBinaryOperatorHasLeftOperand(op ...token.Type) error {
	if p.peekMatch(op...) {
		tok := p.peek()
		return errorFunc(tok, fmt.Sprintf("%s: left operand required", tok.Lexeme))
	}
	return nil
}
