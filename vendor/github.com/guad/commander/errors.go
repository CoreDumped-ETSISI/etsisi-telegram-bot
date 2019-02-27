package commander

import (
	"fmt"
)

var (
	// ErrNotEnoughArgs means user did not pass enough args to the command
	ErrNotEnoughArgs = fmt.Errorf("Not enough arguments provided")
	// ErrTooManyArgs means user passed too many args to the command
	ErrTooManyArgs = fmt.Errorf("Too many arguments provided")
	// ErrArgSyntaxError means an argument failed to parse correctly.
	ErrArgSyntaxError = fmt.Errorf("Syntax error when parsing argument")
)
