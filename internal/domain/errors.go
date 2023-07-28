package domain

import "github.com/pkg/errors"

var (
	ErrUnimplemented     = errors.New("unimplemented")
	ErrNoDataFound       = errors.New("no data found")
	ErrMultipleDataFound = errors.New("multiple data found")
	ErrInvalidValue      = errors.New("invalid value")
)
