package logging

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"log/slog"

	"github.com/fatih/color"
)

type FileHandler struct {
	slog.Handler
	logger *log.Logger
}

func NewFileHandler(out io.Writer, level slog.Level) *FileHandler {
	prefix := ""
	h := &FileHandler{
		Handler: slog.NewJSONHandler(out, &slog.HandlerOptions{
			Level: level,
		}),
		logger: log.New(out, prefix, 0),
	}
	return h
}

func (h *FileHandler) Handle(_ context.Context, record slog.Record) error {
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

	json, _ := json.Marshal(map[string]any{
		"time":    ts,
		"level":   level,
		"message": record.Message,
	})

	h.logger.Print(string(json))
	return nil
}
