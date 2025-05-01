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
		return fmt.Errorf("Usage: golox [script]")
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
	fmt.Println(source)
	return
}

func (l *Lox) error(line int, message string) {
	l.report(line, "", message)
}

func (l *Lox) report(line int, where string, message string) {
	fmt.Fprintln(os.Stderr, "[line ", line, "] Error", where, ": ", message)
	l.hadError = true
}

func main() {
	l := Lox{}
	err := l.Start(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(64)
	}
}
