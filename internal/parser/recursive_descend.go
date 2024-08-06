package parser

import (
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/cmd/generator/out"
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

func (p *Parser) Parse() ast.Expr {
	expr, err := p.Expression()
	if err != nil {
		// TODO report error
		p.errors = append(p.errors, err)
		return nil
	}
	return expr
}

func (p *Parser) Expression() (ast.Expr, error) {
	return p.Equality()
}
func (p *Parser) Equality() (ast.Expr, error) {
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
