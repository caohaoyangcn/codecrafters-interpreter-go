package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/codecrafters-io/interpreter-starter-go/internal/loxscanner"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	//fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	if command != "tokenize" {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	// Uncomment this block to pass the first stage
	//
	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	sc := loxscanner.NewScanner(bytes.NewReader(fileContents))
	tokens := sc.Scan()
	const (
		exitCodeSuccess   = 0
		exitCodeScanError = 65
	)
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
