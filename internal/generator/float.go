package generator

import (
	"log/slog"
	"math"
	"math/rand"
	"reflect"
	"strconv"

	"golang.org/x/exp/constraints"
)

const (
	maxSafeFloat = 1e307
	minSafeFloat = -1e307
)

// FloatGenerator generates floating numbers
// based on min and max values
type FloatGenerator[T constraints.Float] struct {
	min float64
	max float64
	BaseGenerator
}

// NewFloatGenerator creates a new FloatGenerator using "min" and "max" from tags.
// Defaults to min and max values of data type if not provided.
func NewFloatGenerator[T constraints.Float](tags map[string]string, rand *rand.Rand, logger *slog.Logger) Generator[T] {
	var minVal, maxVal float64
	var zero T

	typeInfo := reflect.TypeOf(zero)
	var min, max float64

	switch typeInfo.Kind() {
	case reflect.Float32:
		min, max = -math.MaxFloat32, math.MaxFloat32
	case reflect.Float64:
		// Use safe defaults for float64 to avoid overflow issues when generating values.
		min, max = minSafeFloat, maxSafeFloat
	}

	logger.Debug("Creating FloatGenerator instance", "type", typeInfo.String(), "defaultMin", min, "defaultMax", max)

	minVal, maxVal = min, max

	if tags["min"] != "" {
		val, err := strconv.ParseFloat(tags["min"], 64)
		if err != nil {
			logger.Error("Failed to parse min tag", "min", tags["min"], "error", err)
			panic(err)
		}
		clampedVal := floatClamp(val, min, max)
		minVal = clampedVal
	}

	if tags["max"] != "" {
		val, err := strconv.ParseFloat(tags["max"], 64)
		if err != nil {
			logger.Error("Failed to parse max tag", "max", tags["max"], "error", err)
			panic(err)
		}
		clampedVal := floatClamp(val, min, max)
		maxVal = clampedVal
	}

	logger.Debug("FloatGenerator created", "min", minVal, "max", maxVal)
	return &FloatGenerator[T]{
		min: minVal,
		max: maxVal,
		BaseGenerator: BaseGenerator{
			rand,
			logger,
		},
	}
}

// Evaluate returns a random floating number between min and max values
func (g *FloatGenerator[T]) Evaluate() (T, error) {
	value := g.min + g.rand.Float64()*(g.max-g.min)
	g.logger.Debug("Evaluate generated value", "value", value)
	return T(value), nil
}

func floatClamp(val, min, max float64) float64 {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}
