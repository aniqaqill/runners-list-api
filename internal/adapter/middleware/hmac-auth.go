// Package middleware — HMAC request authentication
//
// Why HMAC over a plain API key?
//
//   A plain API key in a header (X-Internal-Token) proves the caller knows the
//   secret, but if the key ever leaks (logs, error messages, replay attacks),
//   an attacker can replay any request indefinitely.
//
//   HMAC-SHA256 over "timestamp:body" adds two properties:
//   1. Request binding   — the signature covers the request body, so a signed
//      request cannot be modified in transit and reused.
//   2. Replay protection — the timestamp claim is validated within a ±5-minute
//      window, so captured requests cannot be replayed hours later.
//
// How it works (scraper → API):
//   1. Scraper reads the current UTC timestamp (Unix seconds) as a string.
//   2. Scraper computes HMAC-SHA256(key, "<timestamp>:<body>") and hex-encodes it.
//   3. Scraper sends headers: X-Timestamp: <ts>  X-Signature: <hex>
//   4. API validates timestamp window, reconstructs the signed string, compares.
//
// Python scraper example:
//   import hmac, hashlib, time, json
//   body = json.dumps(payload).encode()
//   ts = str(int(time.time()))
//   sig = hmac.new(key.encode(), f"{ts}:{body.decode()}".encode(), hashlib.sha256).hexdigest()
//   headers = {"X-Timestamp": ts, "X-Signature": sig}
package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

const hmacTimestampWindow = 5 * time.Minute

// HMACAuth returns a Fiber middleware that validates an HMAC-SHA256 signature.
// The key is the shared secret (same value as InternalAPIKey).
//
// Expected headers:
//
//	X-Timestamp: <unix seconds as string>
//	X-Signature: <hex(HMAC-SHA256(key, "<timestamp>:<raw body>"))>
func HMACAuth(key string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tsHeader := c.Get("X-Timestamp")
		sigHeader := c.Get("X-Signature")

		if tsHeader == "" || sigHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "missing X-Timestamp or X-Signature header",
			})
		}

		tsUnix, err := strconv.ParseInt(tsHeader, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "invalid X-Timestamp format",
			})
		}

		// Replay protection: reject requests older/newer than 5 minutes
		age := time.Since(time.Unix(tsUnix, 0))
		if math.Abs(age.Seconds()) > hmacTimestampWindow.Seconds() {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   fmt.Sprintf("timestamp out of allowed window (±%s)", hmacTimestampWindow),
			})
		}

		// Reconstruct the signed payload
		body := c.Body() // Fiber reads body once; safe to call here
		message := fmt.Sprintf("%s:%s", tsHeader, string(body))

		mac := hmac.New(sha256.New, []byte(key))
		mac.Write([]byte(message))
		expected := hex.EncodeToString(mac.Sum(nil))

		// Constant-time compare — prevents timing attacks on the hex string
		if subtle.ConstantTimeCompare([]byte(sigHeader), []byte(expected)) != 1 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "invalid HMAC signature",
			})
		}

		return c.Next()
	}
}
