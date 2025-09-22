package main

type Environment struct {
	values map[string]Object
}

func NewEnvironment() *Environment {
	return &Environment{
		values: make(map[string]Object, 10),
	}
}

func (e *Environment) define(name string, value Object) {
	e.values[name] = value
}

func (e *Environment) get(name Token) Object {
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

func (e *Environment) assign(name Token, value Object) {
	_, ok := e.values[name.lexeme]
	if !ok {
		var err error = &RuntimeError{
			"Undefined variable '" + name.lexeme + "'.",
			name,
		}
		panic(err)
	}

	e.values[name.lexeme] = value
}
