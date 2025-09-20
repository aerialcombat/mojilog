package mojilog

import (
	"context"
	"io"
	"log/slog"

	"github.com/mattn/go-runewidth"
)

// EmojiHandler wraps another handler and adds emojis based on log level
type EmojiHandler struct {
	wrapped slog.Handler
}

// NewEmojiHandler creates a new emoji handler that wraps the given handler
func NewEmojiHandler(wrapped slog.Handler) *EmojiHandler {
	return &EmojiHandler{wrapped: wrapped}
}

// Enabled implements slog.Handler
func (h *EmojiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.wrapped.Enabled(ctx, level)
}

// Handle implements slog.Handler
func (h *EmojiHandler) Handle(ctx context.Context, r slog.Record) error {
	// Get contextual emoji first (higher priority)
	contextEmoji := getContextualEmoji(r.Message)

	// Use contextual emoji if available, otherwise use level emoji
	emoji := ""
	if contextEmoji != "" {
		emoji = contextEmoji
	} else {
		emoji = getEmojiForLevel(r.Level)
	}

	// Prepend emoji to the message with appropriate spacing
	// if the emoji takes up double-space like ⚙️, add a double space,and if emoji takes up single-space like 🚀, then add a single space
	if emoji != "" {
		spacing := getEmojiSpacing(emoji)
		r.Message = emoji + spacing + r.Message
	}

	return h.wrapped.Handle(ctx, r)
}

// WithAttrs implements slog.Handler
func (h *EmojiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &EmojiHandler{wrapped: h.wrapped.WithAttrs(attrs)}
}

// WithGroup implements slog.Handler
func (h *EmojiHandler) WithGroup(name string) slog.Handler {
	return &EmojiHandler{wrapped: h.wrapped.WithGroup(name)}
}

// getEmojiForLevel returns an emoji based on the log level
func getEmojiForLevel(level slog.Level) string {
	switch {
	case level >= slog.LevelError:
		return "❌"
	case level >= slog.LevelWarn:
		return "⚠️"
	case level >= slog.LevelInfo:
		return "ℹ️"
	case level >= slog.LevelDebug:
		return "🔍"
	default:
		return "📝"
	}
}

// getEmojiSpacing returns appropriate spacing based on emoji display width
func getEmojiSpacing(emoji string) string {
	// Use runewidth library to determine actual display width
	width := runewidth.StringWidth(emoji)

	// Most emojis are either 1 or 2 characters wide in terminal
	if width > 1 {
		return " " // double space for wide emojis
	}

	return "  " // single space for narrow emojis
}

// getContextualEmoji returns a context-specific emoji based on message content
func getContextualEmoji(msg string) string {
	lowerMsg := toLower(msg)

	// Health status patterns - highest priority
	if stringContains(lowerMsg, "health") {
		if stringContains(lowerMsg, "excellent") {
			return "💚"
		} else if stringContains(lowerMsg, "good") {
			return "🟡"
		} else if stringContains(lowerMsg, "degraded") {
			return "🟠"
		} else if stringContains(lowerMsg, "critical") {
			return "🔴"
		}
	}

	// System states - check for more specific patterns first
	if stringContains(lowerMsg, "shutdown") || stringContains(lowerMsg, "stopping") {
		return "🛑"
	}
	// Only use rocket for main application start, not sub-components
	if stringContains(lowerMsg, "start") || stringContains(lowerMsg, "parser is running") {
		return "🚀"
	}

	if stringContains(lowerMsg, "success") {
		return "🎉"
	}

	if stringContains(lowerMsg, "cleanup") {
		return "🧹"
	}

	// Operations
	if stringContains(lowerMsg, "config") || stringContains(lowerMsg, "setting") {
		return "⚙️"
	}
	if stringContains(lowerMsg, "connect") || stringContains(lowerMsg, "websocket") {
		return "🔌"
	}
	if stringContains(lowerMsg, "failed") {
		return "❌"
	}
	if stringContains(lowerMsg, "table") || stringContains(lowerMsg, "game") || stringContains(lowerMsg, "casino") {
		return "🎰"
	}
	if stringContains(lowerMsg, "statistics") || stringContains(lowerMsg, "metrics") {
		return "📊"
	}
	if stringContains(lowerMsg, "loading") || stringContains(lowerMsg, "processing") {
		return "⏳"
	}

	if stringContains(lowerMsg, "creating") {
		return "🆕"
	}

	return ""
}

// Simple lowercase conversion
func toLower(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if 'A' <= c && c <= 'Z' {
			result[i] = c + 32
		} else {
			result[i] = c
		}
	}
	return string(result)
}

// Simple contains check
func stringContains(s, substr string) bool {
	if len(substr) > len(s) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// SetupLogger sets up a global logger with emoji support
func SetupLogger(w io.Writer, level slog.Level, format string, addSource bool) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: addSource,
	}

	var baseHandler slog.Handler
	if format == "json" {
		baseHandler = slog.NewJSONHandler(w, opts)
	} else {
		baseHandler = slog.NewTextHandler(w, opts)
	}

	// Wrap with emoji handler
	emojiHandler := NewEmojiHandler(baseHandler)

	return slog.New(emojiHandler)
}
