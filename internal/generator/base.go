package generator

import (
	"log/slog"
	"math/rand"
)

// BaseGenerator contains base fields for all generators
type BaseGenerator struct {
	rand   *rand.Rand // for results reproducibility
	logger *slog.Logger
}
