package token

import (
	"fmt"
)

type Type int

func IsKeyword(identifier string) (bool, Type) {
	if v, ok := Keywords[identifier]; ok {
		return true, v
	}
	return false, 0
}

const (
	// Single-character tokens.

	//
	LEFT_PAREN Type = iota + 1
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	SLASH
	STAR

	// One or two character tokens.
	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL

	// Literals.
	IDENTIFIER
	STRING
	NUMBER

	// Keywords.

	AND
	CLASS
	ELSE
	FALSE
	FUN
	FOR
	IF
	NIL
	OR
	PRINT
	RETURN
	SUPER
	THIS
	TRUE
	VAR
	WHILE

	EOF
)

var (
	Keywords = map[string]Type{
		"and":    AND,
		"class":  CLASS,
		"else":   ELSE,
		"false":  FALSE,
		"for":    FOR,
		"fun":    FUN,
		"if":     IF,
		"nil":    NIL,
		"or":     OR,
		"print":  PRINT,
		"return": RETURN,
		"super":  SUPER,
		"this":   THIS,
		"true":   TRUE,
		"var":    VAR,
		"while":  WHILE,
	}

	Keywords2Str = map[Type]string{
		AND:    "and",
		CLASS:  "class",
		ELSE:   "else",
		FALSE:  "false",
		FOR:    "for",
		FUN:    "fun",
		IF:     "if",
		NIL:    "nil",
		OR:     "or",
		PRINT:  "print",
		RETURN: "return",
		SUPER:  "super",
		THIS:   "this",
		TRUE:   "true",
		VAR:    "var",
		WHILE:  "while",
	}
)

func IsKeywordType(t Type) bool {
	return t >= AND && t <= WHILE
}

func (t Type) String() string {
	switch t {
	case LEFT_PAREN:
		return "LEFT_PAREN"
	case RIGHT_PAREN:
		return "RIGHT_PAREN"
	case LEFT_BRACE:
		return "LEFT_BRACE"
	case RIGHT_BRACE:
		return "RIGHT_BRACE"
	case COMMA:
		return "COMMA"
	case DOT:
		return "DOT"
	case MINUS:
		return "MINUS"
	case PLUS:
		return "PLUS"
	case SEMICOLON:
		return "SEMICOLON"
	case SLASH:
		return "SLASH"
	case STAR:
		return "STAR"
	case BANG:
		return "BANG"
	case BANG_EQUAL:
		return "BANG_EQUAL"
	case EQUAL:
		return "EQUAL"
	case EQUAL_EQUAL:
		return "EQUAL_EQUAL"
	case GREATER:
		return "GREATER"
	case GREATER_EQUAL:
		return "GREATER_EQUAL"
	case LESS:
		return "LESS"
	case LESS_EQUAL:
		return "LESS_EQUAL"
	case IDENTIFIER:
		return "IDENTIFIER"
	case STRING:
		return "STRING"
	case NUMBER:
		return "NUMBER"
	case AND:
		return "AND"
	case CLASS:
		return "CLASS"
	case ELSE:
		return "ELSE"
	case FALSE:
		return "FALSE"
	case FUN:
		return "FUN"
	case FOR:
		return "FOR"
	case IF:
		return "IF"
	case NIL:
		return "NIL"
	case OR:
		return "OR"
	case PRINT:
		return "PRINT"
	case RETURN:
		return "RETURN"
	case SUPER:
		return "SUPER"
	case THIS:
		return "THIS"
	case TRUE:
		return "TRUE"
	case VAR:
		return "VAR"
	case WHILE:
		return "WHILE"
	case EOF:
		return "EOF"
	}
	panic("unknown token")
}

// how to represent token in string
func (t Type) Repr(obj any) string {
	if IsKeywordType(t) {
		return Keywords2Str[t]
	}
	switch t {
	case LEFT_PAREN:
		return "("
	case RIGHT_PAREN:
		return ")"
	case LEFT_BRACE:
		return "{"
	case RIGHT_BRACE:
		return "}"
	case COMMA:
		return ","
	case DOT:
		return "."
	case MINUS:
		return "-"
	case PLUS:
		return "+"
	case SEMICOLON:
		return ";"
	case SLASH:
		return "/"
	case STAR:
		return "*"
	case BANG:
		return "!"
	case BANG_EQUAL:
		return "!="
	case EQUAL:
		return "="
	case EQUAL_EQUAL:
		return "=="
	case GREATER:
		return ">"
	case GREATER_EQUAL:
		return ">="
	case LESS:
		return "<"
	case LESS_EQUAL:
		return "<="
	case IDENTIFIER:
		return fmt.Sprintf("%s", obj)
	case STRING:
		return fmt.Sprintf("\"%s\"", obj)
	case NUMBER:
		return fmt.Sprintf("%s", obj)
	case AND:
		return "AND"
	case CLASS:
		return "CLASS"
	case ELSE:
		return "ELSE"
	case FALSE:
		return "false"
	case FUN:
		return "FUN"
	case FOR:
		return "FOR"
	case IF:
		return "IF"
	case NIL:
		return "nil"
	//case OR:
	//	return "OR"
	//case PRINT:
	//	return "PRINT"
	//case RETURN:
	//	return "RETURN"
	//case SUPER:
	//	return "SUPER"
	//case THIS:
	//	return "THIS"
	case TRUE:
		return "true"
	case VAR:
		return "var"
	//case WHILE:
	//	return "WHILE"
	case EOF:
		return ""
	}
	panic("unknown token")
}
