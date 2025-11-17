package usecases

import "fmt"

type MalformedRecordsError struct{ ActualColumns int }

func (err *MalformedRecordsError) Error() string {
	return fmt.Sprintf("expected at least 2 columns but got %d", err.ActualColumns)
}

type ParseTSVError struct {
	Err  error
	Line int
}

func (err *ParseTSVError) Error() string {
	return fmt.Sprintf("parse error at line %d: %s", err.Line, err.Err)
}

func (err *ParseTSVError) Unwrap() error { return err.Err }
