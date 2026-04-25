package main

import (
	"github.com/aniqaqill/runners-list/internal/adapter/http"
	"github.com/aniqaqill/runners-list/internal/adapter/middleware"
	"github.com/aniqaqill/runners-list/internal/config"
	"github.com/aniqaqill/runners-list/internal/platform/observability"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// setupRoutes wires all HTTP routes onto the Fiber app.
// Secrets are passed explicitly from cfg instead of reading os.Getenv inside
// middleware — this makes each handler independently testable.
func setupRoutes(app *fiber.App, db *gorm.DB, cfg *config.Config, eventHandler *http.EventHandler, userHandler *http.UserHandler) {
	// Probes — outside /api/v1 so they are always reachable
	app.Get("/health", observability.HealthHandler())
	app.Get("/ready", observability.ReadyHandler(db))

	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("API v1 is working!")
	})

	v1.Get("/events", middleware.RateLimit(), eventHandler.ListEvents)
	v1.Get("/users", userHandler.ListUsers)

	v1.Post("/register", userHandler.Register)
	v1.Post("/login", userHandler.Login)

	// Internal routes — scraper only.
	// Two layers of auth (both must pass):
	//   1. InternalAPIKeyAuth  — backward-compatible plain key check
	//   2. HMACAuth            — signature + timestamp replay protection
	// Once all scrapers are updated to send X-Timestamp + X-Signature,
	// InternalAPIKeyAuth can be dropped from this group.
	internal := v1.Group("/internal",
		middleware.InternalAPIKeyAuth(cfg.InternalAPIKey),
		middleware.HMACAuth(cfg.InternalAPIKey),
	)
	internal.Post("/sync", eventHandler.SyncEvents)

	// Protected routes — JWT required
	protected := v1.Group("/protected", middleware.JWTProtected(cfg.JWTSecret))
	protected.Post("/events/create-events", middleware.ValidateCreateEventInput, eventHandler.CreateEvent)
	protected.Delete("/events/:id", eventHandler.DeleteEvent)
}
