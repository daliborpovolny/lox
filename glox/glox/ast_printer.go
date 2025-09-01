package main

import (
	"fmt"
	"strings"
)

type AstPrinter struct {
}

func (a AstPrinter) Print(expr Expr) any {
	return expr.Accept(a).(string)
}

func (a AstPrinter) VisitVariableExpr(expr Variable) any {
	return a.parenthesize("variable", expr)
}

func (a AstPrinter) VisitCommaExpr(expr Comma) any {
	return a.parenthesize("comma", expr.exprs...)
}

func (a AstPrinter) VisitTernaryExpr(expr Ternary) any {
	return a.parenthesize("ternary", expr.condition, expr.outcome1, expr.outcome2)
}

func (a AstPrinter) VisitBinaryExpr(expr Binary) any {
	return a.parenthesize(expr.operator.lexeme, expr.left, expr.right)
}

func (a AstPrinter) VisitGroupingExpr(expr Grouping) any {
	return a.parenthesize("group", expr.expression)
}

func (a AstPrinter) VisitLiteralExpr(expr Literal) any {
	if expr.value == nil {
		return "nil"
	}
	return fmt.Sprint(expr.value)
}

func (a AstPrinter) VisitUnaryExpr(expr Unary) any {
	return a.parenthesize(expr.operator.lexeme, expr.right)
}

func (a AstPrinter) parenthesize(name string, exprs ...Expr) string {
	s := strings.Builder{}
	s.WriteString("(")
	s.WriteString(name)

	for _, expr := range exprs {
		s.WriteString(" ")
		s.WriteString(expr.Accept(a).(string))
	}

	s.WriteString(")")
	return s.String()
}
