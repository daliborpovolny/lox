package main

type stmtVisitor interface {
	VisitExpressionStmt(stmt Expression) any
	VisitPrintStmt(stmt Print) any
}

type Stmt interface {
	Accept(visitor stmtVisitor) any
}
type Expression struct {
	expression Expr
}

func (e Expression) Accept(visitor stmtVisitor) any {
	return visitor.VisitExpressionStmt(e)
}

type Print struct {
	expression Expr
}

func (p Print) Accept(visitor stmtVisitor) any {
	return visitor.VisitPrintStmt(p)
}
