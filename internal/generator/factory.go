package generator

import (
	"log/slog"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"golang.org/x/exp/constraints"
)

type GeneratorFactory interface {
	Create(tags map[string]string, rand *rand.Rand, logger *slog.Logger) AnyGenerator
}

var GeneratorFactories = map[string]GeneratorFactory{
	"string":       StringFactory{},
	"&{time Time}": TimeFactory{},
	"&{uuid UUID}": UUIDFactory{},
	"int":          SignedFactory[int]{},
	"int8":         SignedFactory[int8]{},
	"int16":        SignedFactory[int16]{},
	"int32":        SignedFactory[int32]{},
	"int64":        SignedFactory[int64]{},
	"uint":         UnsignedFactory[uint]{},
	"uint8":        UnsignedFactory[uint8]{},
	"uint16":       UnsignedFactory[uint16]{},
	"uint32":       UnsignedFactory[uint32]{},
	"uint64":       UnsignedFactory[uint64]{},
	"float32":      FloatFactory[float32]{},
	"float64":      FloatFactory[float64]{},
}

type StringFactory struct{}

func (f StringFactory) Create(tags map[string]string, rand *rand.Rand, logger *slog.Logger) AnyGenerator {
	return &GenericGenerator[string]{impl: NewStringGenerator(tags, rand, logger)}
}

type TimeFactory struct{}

func (f TimeFactory) Create(tags map[string]string, rand *rand.Rand, logger *slog.Logger) AnyGenerator {
	return &GenericGenerator[time.Time]{impl: NewTimeGenerator(tags, rand, logger)}
}

type UUIDFactory struct{}

func (f UUIDFactory) Create(tags map[string]string, rand *rand.Rand, logger *slog.Logger) AnyGenerator {
	return &GenericGenerator[uuid.UUID]{impl: NewUUIDGenerator(rand, logger)}
}

type SignedFactory[T constraints.Signed] struct{}

func (f SignedFactory[T]) Create(tags map[string]string, rand *rand.Rand, logger *slog.Logger) AnyGenerator {
	return &GenericGenerator[T]{impl: NewSignedGenerator[T](tags, rand, logger)}
}

type UnsignedFactory[T constraints.Unsigned] struct{}

func (f UnsignedFactory[T]) Create(tags map[string]string, rand *rand.Rand, logger *slog.Logger) AnyGenerator {
	return &GenericGenerator[T]{impl: NewUnsignedGenerator[T](tags, rand, logger)}
}

type FloatFactory[T constraints.Float] struct{}

func (f FloatFactory[T]) Create(tags map[string]string, rand *rand.Rand, logger *slog.Logger) AnyGenerator {
	return &GenericGenerator[T]{impl: NewFloatGenerator[T](tags, rand, logger)}
}
