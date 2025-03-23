package generator

import (
	"fmt"
	"log/slog"
	"math/rand"
	"strconv"
)

const defaultLength = 8
const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"

// StringGenerator generates random strings
// based on the given length, prefix and suffix.
type StringGenerator struct {
	len    int
	prefix string
	suffix string
	BaseGenerator
}

// NewStringGenerator creates a new StringGenerator using.
// "prefix", "suffix" and "len" tags. Default length is 8 if not provided.
func NewStringGenerator(tags map[string]string, rand *rand.Rand, logger *slog.Logger) Generator[string] {
	var length int
	lenVal, ok := tags["len"]
	if ok {
		l, err := strconv.Atoi(lenVal)
		if err != nil {
			logger.Error("Failed to parse len tag", "len", lenVal, "error", err)
			panic(err)
		}
		if l < 1 {
			logger.Error("Invalid length provided", "length", l)
			panic(fmt.Sprintf("invalid length provided: %d", l))
		}
		length = l
	} else {
		length = defaultLength
	}

	prefix, ok := tags["prefix"]
	if !ok {
		prefix = ""
	}
	suffix, ok := tags["suffix"]
	if !ok {
		suffix = ""
	}

	logger.Debug("StringGenerator created", "length", length, "prefix", prefix, "suffix", suffix)
	return &StringGenerator{length, prefix, suffix, BaseGenerator{rand, logger}}
}

// Evaluate returns randomly generated string with given length including prefix and suffix.
func (g *StringGenerator) Evaluate() (string, error) {
	result := g.prefix + generateRandomString(g.len, g.rand) + g.suffix
	g.logger.Debug("Evaluate generated value", "value", result)
	return result, nil
}

func generateRandomString(n int, rand *rand.Rand) string {
	if n <= 0 {
		return ""
	}
	result := make([]byte, n)

	for i := 0; i < n; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}

	return string(result)
}

func (g *StringGenerator) Validate() error {
	return nil
}
