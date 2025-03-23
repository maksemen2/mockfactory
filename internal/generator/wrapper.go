package generator

// AnyGenerator is an interface for a generator
// that can return any type.
type AnyGenerator interface {
	EvaluateAny() (any, error)
}

// GenericGenerator is a wrapper for a generator
// that can return any data type.
type GenericGenerator[T any] struct {
	impl Generator[T]
}

// EvaluateAny evaluates the generator and returns the result
func (g *GenericGenerator[T]) EvaluateAny() (any, error) {
	return g.impl.Evaluate()
}
