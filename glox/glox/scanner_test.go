package main

import (
	"testing"
)

func TestMultilineComments(t *testing.T) {
	program := "x+=2/* fsdfsdf" +
		"fsdfl;ksdf sdflsdkfj <3 !=" +
		"fsdfsdf */-3"

	l := Lox{}
	l.run(program)
}
