package utils

import (
	"bytes"
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewPlainHandler(t *testing.T) {
	var buf bytes.Buffer

	handler := NewPlainHandler(&buf, nil)
	assert.NotNil(t, handler)
	assert.Equal(t, handler.w, &buf)
	assert.Nil(t, handler.leveler)

	opts := &slog.HandlerOptions{Level: slog.LevelInfo}
	handler = NewPlainHandler(&buf, opts)
	assert.NotNil(t, handler)
	assert.Equal(t, handler.w, &buf)
	assert.Equal(t, handler.leveler, opts.Level)
}

func TestPlainHandler_Enabled(t *testing.T) {
	var buf bytes.Buffer

	handler := &plainHandler{w: &buf, leveler: nil}
	ctx := context.Background()
	assert.True(t, handler.Enabled(ctx, slog.LevelInfo))

	handler.leveler = slog.LevelInfo
	assert.False(t, handler.Enabled(ctx, slog.LevelDebug))
	assert.True(t, handler.Enabled(ctx, slog.LevelInfo))
}

func TestPlainHandler_Handle(t *testing.T) {
	var buf bytes.Buffer
	handler := &plainHandler{w: &buf, leveler: slog.LevelInfo}
	ctx := context.Background()
	record := slog.NewRecord(time.Now(), slog.LevelInfo, "test", 0)

	err := handler.Handle(ctx, record)
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "test")
}

func TestPlainHandler_WithAttrs(t *testing.T) {
	var buf bytes.Buffer
	handler := &plainHandler{w: &buf, leveler: slog.LevelInfo}
	result := handler.WithAttrs([]slog.Attr{})
	assert.Equal(t, handler, result)
}

func TestPlainHandler_WithGroup(t *testing.T) {
	var buf bytes.Buffer
	handler := &plainHandler{w: &buf, leveler: slog.LevelInfo}
	result := handler.WithGroup("group")
	assert.Equal(t, handler, result)
}

func TestInitLogger(t *testing.T) {
	InitLogger(false)
	InitLogger(true)
	assert.NotPanics(t, func() {
		InitLogger(false)
		InitLogger(true)
	})
}
