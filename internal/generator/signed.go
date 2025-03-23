package generator

import (
	"log/slog"
	"math"
	"math/rand"
	"reflect"
	"strconv"

	"golang.org/x/exp/constraints"
)

// SignedGenerator generates signed integers
// based on min and max values
type SignedGenerator[T constraints.Signed] struct {
	min int64
	max int64
	BaseGenerator
}

// NewSignedGenerator creates a new SignedGenerator using "min" and "max" from tags.
func NewSignedGenerator[T constraints.Signed](tags map[string]string, rand *rand.Rand, logger *slog.Logger) Generator[T] {
	var minVal, maxVal int64
	var zero T

	typeInfo := reflect.TypeOf(zero)
	var min, max int64

	switch typeInfo.Kind() {
	case reflect.Int:
		min, max = math.MinInt, math.MaxInt
	case reflect.Int8:
		min, max = math.MinInt8, math.MaxInt8
	case reflect.Int16:
		min, max = math.MinInt16, math.MaxInt16
	case reflect.Int32:
		min, max = math.MinInt32, math.MaxInt32
	case reflect.Int64:
		min, max = math.MinInt64, math.MaxInt64
	}

	logger.Debug("Creating SignedGenerator instance", "type", typeInfo.String(), "defaultMin", min, "defaultMax", max)
	minVal, maxVal = min, max

	if tags["min"] != "" {
		val, err := strconv.ParseInt(tags["min"], 10, 64)
		if err != nil {
			logger.Error("Failed to parse min tag", "min", tags["min"], "error", err)
			panic(err)
		}
		clampedVal := signedClamp(val, min, max)
		minVal = clampedVal
	}

	if tags["max"] != "" {
		val, err := strconv.ParseInt(tags["max"], 10, 64)
		if err != nil {
			logger.Error("Failed to parse max tag", "max", tags["max"], "error", err)
			panic(err)
		}
		clampedVal := signedClamp(val, min, max)
		maxVal = clampedVal
	}

	logger.Debug("SignedGenerator created", "min", minVal, "max", maxVal)
	return &SignedGenerator[T]{
		min: minVal,
		max: maxVal,
		BaseGenerator: BaseGenerator{
			rand,
			logger,
		},
	}
}

// Evaluate returns randomly generated signed integer within the given range.
func (g *SignedGenerator[T]) Evaluate() (T, error) {
	// special-case for full int64 range to avoid overflow
	if g.min == math.MinInt64 && g.max == math.MaxInt64 {
		// Use a full 64-bit random number.
		r := g.rand.Uint64()
		value := int64(r)
		g.logger.Debug("Evaluate generated value (full range)", "value", value)
		return T(value), nil
	}

	// calculate the range using unsigned arithmetic for smaller ranges
	rangeSize := uint64(g.max-g.min) + 1
	r := g.rand.Uint64() % rangeSize
	value := int64(r) + g.min
	g.logger.Debug("Evaluate generated value", "value", value)
	return T(value), nil
}

func signedClamp(value, min, max int64) int64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
