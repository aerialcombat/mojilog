package mojilog

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/fatih/color"
)

// PrettyHandler is a custom handler that formats logs in a pretty way with colors
type PrettyHandler struct {
	opts      *slog.HandlerOptions
	mu        sync.Mutex
	out       io.Writer
	attrs     []slog.Attr
	groups    []string
	showEmoji bool
}

// Color functions for different levels
var (
	debugColor = color.New(color.FgCyan).SprintFunc()
	infoColor  = color.New(color.FgGreen).SprintFunc()
	warnColor  = color.New(color.FgYellow).SprintFunc()
	errorColor = color.New(color.FgRed, color.Bold).SprintFunc()
	fatalColor = color.New(color.FgRed, color.Bold, color.BgWhite).SprintFunc()

	timeColor     = color.New(color.FgHiBlack).SprintFunc()
	fileColor     = color.New(color.FgBlue).SprintFunc()
	functionColor = color.New(color.FgBlue).SprintFunc()
	attrColor     = color.New(color.FgMagenta).SprintFunc()
	resetColor    = color.New(color.Reset).SprintFunc()
)

// NewPrettyHandler creates a new pretty handler
func NewPrettyHandler(out io.Writer, opts *slog.HandlerOptions) *PrettyHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	return &PrettyHandler{
		out:       out,
		opts:      opts,
		showEmoji: true,
	}
}

// Enabled implements slog.Handler
func (h *PrettyHandler) Enabled(ctx context.Context, level slog.Level) bool {
	minLevel := slog.LevelInfo
	if h.opts.Level != nil {
		minLevel = h.opts.Level.Level()
	}
	return level >= minLevel
}

// Handle implements slog.Handler
func (h *PrettyHandler) Handle(ctx context.Context, r slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Use local time (which is set to KST via TZ environment variable)
	localTime := r.Time.Local()

	// Format timestamp (short format)
	timestamp := localTime.Format("15:04:05.0")

	// Get level and color
	levelStr := h.formatLevel(r.Level)

	// Get source info if requested
	source := ""
	if h.opts.AddSource && r.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		if f.File != "" {
			// Extract just filename and function name
			file := filepath.Base(f.File)
			funcName := filepath.Base(f.Function) + "()"
			// Remove package prefix from function name
			if idx := strings.LastIndex(funcName, "."); idx != -1 {
				funcName = funcName[idx+1:]
			}
			source = fmt.Sprintf("%s:%s:%d",
				fileColor(file),
				functionColor(funcName),
				f.Line)
		}
	}

	// Get emoji if contextual
	emoji := ""
	if h.showEmoji {
		contextEmoji := getContextualEmoji(r.Message)
		if contextEmoji != "" {
			emoji = contextEmoji
		} else {
			emoji = getEmojiForLevel(r.Level)
		}
		if emoji != "" {
			spacing := getEmojiSpacing(emoji)
			emoji = emoji + spacing
		}
	}

	// Format the main message
	var msg strings.Builder
	msg.WriteString(timeColor(timestamp))
	msg.WriteString(" ")
	msg.WriteString(levelStr)
	msg.WriteString(" ")

	if source != "" {
		msg.WriteString(source)
	}

	msg.WriteString(" ")
	msg.WriteString(emoji)
	msg.WriteString(r.Message)

	// Add attributes
	attrs := h.formatAttrs(r)
	if attrs != "" {
		msg.WriteString(" ")
		msg.WriteString(attrColor(attrs))
	}

	msg.WriteString("\n")

	_, err := h.out.Write([]byte(msg.String()))
	return err
}

// formatLevel returns a colored level string
func (h *PrettyHandler) formatLevel(level slog.Level) string {
	switch {
	case level >= slog.LevelError:
		return errorColor("ERROR")
	case level >= slog.LevelWarn:
		return warnColor(" WARN")
	case level >= slog.LevelInfo:
		return infoColor(" INFO")
	case level >= slog.LevelDebug:
		return debugColor("DEBUG")
	default:
		return "TRACE"
	}
}

// formatAttrs formats attributes as key=value pairs
func (h *PrettyHandler) formatAttrs(r slog.Record) string {
	var attrs []string

	// Add handler's attributes
	for _, attr := range h.attrs {
		if attr.Key != "" && !h.shouldSkipAttr(attr.Key) {
			attrs = append(attrs, fmt.Sprintf("%s=%v", attr.Key, attr.Value))
		}
	}

	// Add record's attributes
	r.Attrs(func(a slog.Attr) bool {
		if a.Key != "" && !h.shouldSkipAttr(a.Key) {
			attrs = append(attrs, fmt.Sprintf("%s=%v", a.Key, a.Value))
		}
		return true
	})

	if len(attrs) == 0 {
		return ""
	}

	return strings.Join(attrs, " ")
}

// shouldSkipAttr determines if an attribute should be skipped
func (h *PrettyHandler) shouldSkipAttr(key string) bool {
	// Skip common verbose attributes
	skipKeys := []string{
		"service",
		"version",
		"metric_name",
		"metric_value",
		"environment",
		"pid",
	}

	for _, skip := range skipKeys {
		if key == skip {
			return true
		}
	}
	return false
}

// WithAttrs implements slog.Handler
func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &PrettyHandler{
		out:       h.out,
		opts:      h.opts,
		attrs:     append(h.attrs, attrs...),
		groups:    h.groups,
		showEmoji: h.showEmoji,
	}
}

// WithGroup implements slog.Handler
func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	return &PrettyHandler{
		out:       h.out,
		opts:      h.opts,
		attrs:     h.attrs,
		groups:    append(h.groups, name),
		showEmoji: h.showEmoji,
	}
}

// SetupPrettyLogger sets up a logger with pretty formatting
func SetupPrettyLogger(w io.Writer, level slog.Level, addSource bool) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: addSource,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Remove time attribute as we handle it ourselves
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			// Remove level attribute as we handle it ourselves
			if a.Key == slog.LevelKey {
				return slog.Attr{}
			}
			// Remove message attribute as we handle it ourselves
			if a.Key == slog.MessageKey {
				return slog.Attr{}
			}
			// Remove source attribute as we handle it ourselves
			if a.Key == slog.SourceKey {
				return slog.Attr{}
			}
			return a
		},
	}

	handler := NewPrettyHandler(w, opts)
	return slog.New(handler)
}
