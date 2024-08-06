package main

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/interpreter-starter-go/internal/loxscanner"
	"github.com/codecrafters-io/interpreter-starter-go/internal/parser"
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
		handleTokenize()
	}
	if command == "parse" {
		handleParse()
		os.Exit(exitCodeSuccess)
	}

	fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
	os.Exit(1)
}

func handleParse() {
	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	sc := loxscanner.NewScanner(string(fileContents))
	tokens := sc.ScanAll()
	if sc.Errors() != nil {
		for _, err := range sc.Errors() {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		}
		os.Exit(exitCodeScanError)
	}
	p := parser.NewParser(tokens)
	expr := p.Parse()
	if p.Errors() != nil {
		for _, err := range p.Errors() {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		}
		os.Exit(exitCodeParseError)
	}
	v := &visitor.AstPrinter{}
	fmt.Println(v.Print(expr))
}

const (
	exitCodeSuccess    = 0
	exitCodeScanError  = 65
	exitCodeParseError = 65
)

func handleTokenize() {
	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	sc := loxscanner.NewScanner(string(fileContents))
	tokens := sc.ScanAll()
	exitCode := exitCodeSuccess
	if sc.Errors() != nil {
		exitCode = exitCodeScanError
		for _, err := range sc.Errors() {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		}
	}
	for _, token := range tokens {
		fmt.Println(token.String())
	}
	os.Exit(exitCode)
}
