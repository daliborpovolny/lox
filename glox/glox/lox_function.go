package main

type LoxFunction struct {
	declaration *Function
}

func (f LoxFunction) call(interpreter Interpreter, arguments []any) (result any) {
	defer func() {
		if r := recover(); r != nil {
			if returnErr, ok := r.(*ReturnError); ok {
				result = returnErr.Value
			} else {
				panic(r)
			}
		}
	}()

	env := NewEnvironment(interpreter.globals)
	for i := range len(f.declaration.params) {
		env.initialize(f.declaration.params[i].lexeme)
		env.define(f.declaration.params[i].lexeme, arguments[i])
	}

	interpreter.executeBlock(f.declaration.body, env)

	return nil
}

func (f LoxFunction) arity() int {
	return len(f.declaration.params)
}

func (f LoxFunction) String() string {
	return "<fn " + f.declaration.name.lexeme + ">"
}
