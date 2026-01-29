package main

import (
	"fmt"
	"strconv"
	"strings"
)

type RuntimeError struct {
	Message string
	Token   Token
}

func (e *RuntimeError) Error() string {
	return fmt.Sprintf("Runtime error at '%v': %s", e.Token.lexeme, e.Message)
}

type Interpreter struct {
	environment *Environment
}

func (i *Interpreter) Interpret(statements []Stmt) {
	defer func() {
		if r := recover(); r != nil {
			if runtimeErr, ok := r.(*RuntimeError); ok {
				lox.runTimeError(*runtimeErr)
			} else {
				panic(r)
			}
		}
	}()

	for _, stmt := range statements {
		i.execute(stmt)
	}
}

// prints out the value of expression statements after executing them
func (i *Interpreter) ReplInterpret(statements []Stmt) {
	defer func() {
		if r := recover(); r != nil {
			if runTimeErr, ok := r.(*RuntimeError); ok {
				lox.runTimeError(*runTimeErr)
			} else {
				panic(r)
			}
		}
	}()

	for _, stmt := range statements {
		exprStmt, ok := stmt.(Expression)
		if ok {
			value := i.evaluate(exprStmt.expression)
			fmt.Println(value)
		} else {
			i.execute(stmt)
		}

	}
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		environment: NewEnvironment(nil),
	}
}

func (i *Interpreter) execute(stmt Stmt) {
	stmt.Accept(i)
}

func (i *Interpreter) VisitBlockStmt(stmt Block) any {
	i.executeBlock(stmt.statements, NewEnvironment(i.environment))
	return nil
}

func (i *Interpreter) executeBlock(stmts []Stmt, env *Environment) {

	oldEnv := i.environment

	i.environment = env

	for _, stmt := range stmts {
		i.execute(stmt)
	}

	i.environment = oldEnv
}

func (i *Interpreter) VisitVarStmt(stmt Var) any {
	var value Object

	if stmt.initializer != nil {
		value = i.evaluate(stmt.initializer)
		i.environment.initialize(stmt.name.lexeme)
	}

	i.environment.define(stmt.name.lexeme, value)
	return nil
}

func (i *Interpreter) VisitExpressionStmt(stmt Expression) any {
	i.evaluate(stmt.expression)
	return nil
}

func (i *Interpreter) VisitPrintStmt(stmt Print) any {
	value := i.evaluate(stmt.expression)
	fmt.Println(value)
	return nil
}

func (i *Interpreter) VisitIfStmt(stmt If) any {
	if i.isTruthy(i.evaluate(stmt.condition)) {
		i.execute(stmt.thenBranch)
	} else if stmt.elseBranch != nil {
		i.execute(stmt.elseBranch)
	}

	return nil
}

func (i *Interpreter) VisitWhileStmt(stmt While) any {
	for i.isTruthy(i.evaluate(stmt.condition)) {
		i.execute(stmt.body)
	}

	return nil
}

func (i *Interpreter) VisitVariableExpr(expr Variable) any {
	return i.environment.get(expr.name)
}

func (i *Interpreter) VisitLiteralExpr(expr Literal) any {
	return expr.value
}

func (i *Interpreter) VisitGroupingExpr(expr Grouping) any {
	return i.evaluate(expr.expression)
}

func (i *Interpreter) evaluate(expr Expr) any {
	return expr.Accept(i)
}

func (i *Interpreter) VisitAssignExpr(expr Assign) any {
	value := i.evaluate(expr.value)
	i.environment.assign(expr.name, value)
	return value
}

func (i *Interpreter) VisitLogicalExpr(expr Logical) any {

	left := i.evaluate(expr.left)

	switch expr.operator.tokenType {
	case OR:
		if i.isTruthy(left) {
			return left
		}
	case AND:
		if !i.isTruthy(left) {
			return left
		}
	}

	return i.evaluate(expr.right)

}

func (i *Interpreter) VisitUnaryExpr(expr Unary) any {

	right := i.evaluate(expr.right)

	switch expr.operator.tokenType {
	case MINUS:
		n, ok := right.(float64)
		if !ok {
			i.raiseNumberOperands(expr.operator, right)
		}

		return -n
	case BANG:
		return !i.isTruthy(right)
	default:
		fmt.Println("unknown operator should be unreachable", right)
		return nil
	}
}

func (i *Interpreter) raiseNumberOperands(operator Token, operands ...any) {
	var message string
	if len(operands) == 1 {
		message = "Operand must be a number."
	} else {
		message = "Operands must be numbers."
	}

	var err error = &RuntimeError{
		message,
		operator,
	}

	panic(err)
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

func (i *Interpreter) isEqual(a Object, b Object) bool { //todo simplify, this is go not java
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}
	return a == b
}

func (i *Interpreter) VisitBinaryExpr(expr Binary) any {

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
			i.raiseNumberOperands(expr.operator, left, right)
			return nil
		}

		switch expr.operator.tokenType {
		case MINUS:
			return leftNumber - rightNumber
		case SLASH:
			if rightNumber == 0 {
				var err error = &RuntimeError{
					"Cannot divide by zero.",
					expr.operator,
				}
				panic(err)
			}

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

		if leftStringOk && rightNumberOk {
			rightString = strconv.FormatFloat(rightNumber, 'f', 6, 64)
			rightString = strings.TrimSuffix(rightString, ".000000")
			return leftString + rightString
		}

		if rightStringOk && leftNumberOk {
			leftString = strconv.FormatFloat(leftNumber, 'f', 6, 64)
			leftString = strings.TrimSuffix(leftString, ".000000")
			return leftString + rightString
		}

		var err error = &RuntimeError{
			"Operands must be two numbers or strings and a number.",
			expr.operator,
		}
		panic(err)
	case BANG_EQUAL:
		return !i.isEqual(left, right)
	case EQUAL_EQUAL:
		return i.isEqual(left, right)
	default:
		fmt.Println("Unreachable unknown operator:", expr.operator.lexeme)
		var err RuntimeError = RuntimeError{
			"Unknown operator, should have failed in parsing.",
			expr.operator,
		}
		panic(err)
	}
}

func (i *Interpreter) VisitCommaExpr(expr Comma) any {
	// fmt.Println("visiting comma")

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
	// fmt.Println("visiting ternary")

	condition := i.evaluate(expr.condition)
	if i.isTruthy(condition) {
		return i.evaluate(expr.outcome1)
	}
	return i.evaluate(expr.outcome2)
}
