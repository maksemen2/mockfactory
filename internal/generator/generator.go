package generator

import (
	"errors"
)

// Generator is an interface
// for a generator that can return a specific type.
type Generator[T any] interface {
	Evaluate() (T, error)
}

var ErrValidate = errors.New("failed to validate param")
