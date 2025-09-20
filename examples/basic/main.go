package main

import (
	"github.com/aerialcombat/mojilog"
)

func main() {
	// Initialize with pretty format and source location

	// Log various levels
	mojilog.Debug("This is a debug message", "component", "main", "iteration", 1)
	mojilog.Info("Application started successfully", "version", "1.0.0", "env", "development")
	mojilog.Warn("Cache miss detected", "key", "user:123", "fallback", "database")
	mojilog.Error("Failed to connect to service", "service", "redis", "error", "timeout after 5s")

	// Use With for contextual logging
	userLog := mojilog.With("user_id", "user-123", "request_id", "req-456")
	userLog.Info("User logged in", "ip", "192.168.1.1")
	userLog.Info("User viewed profile", "profile_id", "profile-789")
}
