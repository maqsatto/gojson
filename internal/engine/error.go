package engine

import "errors"

var (
	ErrTableNotFound = errors.New("table not found")
	ErrInvalidOp     = errors.New("invalid operator")
	ErrTypeMismatch  = errors.New("type mismatch")
	ErrBadRequest    = errors.New("bad request")
)
