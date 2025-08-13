package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type Lox struct {
	hadError        bool
	hadRuntimeError bool
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

	if l.hadRuntimeError {
		os.Exit(70)
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

	parser := NewParser(tokens)
	expression := parser.Parse()
	if l.hadError {
		return
	}

	astPrinter := AstPrinter{}
	fmt.Println(astPrinter.Print(expression))

	intepreter := Interpreter{}
	intepreter.Interpret(expression)

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

func (l *Lox) runTimeError(rErr RuntimeError) {
	fmt.Fprintln(os.Stderr, rErr.Message+"\n[line + "+strconv.Itoa(rErr.Token.line)+"]")
	l.hadRuntimeError = false
}

func main() {
	err := lox.Start(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(64)
	}
}
