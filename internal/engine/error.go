package engine

import "errors"

var (
	ErrTableNotFound = errors.New("table not found")
	ErrInvalidOp     = errors.New("invalid operator")
)
