package main

type Visitor interface {
	VisitBinaryExpr(expr Binary) any
	VisitGroupingExpr(expr Grouping) any
	VisitLiteralExpr(expr Literal) any
	VisitUnaryExpr(expr Unary) any
	VisitTernaryExpr(expr Ternary) any
}

type Expr interface {
	Accept(visitor Visitor) any
}
type Binary struct {
	left     Expr
	operator Token
	right    Expr
}

func (b Binary) Accept(visitor Visitor) any {
	return visitor.VisitBinaryExpr(b)
}

type Grouping struct {
	expression Expr
}

func (g Grouping) Accept(visitor Visitor) any {
	return visitor.VisitGroupingExpr(g)
}

type Literal struct {
	value Object
}

func (l Literal) Accept(visitor Visitor) any {
	return visitor.VisitLiteralExpr(l)
}

type Unary struct {
	operator Token
	right    Expr
}

func (u Unary) Accept(visitor Visitor) any {
	return visitor.VisitUnaryExpr(u)
}

type Ternary struct {
	condition Expr
	outcome1  Expr
	outcome2  Expr
}

func (t Ternary) Accept(visitor Visitor) any {
	return visitor.VisitTernaryExpr(t)
}
