package http

import (
	"strconv"
	"time"

	"github.com/aniqaqill/runners-list/internal/core/domain"
	"github.com/aniqaqill/runners-list/internal/core/service"
	"github.com/aniqaqill/runners-list/internal/port"
	"github.com/gofiber/fiber/v2"
)

type EventHandler struct {
	eventService *service.EventService
}

func NewEventHandler(eventService *service.EventService) *EventHandler {
	return &EventHandler{eventService: eventService}
}

// CreateEvent handles the creation of a new event.
func (h *EventHandler) CreateEvent(c *fiber.Ctx) error {
	var event domain.Events
	if err := c.BodyParser(&event); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "invalid input format",
		})
	}

	if err := h.eventService.CreateEvent(&event); err != nil {
		switch err {
		case service.ErrEventDateInPast:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   true,
				"message": "event date must be in the future",
			})
		case service.ErrEventNameNotUnique:
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error":   true,
				"message": "event name must be unique",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   true,
				"message": "failed to create event",
			})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"message": "event created successfully",
	})
}

// ListEvents returns events with optional filtering and pagination.
//
// Query parameters (all optional):
//
//	state=Selangor          — filter by state (case-insensitive)
//	from=2026-01-01         — events on or after date (YYYY-MM-DD)
//	to=2026-12-31           — events on or before date (YYYY-MM-DD)
//	limit=50                — max results (default 50)
//	offset=0                — skip N results (for paging)
//
// Cache-Control header tells Cloudflare and the browser to cache for 60s,
// matching the Next.js ISR revalidation interval.
func (h *EventHandler) ListEvents(c *fiber.Ctx) error {
	filter := port.EventFilter{
		State:  c.Query("state"),
		Limit:  parseIntQuery(c, "limit", 50),
		Offset: parseIntQuery(c, "offset", 0),
	}

	if from := c.Query("from"); from != "" {
		if t, err := time.Parse("2006-01-02", from); err == nil {
			filter.From = t
		}
	}
	if to := c.Query("to"); to != "" {
		if t, err := time.Parse("2006-01-02", to); err == nil {
			filter.To = t
		}
	}

	events, err := h.eventService.ListEvents(filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "failed to retrieve events",
		})
	}

	// Cache-Control: allow Cloudflare/browser to cache for 60s.
	// s-maxage controls CDN caching; max-age controls browser caching.
	c.Set("Cache-Control", "public, max-age=60, s-maxage=60")

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error":  false,
		"data":   events,
		"total":  len(events),
		"limit":  filter.Limit,
		"offset": filter.Offset,
	})
}

// SyncEvents handles bulk event synchronization from the scraper.
func (h *EventHandler) SyncEvents(c *fiber.Ctx) error {
	var syncReq SyncRequest
	if err := c.BodyParser(&syncReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(SyncResponse{
			Success: false,
			Error:   "invalid request format",
		})
	}

	if len(syncReq.Events) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(SyncResponse{
			Success: false,
			Error:   "no events provided",
		})
	}

	var (
		events     []domain.Events
		rowErrors  []SyncRowError
	)

	for i, eventInput := range syncReq.Events {
		event, err := eventInput.ToEvent()
		if err != nil {
			// Collect per-row errors instead of silently skipping
			rowErrors = append(rowErrors, SyncRowError{
				Index:  i,
				Reason: "invalid date format: " + err.Error(),
			})
			continue
		}
		events = append(events, event)
	}

	inserted, updated, err := h.eventService.BulkUpsertEvents(events)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(SyncResponse{
			Success: false,
			Error:   "failed to sync events: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(SyncResponse{
		Success:   true,
		Inserted:  inserted,
		Updated:   updated,
		Total:     len(syncReq.Events),
		RowErrors: rowErrors,
	})
}

// DeleteEvent handles the deletion of an event by its ID.
func (h *EventHandler) DeleteEvent(c *fiber.Ctx) error {
	id := c.Params("id")
	eventID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "invalid event ID",
		})
	}
	if err := h.eventService.DeleteEvent(uint(eventID)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "failed to delete event",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error":   false,
		"message": "event deleted successfully",
	})
}

func parseIntQuery(c *fiber.Ctx, key string, def int) int {
	v := c.Query(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil || n < 0 {
		return def
	}
	return n
}
