package main

import (
	"bufio"
	"fmt"
	"os"
)

type Lox struct {
	hadError bool
}

var lox Lox = Lox{}

func (l *Lox) Start(args []string) error {
	if len(args) > 1 {
		return fmt.Errorf("usage: golox [script]")
	} else if len(args) == 1 {
		err := l.runFile(args[0])
		if err != nil {
			return err
		}
	} else {
		l.runPrompt()
	}
	return nil
}

func (l *Lox) runFile(path string) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	l.run(string(bytes))

	if l.hadError {
		os.Exit(65)
	}

	return nil
}

func (l *Lox) runPrompt() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		l.run(input)
		l.hadError = false
	}
}

func (l *Lox) run(source string) {
	scanner := NewScanner(source)
	tokens := scanner.scanTokens()
	// for _, token := range tokens {
	// 	fmt.Println(token.toString())
	// }
	// fmt.Printf("scanned\nparsing...\n")

	parser := NewParser(tokens)
	expression := parser.Parse()
	if l.hadError {
		return
	}

	astPrinter := AstPrinter{}
	fmt.Println(astPrinter.Print(expression))
}

func (l *Lox) error(line int, message string) {
	l.report(line, "", message)
}

func (l *Lox) report(line int, where string, message string) {
	fmt.Fprintln(os.Stderr, "[line ", line, "] Error", where, ": ", message)
	l.hadError = true
}

func (l *Lox) errorToken(token Token, message string) {
	if token.tokenType == EOF {
		l.report(token.line, " at end", message)
	} else {
		l.report(token.line, "at '"+token.lexeme+"'", message)
	}
}

func main() {
	// // Represents: (1 + 2)
	// left := &Binary{
	// 	left:     &Literal{value: 1},
	// 	operator: Token{lexeme: "+"},
	// 	right:    &Literal{value: 2},
	// }

	// // Represents: (3 - 4)
	// right := &Binary{
	// 	left:     &Literal{value: 3},
	// 	operator: Token{lexeme: "-"},
	// 	right:    &Literal{value: 4},
	// }

	// // Represents: ( (1 + 2) * (3 - 4) )
	// expr := &Binary{
	// 	left:     left,
	// 	operator: Token{lexeme: "*"},
	// 	right:    right,
	// }

	// visitor := AstPrinter{}
	// fmt.Println(visitor.Print(expr))
	// os.Exit(0)

	err := lox.Start(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(64)
	}
}
