package dtos

import "fmt"

var ErrOutOfMonthRange OutOfMonthRangeError

type OutOfMonthRangeError struct{}

func (OutOfMonthRangeError) Error() string { return "out of month range" }

type UnknownMonthError struct{ Value string }

func (e *UnknownMonthError) Error() string { return fmt.Sprintf("unknown month: %q", e.Value) }
