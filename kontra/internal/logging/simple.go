package logging

import (
	"context"
	"io"
	"log"
	"log/slog"

	"github.com/fatih/color"
)

type SimpleHandler struct {
	slog.Handler
	logger *log.Logger
}

func NewSimpleHandler(out io.Writer) *SimpleHandler {
	prefix := ""
	h := &SimpleHandler{
		Handler: slog.NewJSONHandler(out, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}),
		logger: log.New(out, prefix, 0),
	}
	return h
}

func (h *SimpleHandler) Handle(_ context.Context, record slog.Record) error {
	ts := record.Time.Format("2006-01-02 15:04:05")
	level := func() string {
		switch record.Level {
		case slog.LevelInfo:
			return color.GreenString("STDOUT")
		case slog.LevelWarn:
			return color.YellowString("WARNING")
		case slog.LevelError:
			return color.RedString("STDERR")
		default:
			return record.Level.String()
		}
	}()

	h.logger.Printf("%s %s: %s", color.CyanString(ts), level, record.Message)
	return nil
}
