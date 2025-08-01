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

func (p *Parser) Parse() (expr Expr) {
	defer func() {
		if r := recover(); r != nil {
			expr = nil
		}
	}()

	return p.expression()
}

func (p *Parser) expression() Expr {
	expr := p.nonCommaExpression()
	if p.peek().tokenType == COMMA {
		commaExpr := Comma{
			make([]Expr, 0),
		}

		for p.match(COMMA) {
			commaExpr.exprs = append(commaExpr.exprs, p.nonCommaExpression())
		}
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
		outcome2 := p.equality()
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
	panic("Expect expression")
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
