package generator

import (
	"log/slog"
	"math/rand"
	"time"
)

// TimeGenerator generates time.Time objects
// Can generate past, future or current time
type TimeGenerator struct {
	isPast   bool
	isFuture bool
	BaseGenerator
}

// NewTimeGenerator creates a TimeGenerator based on the "range" tag.
// "range" tag can be "past" or "future". Default is current time.
func NewTimeGenerator(tags map[string]string, rand *rand.Rand, logger *slog.Logger) Generator[time.Time] {
	past := false
	future := false
	timeRange, ok := tags["range"]
	if ok {
		if timeRange == "past" {
			past = true
		} else if timeRange == "future" {
			future = true
		}
	}
	logger.Debug("TimeGenerator created", "past", past, "future", future)
	return &TimeGenerator{past, future, BaseGenerator{rand, logger}}
}

// Evaluate returns a random time based on the range.
// If range is not provided, it returns the current time.
func (g *TimeGenerator) Evaluate() (time.Time, error) {
	if g.isPast {
		g.logger.Debug("Generating a past time value")
		return time.Now().Add(-time.Duration(g.rand.Intn(1000000)) * time.Hour), nil
	} else if g.isFuture {
		g.logger.Debug("Generating a future time value")
		return time.Now().Add(time.Duration(g.rand.Intn(1000000)) * time.Hour), nil
	}
	g.logger.Debug("Returning current time")
	return time.Now(), nil
}

func (g *TimeGenerator) Validate() error {
	g.logger.Debug("TimeGenerator validation passed")
	return nil
}
