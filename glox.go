package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/skusel/glox/interpreter"
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
		errorHandler := langerr.NewHandler()
		interpreter := interpreter.NewInterpreter(errorHandler)
		run(string(source), interpreter, errorHandler)
		if errorHandler.HadError {
			os.Exit(65)
		}
		if errorHandler.HadRuntimeError {
			os.Exit(70)
		}
	}
}

func runPrompt() {
	errorHandler := langerr.NewHandler()
	interpreter := interpreter.NewInterpreter(errorHandler)
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
		} else {
			run(line, interpreter, errorHandler)
			errorHandler.HadError = false
			errorHandler.HadRuntimeError = false
		}
	}
}

func run(source string, interpreter *interpreter.Interpreter, errorHandler *langerr.Handler) {
	scanner := scanner.NewScanner(source, errorHandler)
	tokens := scanner.ScanTokens()
	parser := parser.NewParser(tokens, errorHandler)
	expr := parser.Parse()

	if errorHandler.HadError {
		return
	}

	//var printer ast.Printer
	//fmt.Println(printer.Print(expr))

	interpreter.Interpret(expr)

	if errorHandler.HadRuntimeError {
		return
	}
}
