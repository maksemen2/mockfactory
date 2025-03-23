package generator

import (
	"log/slog"
	"math"
	"math/rand"
	"reflect"
	"strconv"

	"golang.org/x/exp/constraints"
)

// UnsignedGenerator generates unsigned integers
// based on min and max values
type UnsignedGenerator[T constraints.Unsigned] struct {
	min uint64
	max uint64
	BaseGenerator
}

// NewUnsignedGenerator creates a new UnsignedGenerator using "min" and "max" from tags.
// Defaults to min and max values of data type if not provided.
func NewUnsignedGenerator[T constraints.Unsigned](tags map[string]string, rand *rand.Rand, logger *slog.Logger) Generator[T] {
	var minVal, maxVal uint64
	var zero T

	typeInfo := reflect.TypeOf(zero)
	var min, max uint64

	switch typeInfo.Kind() {
	case reflect.Uint:
		max = math.MaxUint
	case reflect.Uint8:
		max = math.MaxUint8
	case reflect.Uint16:
		max = math.MaxUint16
	case reflect.Uint32:
		max = math.MaxUint32
	case reflect.Uint64:
		max = math.MaxUint64
	}

	minVal, maxVal = min, max

	if tags["min"] != "" {
		val, err := strconv.ParseUint(tags["min"], 10, 64)
		if err != nil {
			logger.Error("Failed to parse min value", "min", tags["min"], "error", err)
			panic(err)
		}
		minVal = unsignedClamp(val, min, max)
	}

	if tags["max"] != "" {
		val, err := strconv.ParseUint(tags["max"], 10, 64)
		if err != nil {
			logger.Error("Failed to parse max value", "max", tags["max"], "error", err)
			panic(err)
		}
		maxVal = unsignedClamp(val, min, max)
		logger.Debug("Parsed and clamped max value", "max", maxVal)
	}

	logger.Debug("UnsignedGenerator created", "min", minVal, "max", maxVal)
	return &UnsignedGenerator[T]{
		min:           minVal,
		max:           maxVal,
		BaseGenerator: BaseGenerator{rand, logger},
	}
}

// Evaluate returns randomly generated unsigned integer
func (g *UnsignedGenerator[T]) Evaluate() (T, error) {
	value := T(g.min + g.rand.Uint64()%(g.max-g.min))
	g.logger.Debug("Evaluated unsigned value", "value", value, "min", g.min, "max", g.max)
	return value, nil
}

func unsignedClamp(value, min, max uint64) uint64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
