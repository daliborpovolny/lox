package main

type stmtVisitor interface {
	VisitBlockStmt(stmt Block) any
	VisitExpressionStmt(stmt Expression) any
	VisitPrintStmt(stmt Print) any
	VisitVarStmt(stmt Var) any
	VisitIfStmt(stmt If) any
	VisitWhileStmt(stmt While) any
}

type Stmt interface {
	Accept(visitor stmtVisitor) any
}
type Block struct {
	statements []Stmt
}

func (b Block) Accept(visitor stmtVisitor) any {
	return visitor.VisitBlockStmt(b)
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

type Var struct {
	name        Token
	initializer Expr
}

func (v Var) Accept(visitor stmtVisitor) any {
	return visitor.VisitVarStmt(v)
}

type If struct {
	condition  Expr
	thenBranch Stmt
	elseBranch Stmt
}

func (i If) Accept(visitor stmtVisitor) any {
	return visitor.VisitIfStmt(i)
}

type While struct {
	condition Expr
	body      Stmt
}

func (w While) Accept(visitor stmtVisitor) any {
	return visitor.VisitWhileStmt(w)
}
