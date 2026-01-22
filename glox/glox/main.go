package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

type Lox struct {
	hadError        bool
	hadRuntimeError bool

	interpreter *Interpreter
}

func NewLox() *Lox {
	return &Lox{
		hadError:        false,
		hadRuntimeError: false,
		interpreter:     NewInterpreter(),
	}
}

var printParseTree bool = false

var lox *Lox = NewLox()

func (l *Lox) Start(args []string) error {
	if len(args) > 1 {
		return fmt.Errorf("usage: golox [script] | test")
	} else if len(args) == 1 {
		if args[0] == "test" {
			l.runTests()
			return nil
		}

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
	statements := parser.Parse()
	if l.hadError {
		return
	}
	// fmt.Println("succesfully parsed")

	if printParseTree {
		astPrinter := AstPrinter{}
		fmt.Println(astPrinter.Print(statements))
	}

	l.interpreter.Interpret(statements)
	// fmt.Println("succesfully interpreted")

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
	fmt.Fprintln(os.Stderr, rErr.Message+"\n[line "+strconv.Itoa(rErr.Token.line)+"]")
	l.hadRuntimeError = false
}

// TESTFILES is a list of strings of test file names, for a test to work
// it has a have a {name}.lox and {name}.out in the tests folder
// the {name}.out file is what the output will be compared agains
// the {name}.out file should end with a new line
var TESTFILES []string = []string{
	"hello",
	"scope",
	"math",
	"precedance",
	"comments",
}

// runs the test included in TESTFILES
func (l *Lox) runTests() bool {

	exec.Command("go", "build").Run()

	for _, path := range TESTFILES {
		cmd := exec.Command("./glox", "../tests/"+path+".lox")

		outputBytes, err := cmd.Output()
		if err != nil {
			fmt.Println("error during test of", path, ":", err)
			return false
		}

		desiredBytes, err := os.ReadFile("../tests/" + path + ".out")

		output := string(outputBytes)
		desired := string(desiredBytes)

		if output != desired {
			fmt.Println("test", path, "failed")
			fmt.Println("\t expected:\n" + desired)
			fmt.Println("\t actual:\n" + output)
		} else {
			fmt.Println("test", path, "passed")
		}
	}

	return true
}

func main() {
	err := lox.Start(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(64)
	}
}
