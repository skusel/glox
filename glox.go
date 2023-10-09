package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/skusel/glox/error"
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
	source, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	} else {
		var errorHandler error.Handler
		run(string(source), &errorHandler)
		if errorHandler.HadError {
			os.Exit(65)
		}
	}
}

func runPrompt() {
	var errorHandler error.Handler
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

func run(source string, errorHandler *error.Handler) {
	scanner := scanner.NewScanner(source, errorHandler)
	tokens := scanner.ScanTokens()

	for _, token := range tokens {
		fmt.Println(token)
	}
}
