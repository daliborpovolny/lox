package main

import "fmt"

type Environment struct {
	values map[string]Object
}

func NewEnvironment() *Environment {
	return &Environment{
		values: make(map[string]Object, 10),
	}
}

func (e *Environment) define(name string, value Object) {
	fmt.Println("defining...", name)
	e.values[name] = value
	fmt.Println(e.values)
}

func (e *Environment) get(name Token) Object {
	fmt.Println(e.values, "looking for", name.lexeme)

	value, ok := e.values[name.lexeme]
	if !ok {
		var err error = &RuntimeError{
			"Undefined variable '" + name.lexeme + "'.",
			name,
		}
		panic(err)
	}
	return value
}
