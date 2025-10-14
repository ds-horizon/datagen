package utils

import (
	"log/slog"
	"os"
)

func InitLogger(verbose bool) {
	logLevel := slog.LevelWarn
	if verbose {
		logLevel = slog.LevelDebug
	}

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	})))
}
