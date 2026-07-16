package main

type LoxCallable interface {
	call(interpreter Interpreter, arguments []any) any
	arity() int
	String() string
}
