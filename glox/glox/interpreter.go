package main

import "fmt"

type Interpreter struct{}

func (i Interpreter) VisitLiteralExpr(expr Literal) Object {
	return expr.value
}

func (i Interpreter) VisitGroupingExpr(expr Grouping) Object {
	return i.evaluate(expr)
}

func (i Interpreter) evaluate(expr Expr) Object {
	return expr.Accept(i)
}

func (i Interpreter) VisitUnaryExpr(expr Unary) Object {
	right := i.evaluate(expr)

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

func (i Interpreter) isTruthy(obj Object) bool {
	if obj == nil {
		return false
	}

	b, ok := obj.(bool)
	if ok {
		return b
	}

	return true
}

func (i Interpreter) VisitBinaryExpr(expr Binary) Object {
	left := i.evaluate(expr.left)
	right := i.evaluate(expr.right)

	switch expr.operator.tokenType {
	case MINUS:
		n_left, _ := left.(float64)
		n_right, _ := right.(float64)

		return n_left - n_right
	case SLASH:
		n_left, _ := left.(float64)
		n_right, _ := right.(float64)

		return n_left / n_right
	case STAR:
		n_left, _ := left.(float64)
		n_right, _ := right.(float64)

		return n_left * n_right
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
