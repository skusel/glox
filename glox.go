package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/skusel/glox/ast"
	"github.com/skusel/glox/langerr"
	"github.com/skusel/glox/parser"
	"github.com/skusel/glox/scanner"
)

func main() {
	numArgs := len(os.Args[1:])
	if numArgs > 1 {
		fmt.Println("Usage: glox [script]")
		os.Exit(64)
	} else if numArgs == 1 {
		runFile(os.Args[1])
	} else {
		runPrompt()
	}
}

func runFile(path string) {
	source, readErr := os.ReadFile(path)
	if readErr != nil {
		fmt.Println(readErr)
		os.Exit(2)
	} else {
		errorHandler := langerr.Handler{HadError: false}
		run(string(source), &errorHandler)
		if errorHandler.HadError {
			os.Exit(65)
		}
	}
}

func runPrompt() {
	errorHandler := langerr.Handler{HadError: false}
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
		} else {
			run(line, &errorHandler)
			errorHandler.HadError = false
		}
	}
}

func run(source string, errorHandler *langerr.Handler) {
	scanner := scanner.NewScanner(source, errorHandler)
	tokens := scanner.ScanTokens()
	parser := parser.NewParser(tokens, errorHandler)
	expr := parser.Parse()

	if errorHandler.HadError {
		return
	}

	var printer ast.Printer
	fmt.Println(printer.Print(expr))
}
