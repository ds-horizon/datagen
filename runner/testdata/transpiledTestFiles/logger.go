package main

import (
    "context"
    "fmt"
    "io"
    "log/slog"
    "os"
    "strings"
)

type __dgi_plainHandler struct {
    w       io.Writer
    leveler slog.Leveler
}

func __dgi_NewPlainHandler(w io.Writer, opts *slog.HandlerOptions) *__dgi_plainHandler {
    if opts == nil {
        opts = &slog.HandlerOptions{}
    }
    return &__dgi_plainHandler{w: w, leveler: opts.Level}
}

func (h *__dgi_plainHandler) Enabled(ctx context.Context, level slog.Level) bool {
    if h.leveler == nil {
        return true
    }
    return level >= h.leveler.Level()
}

func (h *__dgi_plainHandler) Handle(ctx context.Context, r slog.Record) error {
    ts := r.Time.Local().Format("2006-01-02T15:04:05")
    lvl := strings.ToUpper(r.Level.String())
    _, err := fmt.Fprintf(h.w, "%s %s %s\n", ts, lvl, r.Message)
    return err
}

func (h *__dgi_plainHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
    return h
}

func (h *__dgi_plainHandler) WithGroup(name string) slog.Handler {
    return h
}

func __dgi_InitLogger(verbose bool) {
    logLevel := slog.LevelInfo
    if verbose {
        logLevel = slog.LevelDebug
    }
    slog.SetDefault(slog.New(__dgi_NewPlainHandler(os.Stdout, &slog.HandlerOptions{
            Level: logLevel,
    })))
}