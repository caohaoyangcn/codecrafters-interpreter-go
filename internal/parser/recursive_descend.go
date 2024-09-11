package parser

import (
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/internal/ast"
	"github.com/codecrafters-io/interpreter-starter-go/internal/token"
)

type Parser struct {
	tokens []*token.Token
	curr   int
	errors []error
}

func NewParser(tokens []*token.Token) *Parser {
	return &Parser{tokens: tokens}
}
func (p *Parser) Errors() []error {
	return p.errors
}

func (p *Parser) Parse() []ast.Stmt {
	var stmts []ast.Stmt
	for !p.atEnd() {
		stmt, err := p.Statement()
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
//	| printStmt ;
func (p *Parser) Statement() (ast.Stmt, error) {
	if p.match(token.PRINT) {
		return p.PrintStatement()
	}
	expr, err := p.Expression()
	if err != nil {
		return nil, fmt.Errorf("statement: %w", err)
	}
	if _, err := p.consume(token.SEMICOLON, "Expect ';' after expression."); err != nil {
		return nil, err
	}
	return ast.NewStmtExpression(expr), nil
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
	return p.Comma()
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
		rightExpr, err := p.Equality()
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
