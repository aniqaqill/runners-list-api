package middleware

import (
	"crypto/subtle"

	"github.com/gofiber/fiber/v2"
)

// InternalAPIKeyAuth returns a Fiber middleware that validates the
// X-Internal-Token header against the expected key.
//
// The expected key is injected at construction time (from Config) instead of
// being read from os.Getenv on every request. This makes the middleware
// trivially testable and removes a hidden global dependency.
//
// subtle.ConstantTimeCompare is used to prevent timing-based side-channel
// attacks: the comparison always takes the same amount of time regardless of
// where the strings differ.
func InternalAPIKeyAuth(expectedKey string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if expectedKey == "" {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "internal API key not configured on server",
			})
		}

		apiKey := c.Get("X-Internal-Token")
		if apiKey == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "missing X-Internal-Token header",
			})
		}

		// constant-time byte comparison — same duration whether keys match or not
		if subtle.ConstantTimeCompare([]byte(apiKey), []byte(expectedKey)) != 1 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "invalid API key",
			})
		}

		return c.Next()
	}
}
