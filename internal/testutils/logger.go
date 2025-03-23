package testutils

import (
	"io"
	"log/slog"
)

func TestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
}
