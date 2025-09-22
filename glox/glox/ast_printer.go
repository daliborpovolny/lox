package main

import (
	"fmt"
	"strings"
)

type AstPrinter struct{}

func (a AstPrinter) Print(statements []Stmt) string {
	var builder strings.Builder
	for _, stmt := range statements {
		builder.WriteString(stmt.Accept(a).(string))
		builder.WriteString("\n")
	}
	return builder.String()
}

func (a AstPrinter) VisitExpressionStmt(stmt Expression) any {
	return a.parenthesize("expression", stmt.expression)
}

func (a AstPrinter) VisitPrintStmt(stmt Print) any {
	return a.parenthesize("print", stmt.expression)
}

func (a AstPrinter) VisitVarStmt(stmt Var) any {
	return a.parenthesize("var "+stmt.name.lexeme, stmt.initializer)
}

func (a AstPrinter) VisitVariableExpr(expr Variable) any {
	return expr.name.lexeme
}

func (a AstPrinter) VisitAssignExpr(expr Assign) any {
	return a.parenthesize("assign "+expr.name.lexeme, expr.value)
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
	var s strings.Builder
	s.WriteString("(")
	s.WriteString(name)
	for _, expr := range exprs {
		if expr != nil {
			s.WriteString(" ")
			s.WriteString(expr.Accept(a).(string))
		}
	}
	s.WriteString(")")
	return s.String()
}
