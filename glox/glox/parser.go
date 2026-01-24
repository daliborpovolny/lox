package main

import "fmt"

type ParseError struct {
	Message string
	Token   Token
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("Parse error at '%v': %s", e.Token.lexeme, e.Message)
}

type Parser struct {
	tokens  []Token
	current int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens,
		0,
	}
}

func (p *Parser) Parse() []Stmt {
	statements := make([]Stmt, 0, 10)
	for !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}

	return statements
}

func (p *Parser) declaration() Stmt {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(*ParseError); ok {
				// fmt.Println("synchronizing...")
				p.synchronize()
			} else {
				panic(r)
			}
		}
	}()

	if p.match(VAR) {
		return p.varDeclaration()
	}

	return p.statement()
}

func (p *Parser) varDeclaration() Stmt {
	var name Token = p.consume(IDENTIFIER, "Expect variable name.")

	var initializer Expr
	if p.match(EQUAL) {
		initializer = p.expression()
	}

	p.consume(SEMICOLON, "Expect ';' after variable declaration.")

	return Var{
		name,
		initializer,
	}
}

func (p *Parser) statement() Stmt {
	if p.match(PRINT) {
		return p.printStatement()
	}
	if p.match(LEFT_BRACE) {
		return Block{p.block()}
	}
	if p.match(IF) {
		return p.ifStatement()
	}
	return p.expressionStatement()
}

func (p *Parser) block() []Stmt {
	statements := []Stmt{}

	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}

	p.consume(RIGHT_BRACE, "Expect '}' after block")
	return statements
}

func (p *Parser) printStatement() Stmt {
	value := p.expression()
	p.consume(SEMICOLON, "Expect ';' after value")
	return Print{value}
}

func (p *Parser) expressionStatement() Stmt {
	expr := p.expression()
	p.consume(SEMICOLON, "Expect ';' after expression")
	return Expression{expr}
}

func (p *Parser) ifStatement() Stmt {

	p.consume(LEFT_PAREN, "Expect '(' after if")
	condition := p.expression()
	p.consume(RIGHT_PAREN, "Expect ')' after if")

	thenStmt := p.statement()

	var elseStmt Stmt
	if p.match(ELSE) {
		elseStmt = p.statement()
	}

	return If{
		condition,
		thenStmt,
		elseStmt,
	}
}

func (p *Parser) expression() Expr {
	return p.comma()
}

func (p *Parser) comma() Expr {
	expr := p.nonCommaExpression()
	if p.peek().tokenType == COMMA {
		commaExpr := Comma{
			[]Expr{expr},
		}
		for p.match(COMMA) {
			commaExpr.exprs = append(commaExpr.exprs, p.nonCommaExpression())
		}
		return commaExpr
	}
	return expr
}

func (p *Parser) nonCommaExpression() Expr {
	return p.assignment()
}

func (p *Parser) assignment() Expr {
	expr := p.ternary()
	if p.match(EQUAL) {
		equals := p.previous()
		value := p.assignment()

		if v, ok := expr.(Variable); ok {
			name := v.name
			return Assign{name, value}
		}
		p.error(equals, "Invalid assignment target")
	}

	return expr
}

func (p *Parser) ternary() Expr {
	expr := p.logicOr()

	if p.match(QUESTION_MARK) {
		outcome1 := p.expression()
		p.consume(COLON, "? denotes a ternary operator: expected expr ? expr : expr")
		outcome2 := p.ternary()
		return Ternary{
			condition: expr,
			outcome1:  outcome1,
			outcome2:  outcome2,
		}
	}
	return expr
}

func (p *Parser) logicOr() Expr {
	expr := p.logicAnd()

	for p.match(OR) {
		operator := p.previous()
		right := p.logicAnd()

		expr = Logical{expr, operator, right}
	}

	return expr
}

func (p *Parser) logicAnd() Expr {
	expr := p.equality()

	for p.match(AND) {
		operator := p.previous()
		right := p.equality()

		expr = Logical{expr, operator, right}
	}

	return expr
}

func (p *Parser) equality() Expr {
	expr := p.comparison()
	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = Binary{expr, operator, right}
	}
	return expr
}

func (p *Parser) comparison() Expr {
	expr := p.term()
	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = Binary{expr, operator, right}
	}
	return expr
}

func (p *Parser) term() Expr {
	expr := p.factor()
	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = Binary{expr, operator, right}
	}
	return expr
}

func (p *Parser) factor() Expr {
	expr := p.unary()
	for p.match(SLASH, STAR) {
		operator := p.previous()
		right := p.unary()
		expr = Binary{expr, operator, right}
	}
	return expr
}

func (p *Parser) unary() Expr {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right := p.unary()
		return Unary{operator, right}
	}
	return p.primary()
}

func (p *Parser) primary() Expr {
	if p.match(FALSE) {
		return Literal{false}
	} else if p.match(TRUE) {
		return Literal{true}
	} else if p.match(NIL) {
		return Literal{nil}
	}

	if p.match(NUMBER, STRING) {
		return Literal{p.previous().object}
	}

	if p.match(IDENTIFIER) {
		return Variable{p.previous()}
	}

	if p.match(LEFT_PAREN) {
		expr := p.expression()
		p.consume(RIGHT_PAREN, "Expect ')' after expression")
		return Grouping{expr}
	}

	p.error(p.peek(), "Expect expression.")
	panic("") // unreachable
}

func (p *Parser) consume(tokenType TokenType, message string) Token {
	if p.check(tokenType) {
		return p.advance()
	}

	p.error(p.peek(), message)
	return Token{}
}

func (p *Parser) error(token Token, message string) {
	var err = &ParseError{
		message,
		token,
	}

	lox.errorToken(token, message)
	panic(err)
}

func (p *Parser) synchronize() {
	p.advance()
	for !p.isAtEnd() {
		if p.previous().tokenType == SEMICOLON {
			return
		}

		switch p.peek().tokenType {
		case CLASS:
		case FOR:
		case FUN:
		case IF:
		case PRINT:
		case RETURN:
		case VAR:
		case WHILE:
			return
		}

		p.advance()
	}
}

func (p *Parser) match(tokenTypes ...TokenType) bool {
	for _, tokenT := range tokenTypes {
		if p.check(tokenT) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(tokenType TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().tokenType == tokenType
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current += 1
	}
	return p.previous()
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

func (p *Parser) isAtEnd() bool {
	return p.peek().tokenType == EOF
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}
