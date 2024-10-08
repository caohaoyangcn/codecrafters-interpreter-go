package main

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/interpreter-starter-go/internal/ast"
	"github.com/codecrafters-io/interpreter-starter-go/internal/loxscanner"
	"github.com/codecrafters-io/interpreter-starter-go/internal/parser"
	"github.com/codecrafters-io/interpreter-starter-go/internal/token"
	"github.com/codecrafters-io/interpreter-starter-go/internal/visitor"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	//fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	if command == "tokenize" {
		tokens, errs := handleTokenize()
		for _, t := range tokens {
			fmt.Println(t.String())
		}
		if errs != nil {
			for _, err := range errs {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			}
			os.Exit(exitCodeScanError)
		}
		os.Exit(exitCodeSuccess)
	}
	if command == "parse" {
		tokens, errs := handleTokenize()
		if errs != nil {
			for _, err := range errs {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			}
			os.Exit(exitCodeScanError)
		}
		if stmts, errs := handleParse(tokens); errs != nil {
			for _, err := range errs {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			}
			os.Exit(exitCodeParseError)
		} else {
			v := &visitor.AstPrinter{}
			for _, stmt := range stmts {
				fmt.Println(v.PrintStmt(stmt))
			}
		}
		os.Exit(exitCodeSuccess)
	}
	if command == "evaluate" || command == "run" {
		tokens, errs := handleTokenize()
		if errs != nil {
			for _, err := range errs {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			}
			os.Exit(exitCodeScanError)
		}
		if expr, errs := handleParse(tokens); errs != nil {
			for _, err := range errs {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			}
			os.Exit(exitCodeParseError)
		} else {
			handleInterpret(expr)
			os.Exit(exitCodeSuccess)
		}
	}

	fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
	os.Exit(1)
}

func handleParse(tokens []*token.Token) ([]ast.Stmt, []error) {
	p := parser.NewParser(tokens)
	stmts := p.Parse()
	if p.Errors() != nil {
		return nil, p.Errors()
	}
	return stmts, nil
}

const (
	exitCodeSuccess    = 0
	exitCodeScanError  = 65
	exitCodeParseError = 65
	interpreterError   = 70
)

func handleTokenize() ([]*token.Token, []error) {
	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	sc := loxscanner.NewScanner(string(fileContents))
	tokens := sc.ScanAll()
	return tokens, sc.Errors()
}

func handleInterpret(expr []ast.Stmt) {
	i := visitor.NewInterpreter()
	_, err := i.Interpret(expr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(interpreterError)
	}
}
