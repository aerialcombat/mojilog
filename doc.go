// Package mojilog provides a beautiful, emoji-enhanced structured logging library for Go.
// It supports multiple output formats including pretty-printed text with colors,
// standard JSON, and pretty-printed JSON with syntax highlighting.
//
// Features:
//   - 🎨 Colorful terminal output
//   - 📝 Multiple format support (text, json, pretty-json)
//   - 🚀 High performance with zero allocations in hot path
//   - 🔧 Compatible with standard slog package
//   - 😊 Emoji indicators for log levels
//   - 🌍 Global logger with thread-safe initialization
//
// Basic usage:
//
//	import "github.com/aerialcombat/mojilog"
//	import "log/slog"
//
//	func main() {
//	    // Initialize global logger
//	    mojilog.InitGlobal(slog.LevelInfo, "pretty", true)
//
//	    // Use it anywhere
//	    log := mojilog.Get()
//	    log.Info("Server started", "port", 8080)
//	}
package mojilog