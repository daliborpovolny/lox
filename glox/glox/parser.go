package main

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
	defer func() {
		if r := recover(); r != nil {
			// fmt.Println(r)
		}
	}()

	statements := make([]Stmt, 0, 10)
	for !p.isAtEnd() {
		statements = append(statements, p.statement())
	}

	return statements
}

func (p *Parser) statement() Stmt {
	if p.match(PRINT) {
		return p.printStatement()
	}
	return p.expressionStatement()
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

func (p *Parser) expression() Expr {
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
	return p.ternary()
}

func (p *Parser) ternary() Expr {
	expr := p.equality()

	if p.match(QUESTION_MARK) {
		outcome1 := p.equality()
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

	if p.match(LEFT_PAREN) {
		expr := p.expression()
		p.consume(RIGHT_PAREN, "Expect ')' after expression")
		return Grouping{expr}
	}

	p.err(p.peek(), "Expect expression")
	panic("")
}

func (p *Parser) consume(tokenType TokenType, message string) Token {
	if p.check(tokenType) {
		return p.advance()
	}

	p.err(p.peek(), message)
	return Token{}
}

func (p *Parser) err(token Token, message string) {
	lox.errorToken(token, message)
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
