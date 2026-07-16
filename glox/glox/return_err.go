package main

import "fmt"

type ReturnError struct {
	Value any
}

func (e *ReturnError) Error() string {
	return fmt.Sprintf("Returning value %v", e.Value)
}
