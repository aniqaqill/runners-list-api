// Package observability provides health and readiness HTTP handlers.
//
// GET /health  — liveness probe (always 200 when the process is up).
// GET /ready   — readiness probe (pings the DB; fails if DB is unreachable).
//
// Cloud Run uses these probes to decide whether to route traffic to the
// instance. ECS can also consume them via a target-group health check.
package observability

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// HealthHandler returns a Fiber handler that always responds 200 OK.
// Use this for the liveness probe — it only answers "is the process alive?"
func HealthHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "ok",
		})
	}
}

// ReadyHandler returns a Fiber handler that pings the database.
// Use this for the readiness probe — it answers "is the service ready to
// accept traffic?". A 1-second timeout prevents the probe from hanging.
func ReadyHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		sqlDB, err := db.DB()
		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status": "unavailable",
				"reason": "could not obtain db pool",
			})
		}
		if err := sqlDB.PingContext(ctx); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status": "unavailable",
				"reason": "db ping failed",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "ready",
		})
	}
}
