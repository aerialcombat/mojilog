package main

import (
	"log/slog"

	"github.com/aerialcombat/mojilog"
)

func main() {
	// Initialize with pretty format and source location
	mojilog.InitGlobal(slog.LevelDebug, "pretty", true)

	// Get the global logger
	log := mojilog.Get()

	// Log various levels
	log.Debug("This is a debug message", "component", "main", "iteration", 1)
	log.Info("Application started successfully", "version", "1.0.0", "env", "development")
	log.Warn("Cache miss detected", "key", "user:123", "fallback", "database")
	log.Error("Failed to connect to service", "service", "redis", "error", "timeout after 5s")

	// Use With for contextual logging
	userLog := mojilog.With("user_id", "user-123", "request_id", "req-456")
	userLog.Info("User logged in", "ip", "192.168.1.1")
	userLog.Info("User viewed profile", "profile_id", "profile-789")
}