package main

type Environment struct {
	enclosing *Environment

	values      map[string]Object
	initialized map[string]bool
}

func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		enclosing:   enclosing,
		values:      make(map[string]Object, 10),
		initialized: make(map[string]bool, 10),
	}
}

func (e *Environment) initialize(name string) {
	e.initialized[name] = true
}

func (e *Environment) define(name string, value Object) {
	e.values[name] = value
}

func (e *Environment) get(name Token) Object {
	value, ok := e.values[name.lexeme]
	if !ok {
		if e.enclosing != nil {
			return e.enclosing.get(name)
		}

		var err error = &RuntimeError{
			"Undefined variable '" + name.lexeme + "'.",
			name,
		}
		panic(err)
	}

	if !e.initialized[name.lexeme] {
		var err error = &RuntimeError{
			"Uninitialized variable '" + name.lexeme + "'.",
			name,
		}
		panic(err)
	}

	return value
}

func (e *Environment) assign(name Token, value Object) {
	_, ok := e.values[name.lexeme]
	if !ok {

		if e.enclosing != nil {
			e.enclosing.assign(name, value)
			return
		}

		var err error = &RuntimeError{
			"Undefined variable '" + name.lexeme + "'.",
			name,
		}
		panic(err)
	}

	e.values[name.lexeme] = value
	e.initialized[name.lexeme] = true
}
