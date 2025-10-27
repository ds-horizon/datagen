package utils

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
)

type plainHandler struct {
	w       io.Writer
	leveler slog.Leveler
}

func NewPlainHandler(w io.Writer, opts *slog.HandlerOptions) *plainHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	return &plainHandler{w: w, leveler: opts.Level}
}

func (h *plainHandler) Enabled(ctx context.Context, level slog.Level) bool {
	if h.leveler == nil {
		return true
	}
	return level >= h.leveler.Level()
}

func (h *plainHandler) Handle(ctx context.Context, r slog.Record) error {
	ts := r.Time.Local().Format("2006-01-02T15:04:05")
	lvl := strings.ToUpper(r.Level.String())
	_, err := fmt.Fprintf(h.w, "%s %s %s\n", ts, lvl, r.Message)
	return err
}

func (h *plainHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *plainHandler) WithGroup(name string) slog.Handler {
	return h
}

func InitLogger(verbose bool) {
	logLevel := slog.LevelInfo
	if verbose {
		logLevel = slog.LevelDebug
	}

	slog.SetDefault(slog.New(NewPlainHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	})))
}
