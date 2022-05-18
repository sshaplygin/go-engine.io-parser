package go_engine_io_parser

import (
	"errors"
	"fmt"
)

var (
	//
	errPaused = errors.New("paused")

	//
	errTimeout = errors.New("timeout")

	//
	errInvalidPayload = errors.New("invalid payload")

	//
	errOverlap = errors.New("overlap")
)

// OperationError is operation error.
type OperationError struct {
	Operation string
	Errors    error
}

func newOperationError(operation string, err error) error {
	return &OperationError{
		Operation: operation,
		Errors:    err,
	}
}

// Error implemented Error interface.
func (e *OperationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Operation, e.Errors.Error())
}

// Temporary returns true if error can retry.
func (e *OperationError) Temporary() bool {
	return e.Errors != nil
}
