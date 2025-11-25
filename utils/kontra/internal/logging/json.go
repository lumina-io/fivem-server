package logging

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"regexp"
)

// ANSI color code regex pattern - matches escape sequences
var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*m`)

type FileHandler struct {
	slog.Handler
	logger *log.Logger
}

// stripANSIColors removes ANSI color codes from text
func stripANSIColors(text string) string {
	return ansiRegex.ReplaceAllString(text, "")
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
			return "STDOUT"
		case slog.LevelWarn:
			return "WARNING"
		case slog.LevelError:
			return "STDERR"
		default:
			return record.Level.String()
		}
	}()

	// Strip ANSI color codes from the message for clean JSON output
	cleanMessage := stripANSIColors(record.Message)

	json, _ := json.Marshal(map[string]any{
		"time":    ts,
		"level":   level,
		"message": cleanMessage,
	})

	h.logger.Print(string(json))
	return nil
}
