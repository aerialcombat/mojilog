package mojilog

import (
	"bytes"
	"log/slog"
	"testing"
)

func TestEmojiSpacing(t *testing.T) {
	// Test different emojis and their spacing
	testCases := []struct {
		emoji    string
		expected string
		desc     string
	}{
		{"🚀", " ", "rocket (narrow)"},
		{"⚙️", "  ", "gear (wide)"},
		{"🎰", "  ", "slot machine (wide)"},
		{"📊", "  ", "bar chart (wide)"},
		{"💚", "  ", "green heart (wide)"},
		{"🔴", "  ", "red circle (wide)"},
		{"❌", "  ", "cross mark (wide)"},
		{"⚠️", "  ", "warning (wide)"},
		{"ℹ️", "  ", "information (wide)"},
		{"🔍", "  ", "magnifying glass (wide)"},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			spacing := getEmojiSpacing(tc.emoji)
			if spacing != tc.expected {
				t.Errorf("Emoji %s: expected spacing %q, got %q", tc.emoji, tc.expected, spacing)
			}
		})
	}
}

func TestEmojiHandlerIntegration(t *testing.T) {
	// Test the full emoji handler with different emojis
	var buf bytes.Buffer

	// Create a text handler that writes to buffer
	textHandler := slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	// Wrap with emoji handler
	emojiHandler := NewEmojiHandler(textHandler)
	logger := slog.New(emojiHandler)

	// Test different log levels and messages
	logger.Info("Starting application")
	logger.Warn("Configuration issue")
	logger.Error("Connection failed")
	logger.Debug("Processing data")

	output := buf.String()
	t.Logf("Output:\n%s", output)

	// Verify that emojis are present
	if !containsEmoji(output, "ℹ️") {
		t.Error("Info emoji not found in output")
	}
	if !containsEmoji(output, "⚠️") {
		t.Error("Warning emoji not found in output")
	}
	if !containsEmoji(output, "❌") {
		t.Error("Error emoji not found in output")
	}
	if !containsEmoji(output, "🔍") {
		t.Error("Debug emoji not found in output")
	}
}

func containsEmoji(s, emoji string) bool {
	for _, r := range s {
		if string(r) == emoji {
			return true
		}
	}
	return false
}
