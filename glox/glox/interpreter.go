package main

import "fmt"

type RuntimeError struct {
	Message string
	Token   Token // optional, if you want context
}

func (e *RuntimeError) Error() string {
	return fmt.Sprintf("Runtime error at '%v': %s", e.Token.lexeme, e.Message)
}

type Interpreter struct{}

func (i *Interpreter) Interpret(expr Expr) {
	defer func() {
		if r := recover(); r != nil {
			if runtimeErr, ok := r.(*RuntimeError); ok {
				lox.runTimeError(*runtimeErr)
				// fmt.Println(runtimeErr.Error())
			} else {
				panic(r)
			}
		}
	}()

	fmt.Println(i.evaluate(expr))
}

func (i *Interpreter) VisitLiteralExpr(expr Literal) any {
	return expr.value
}

func (i *Interpreter) VisitGroupingExpr(expr Grouping) any {
	fmt.Println("visiting grouping")
	return i.evaluate(expr)
}

func (i *Interpreter) evaluate(expr Expr) any {
	return expr.Accept(i)
}

func (i *Interpreter) VisitUnaryExpr(expr Unary) any {
	fmt.Println("visiting unary")
	right := i.evaluate(expr.right)

	switch expr.operator.tokenType {
	case MINUS:
		n, _ := right.(float64)
		return -n
	case BANG:
		return !i.isTruthy(right)
	default:
		fmt.Println("unknown operator should be unreachable", right)
		return nil
	}
}

func (i *Interpreter) isTruthy(obj Object) bool {
	if obj == nil {
		return false
	}

	b, ok := obj.(bool)
	if ok {
		return b
	}

	return true
}

func (i *Interpreter) VisitBinaryExpr(expr Binary) any {
	fmt.Println("visiting binary")
	left := i.evaluate(expr.left)
	right := i.evaluate(expr.right)

	switch expr.operator.tokenType {
	case MINUS,
		SLASH,
		STAR,
		GREATER,
		GREATER_EQUAL,
		LESS,
		LESS_EQUAL:

		leftNumber, leftNumberOk := left.(float64)
		rightNumber, rightNumberOk := right.(float64)

		if !(leftNumberOk && rightNumberOk) {
			fmt.Println("oh no")
			return nil
		}

		switch expr.operator.tokenType {
		case MINUS:
			return leftNumber - rightNumber
		case SLASH:
			return leftNumber / rightNumber
		case STAR:
			return leftNumber * rightNumber
		case GREATER:
			return leftNumber > rightNumber
		case GREATER_EQUAL:
			return leftNumber >= rightNumber
		case LESS:
			return leftNumber < rightNumber
		case LESS_EQUAL:
			return leftNumber <= rightNumber
		default:
			return nil
		}
	case PLUS:

		leftString, leftStringOk := left.(string)
		rightString, rightStringOk := right.(string)

		if leftStringOk && rightStringOk {
			return leftString + rightString
		}

		leftNumber, leftNumberOk := left.(float64)
		rightNumber, rightNumberOk := right.(float64)

		if leftNumberOk && rightNumberOk {
			return leftNumber + rightNumber
		}

		fmt.Println("uh oh")
		return nil
	default:
		fmt.Println("unreachable unkonwn op")
		return nil
	}
}

func (i *Interpreter) VisitCommaExpr(expr Comma) any {
	fmt.Println("visiting comma")

	for index, e := range expr.exprs {
		res := i.evaluate(e)
		if index == len(expr.exprs)-1 {
			return res
		}
	}

	fmt.Println("unreachable")
	return nil
}

func (i *Interpreter) VisitTernaryExpr(expr Ternary) any {
	fmt.Println("visiting ternary")

	condition := i.evaluate(expr.condition)
	if i.isTruthy(condition) {
		return i.evaluate(expr.outcome1)
	}
	return i.evaluate(expr.outcome2)
}
