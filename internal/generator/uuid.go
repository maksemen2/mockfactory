package generator

import (
	"log/slog"
	"math/rand"

	"github.com/google/uuid"
)

// UUIDGenerator generates google UUIDs
type UUIDGenerator struct {
	BaseGenerator
}

// NewUUIDGenerator returns a new UUIDGenerator.
func NewUUIDGenerator(rand *rand.Rand, logger *slog.Logger) Generator[uuid.UUID] {
	logger.Debug("UUIDGenerator created")
	return &UUIDGenerator{BaseGenerator{rand, logger}}
}

// Evaluate generates a new UUID.
func (g *UUIDGenerator) Evaluate() (uuid.UUID, error) {
	g.logger.Debug("Generating a new UUID")
	id, err := uuid.NewRandomFromReader(g.rand)
	if err != nil {
		g.logger.Error("Failed to generate UUID", "error", err)
		return id, err
	}
	g.logger.Debug("UUID generated successfully", "uuid", id.String())
	return id, nil
}
