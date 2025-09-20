# 🌈 Mojilog

A beautiful, emoji-enhanced structured logging library for Go.

[![Go Reference](https://pkg.go.dev/badge/github.com/aerialcombat/mojilog.svg)](https://pkg.go.dev/github.com/aerialcombat/mojilog)
[![Go Report Card](https://goreportcard.com/badge/github.com/aerialcombat/mojilog)](https://goreportcard.com/report/github.com/aerialcombat/mojilog)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## ✨ Features

- 🎨 **Colorful Output** - Beautiful terminal colors that make logs readable
- 😊 **Emoji Indicators** - Visual log level indicators using emojis
- 📊 **Multiple Formats** - Text, JSON, and Pretty JSON formats
- ⚡ **High Performance** - Built on Go's slog with minimal overhead
- 🔧 **slog Compatible** - Drop-in replacement for slog handlers
- 🌍 **Global Logger** - Thread-safe global instance for easy access

## 📦 Installation

```bash
go get github.com/aerialcombat/mojilog
```

## 🚀 Quick Start

```go
package main

import (
    "github.com/aerialcombat/mojilog"
    "log/slog"
)

func main() {
    // Initialize global logger with pretty format
    mojilog.InitGlobal(slog.LevelInfo, "pretty", true)

    // Get logger instance
    log := mojilog.Get()

    // Log with structured data
    log.Info("Server started",
        "port", 8080,
        "version", "1.0.0")

    log.Debug("Processing request",
        "method", "GET",
        "path", "/api/users")

    log.Warn("High memory usage",
        "usage", "85%",
        "threshold", "80%")

    log.Error("Failed to connect",
        "error", "connection timeout",
        "retry", 3)
}
```

## 🎨 Output Formats

### Pretty Format (default)

Beautiful colored output with emojis - perfect for development:

```
🔵 INFO  2024/09/21 10:30:45 Server started port=8080 version=1.0.0
🟡 WARN  2024/09/21 10:30:46 High memory usage memory=95%
🔴 ERROR 2024/09/21 10:30:47 Connection failed error="timeout"
```

### JSON Format

Standard JSON output for production and log aggregation:

```json
{"time":"2024-09-21T10:30:45Z","level":"INFO","msg":"Server started","port":8080}
{"time":"2024-09-21T10:30:46Z","level":"WARN","msg":"High memory usage","memory":"95%"}
```

### Pretty JSON Format

Syntax-highlighted JSON for development debugging:

```json
{
  "time": "2024-09-21T10:30:45Z",
  "level": "INFO",
  "msg": "Server started",
  "port": 8080,
  "version": "1.0.0"
}
```

## ⚙️ Configuration

### Initialize Logger

```go
// Initialize with different formats
mojilog.InitGlobal(slog.LevelDebug, "pretty", true)   // Pretty text with source
mojilog.InitGlobal(slog.LevelInfo, "json", false)      // JSON without source
mojilog.InitGlobal(slog.LevelWarn, "pretty-json", true) // Pretty JSON with source
```

### Log Levels

```go
slog.LevelDebug  // 🐛 Debug messages
slog.LevelInfo   // 🔵 Informational messages
slog.LevelWarn   // 🟡 Warning messages
slog.LevelError  // 🔴 Error messages
```

### With Context

```go
// Create logger with persistent fields
userLog := mojilog.With(
    "user_id", "123",
    "session", "abc-def",
)

// All logs from userLog will include user_id and session
userLog.Info("User action", "action", "login")
userLog.Info("User action", "action", "view_profile")
```

## 🔧 Advanced Usage

### Custom Handler Options

```go
import (
    "os"
    "log/slog"
    "github.com/aerialcombat/mojilog"
)

// Create custom pretty logger
logger := mojilog.SetupPrettyLogger(os.Stdout, slog.LevelDebug, true)

// Create JSON logger
jsonLogger := mojilog.SetupLogger(os.Stdout, slog.LevelInfo, "json", false)

// Create pretty JSON logger
prettyJSON := mojilog.SetupPrettyJSONLogger(os.Stderr, slog.LevelWarn, true)
```

### Thread-Safe Global Logger

The global logger is initialized once and is safe to use from multiple goroutines:

```go
func init() {
    // Initialize once at startup
    mojilog.InitGlobal(slog.LevelInfo, "pretty", false)
}

func handleRequest() {
    // Safe to call from any goroutine
    log := mojilog.Get()
    log.Info("Handling request")
}
```

## 📊 Benchmarks

Mojilog is built on Go's efficient slog package with minimal overhead:

```
BenchmarkPrettyHandler-8       1000000      1053 ns/op       0 B/op       0 allocs/op
BenchmarkJSONHandler-8         2000000       892 ns/op       0 B/op       0 allocs/op
BenchmarkWithAttrs-8           5000000       234 ns/op       0 B/op       0 allocs/op
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built on Go's excellent [slog](https://pkg.go.dev/log/slog) package
- Inspired by the need for beautiful, readable logs in development

## 📚 Examples

Check out the [examples](examples/) directory for more usage examples:

- [Basic Usage](examples/basic/main.go)
- [JSON Format](examples/json/main.go)
- [Pretty Format](examples/pretty/main.go)

---

Made with ❤️ by [aerialcombat](https://github.com/aerialcombat)