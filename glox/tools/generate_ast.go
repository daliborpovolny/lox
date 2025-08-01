package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: generate_ast <output directory>")
		os.Exit(64)
	}

	outputDir := os.Args[1]
	err := defineAst(outputDir, "Expr", []string{
		"Binary		: Expr left, Token operator, Expr right",
		"Grouping	: Expr expression",
		"Literal	: Object value",
		"Unary 		: Token operator, Expr right",
		"Ternary	: Expr condition, Expr outcome1, Expr outcome2",
	})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func defineAst(outDir, baseName string, types []string) error {

	path := outDir + "/" + strings.ToLower(baseName) + ".go"
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintln(f, "package main")

	fmt.Fprintln(f, "")

	// visitor interface
	fmt.Fprintln(f, "type Visitor interface {")
	for _, t := range types {
		name := strings.TrimSpace(strings.Split(t, ":")[0])
		fmt.Fprintln(f, "Visit"+name+"Expr(expr "+name+") any")
	}
	fmt.Fprintln(f, "}")

	fmt.Fprintln(f, "")

	// expr interface
	fmt.Fprintln(f, "type "+baseName+" interface{")
	fmt.Fprintln(f, "Accept(visitor Visitor) any")
	fmt.Fprintln(f, "}")

	// structs and their functions
	for _, t := range types {
		parts := strings.Split(t, ":")
		name := strings.TrimSpace(parts[0])
		params := strings.Split(parts[1], ",")

		// struct
		fmt.Fprintln(f, "type "+name+" struct {")
		for _, param := range params {
			param = strings.TrimSpace(param)

			param_parts := strings.Split(param, " ")
			param_type := param_parts[0]
			param_name := param_parts[1]

			fmt.Fprintln(f, param_name+" "+param_type)
		}
		fmt.Fprintf(f, "}\n\n")

		// function
		fmt.Fprintln(f, "func ("+strings.ToLower(name)[:1]+" "+name+") Accept(visitor Visitor) any {")
		fmt.Fprintln(f, "return visitor.Visit"+name+"Expr("+strings.ToLower(name)[:1]+")")
		fmt.Fprintln(f, "}")

	}
	cmd := exec.Command("gofmt", "-w", path)
	err = cmd.Run()
	return err
}
