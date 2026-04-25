// Package middleware — structured request logging + request ID propagation
//
// Every HTTP request gets a unique ID (X-Request-ID header). The ID is either
// read from an incoming header (useful when a load balancer sets it) or
// generated fresh using Go's crypto/rand UUID. The ID is stored in Fiber's
// context locals so handlers can include it in their own log lines.
//
// The logger middleware emits one structured slog record per request,
// containing: method, path, status, latency, request_id. Cloud Logging
// renders these as indexed fields for easy filtering:
//
//   resource.labels.service_name="runners-list-api"
//   jsonPayload.status=500
package middleware

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const requestIDKey = "request_id"

// RequestID generates (or reads) a unique ID for each request and stores it
// in c.Locals("request_id"). It also echoes it back in the response header
// so clients can correlate requests with server logs.
func RequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Get("X-Request-ID")
		if id == "" {
			id = uuid.New().String()
		}
		c.Locals(requestIDKey, id)
		c.Set("X-Request-ID", id)
		return c.Next()
	}
}

// RequestLogger logs every request as a structured slog record.
// Add this as a global middleware in setupRoutes, before route groups.
//
// Log format (JSON via slog):
//
//	{"time":"...","level":"INFO","msg":"request","method":"GET","path":"/api/v1/events","status":200,"latency_ms":4,"request_id":"..."}
func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		latency := time.Since(start)
		status := c.Response().StatusCode()

		reqID, _ := c.Locals(requestIDKey).(string)

		level := slog.LevelInfo
		if status >= 500 {
			level = slog.LevelError
		} else if status >= 400 {
			level = slog.LevelWarn
		}

		slog.Log(c.Context(), level, "request",
			"method", c.Method(),
			"path", c.Path(),
			"status", status,
			"latency_ms", latency.Milliseconds(),
			"request_id", reqID,
		)

		return err
	}
}
