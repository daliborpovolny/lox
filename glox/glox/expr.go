package main

type exprVisitor interface {
	VisitBinaryExpr(expr Binary) any
	VisitGroupingExpr(expr Grouping) any
	VisitLiteralExpr(expr Literal) any
	VisitUnaryExpr(expr Unary) any
	VisitTernaryExpr(expr Ternary) any
	VisitCommaExpr(expr Comma) any
	VisitVariableExpr(expr Variable) any
}

type Expr interface {
	Accept(visitor exprVisitor) any
}
type Binary struct {
	left     Expr
	operator Token
	right    Expr
}

func (b Binary) Accept(visitor exprVisitor) any {
	return visitor.VisitBinaryExpr(b)
}

type Grouping struct {
	expression Expr
}

func (g Grouping) Accept(visitor exprVisitor) any {
	return visitor.VisitGroupingExpr(g)
}

type Literal struct {
	value Object
}

func (l Literal) Accept(visitor exprVisitor) any {
	return visitor.VisitLiteralExpr(l)
}

type Unary struct {
	operator Token
	right    Expr
}

func (u Unary) Accept(visitor exprVisitor) any {
	return visitor.VisitUnaryExpr(u)
}

type Ternary struct {
	condition Expr
	outcome1  Expr
	outcome2  Expr
}

func (t Ternary) Accept(visitor exprVisitor) any {
	return visitor.VisitTernaryExpr(t)
}

type Comma struct {
	exprs []Expr
}

func (c Comma) Accept(visitor exprVisitor) any {
	return visitor.VisitCommaExpr(c)
}

type Variable struct {
	name Token
}

func (v Variable) Accept(visitor exprVisitor) any {
	return visitor.VisitVariableExpr(v)
}
