// Package middleware — rate limiting
//
// Applied to the public GET /events endpoint to protect against runaway
// scrapers or accidental DDoS.
//
// Fiber's built-in limiter uses an in-memory sliding window. This is suitable
// for single-instance Cloud Run (one process). For multi-instance deployments,
// move to a Redis-backed limiter (e.g. gofiber/storage/redis).
//
// Cloudflare's free tier also handles rate limiting at the edge — this is a
// defence-in-depth layer at the application level.
package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// RateLimit returns a Fiber middleware that caps requests per IP.
//
// Default: 60 requests per minute per IP.
// Exceeding the limit → 429 Too Many Requests with a Retry-After header.
func RateLimit() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        60,             // requests per window
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			// Use X-Forwarded-For if present (Cloud Run / Cloudflare set this)
			if ip := c.Get("X-Forwarded-For"); ip != "" {
				return ip
			}
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":   true,
				"message": "rate limit exceeded, please retry after 60 seconds",
			})
		},
	})
}
