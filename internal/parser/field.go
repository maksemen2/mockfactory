package parser

import (
	"errors"
	"log/slog"
	"math/rand"
	"time"

	"github.com/maksemen2/mockfactory/internal/generator"
)

// StructField is a parsed field of an struct
type StructField struct {
	Name     string
	Type     string
	MockTags map[string]string // parsed "mock" tags
}

// ToGenerator converts a StructField to an AnyGenerator
func (f StructField) ToGenerator(seed int64, logger *slog.Logger) (generator.AnyGenerator, error) {
	logger.Debug("Converting StructField to Generator",
		"fieldName", f.Name,
		"fieldType", f.Type,
		"mockTags", f.MockTags,
	)
	factory, ok := generator.GeneratorFactories[f.Type]
	if !ok {
		logger.Error("Unknown generator type provided", "fieldType", f.Type)
		return nil, errors.New("unknown generator type provided: " + f.Type)
	}

	if seed == 0 {
		seed = time.Now().UnixNano()
		logger.Info("Seed not provided, using current time as seed", "seed", seed)
	}

	gen := factory.Create(f.MockTags, rand.New(rand.NewSource(seed)), logger)
	return gen, nil
}
