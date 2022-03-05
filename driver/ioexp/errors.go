package ioexp

import "errors"

var (
	ErrWriteOnly = errors.New("io does not support reads")
	ErrReadOnly  = errors.New("io does not support writes")
)
