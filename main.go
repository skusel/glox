package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/skusel/glox/lang"
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
		errorHandler := lang.NewErrorHandler()
		interpreter := lang.NewInterpreter(errorHandler)
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
	errorHandler := lang.NewErrorHandler()
	interpreter := lang.NewInterpreter(errorHandler)
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

func run(source string, interpreter *lang.Interpreter, errorHandler *lang.ErrorHandler) {
	scanner := lang.NewScanner(source, errorHandler)
	tokens := scanner.ScanTokens()
	parser := lang.NewParser(tokens, errorHandler)
	statements := parser.Parse()

	if errorHandler.HadError {
		return
	}

	interpreter.Interpret(statements)

	if errorHandler.HadRuntimeError {
		return
	}
}
