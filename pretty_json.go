package mojilog

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fatih/color"
)

// PrettyJSONHandler formats logs as indented JSON with colors
type PrettyJSONHandler struct {
	out  io.Writer
	opts *slog.HandlerOptions
}

// NewPrettyJSONHandler creates a new pretty JSON handler
func NewPrettyJSONHandler(out io.Writer, opts *slog.HandlerOptions) *PrettyJSONHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	return &PrettyJSONHandler{
		out:  out,
		opts: opts,
	}
}

// Enabled implements slog.Handler
func (h *PrettyJSONHandler) Enabled(ctx context.Context, level slog.Level) bool {
	minLevel := slog.LevelInfo
	if h.opts.Level != nil {
		minLevel = h.opts.Level.Level()
	}
	return level >= minLevel
}

// Handle implements slog.Handler
func (h *PrettyJSONHandler) Handle(ctx context.Context, r slog.Record) error {
	// Use local time (KST)
	localTime := r.Time.Local()

	// Create JSON structure
	logEntry := make(map[string]interface{})

	// Basic fields
	logEntry["time"] = localTime.Format("2006-01-02 15:04:05.000")
	logEntry["level"] = r.Level.String()

	// Add emoji based on level or context
	emoji := getContextualEmoji(r.Message)
	if emoji == "" {
		emoji = getEmojiForLevel(r.Level)
	}
	if emoji != "" {
		logEntry["emoji"] = emoji
	}

	logEntry["msg"] = r.Message

	// Add source if requested
	if h.opts.AddSource && r.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		if f.File != "" {
			source := make(map[string]interface{})
			source["file"] = filepath.Base(f.File)
			source["line"] = f.Line

			funcName := filepath.Base(f.Function)
			if idx := strings.LastIndex(funcName, "."); idx != -1 {
				funcName = funcName[idx+1:]
			}
			source["function"] = funcName
			logEntry["source"] = source
		}
	}

	// Add attributes
	attrs := make(map[string]interface{})
	r.Attrs(func(a slog.Attr) bool {
		// Skip verbose attributes
		if !shouldSkipJSONAttr(a.Key) {
			// Handle special types
			switch v := a.Value.Any().(type) {
			case json.RawMessage:
				// Try to unmarshal as JSON
				var parsed interface{}
				if err := json.Unmarshal(v, &parsed); err == nil {
					attrs[a.Key] = parsed
				} else {
					attrs[a.Key] = string(v)
				}
			case []byte:
				// Try to parse as JSON first
				var parsed interface{}
				if err := json.Unmarshal(v, &parsed); err == nil {
					attrs[a.Key] = parsed
				} else {
					attrs[a.Key] = string(v)
				}
			case string:
				// Check if it looks like JSON
				if strings.HasPrefix(v, "{") || strings.HasPrefix(v, "[") {
					var parsed interface{}
					if err := json.Unmarshal([]byte(v), &parsed); err == nil {
						attrs[a.Key] = parsed
					} else {
						attrs[a.Key] = v
					}
				} else {
					attrs[a.Key] = v
				}
			default:
				attrs[a.Key] = a.Value.Any()
			}
		}
		return true
	})

	if len(attrs) > 0 {
		logEntry["attrs"] = attrs
	}

	// Marshal with indentation
	output, err := json.MarshalIndent(logEntry, "", "  ")
	if err != nil {
		return err
	}

	// Add color based on level
	var coloredOutput string
	switch r.Level {
	case slog.LevelError:
		coloredOutput = color.RedString(string(output))
	case slog.LevelWarn:
		coloredOutput = color.YellowString(string(output))
	case slog.LevelInfo:
		coloredOutput = color.GreenString(string(output))
	case slog.LevelDebug:
		coloredOutput = color.CyanString(string(output))
	default:
		coloredOutput = string(output)
	}

	_, err = h.out.Write([]byte(coloredOutput + "\n"))
	return err
}

// WithAttrs implements slog.Handler
func (h *PrettyJSONHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// For simplicity, we'll just return self
	// In production, you'd want to store and apply these
	return h
}

// WithGroup implements slog.Handler
func (h *PrettyJSONHandler) WithGroup(name string) slog.Handler {
	// For simplicity, we'll just return self
	return h
}

// shouldSkipJSONAttr determines if an attribute should be skipped in JSON output
func shouldSkipJSONAttr(key string) bool {
	// Skip less important attributes in JSON mode
	skipKeys := []string{
		"service",
		"version",
		"environment",
		"pid",
		"metric_name",
		"metric_value",
	}

	for _, skip := range skipKeys {
		if key == skip {
			return true
		}
	}
	return false
}

// SetupPrettyJSONLogger sets up a logger with pretty JSON formatting
func SetupPrettyJSONLogger(w io.Writer, level slog.Level, addSource bool) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: addSource,
	}

	handler := NewPrettyJSONHandler(w, opts)
	return slog.New(handler)
}