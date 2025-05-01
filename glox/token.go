package main

import "fmt"

type Object any

type Token struct {
	tokenType TokenType
	lexeme    string
	object    Object
	line      int
}

func (t Token) toString() string {
	return fmt.Sprintf("%d %s %v", t.tokenType, t.lexeme, t.object)
}
