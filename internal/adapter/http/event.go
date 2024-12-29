package http

import (
	"strconv"

	"github.com/aniqaqill/runners-list/internal/core/domain"
	"github.com/aniqaqill/runners-list/internal/core/service"
	"github.com/gofiber/fiber/v2"
)

type EventHandler struct {
	eventService *service.EventService
}

func NewEventHandler(eventService *service.EventService) *EventHandler {
	return &EventHandler{eventService: eventService}
}

func (h *EventHandler) CreateEvent(c *fiber.Ctx) error {
	var event domain.RunningEvents
	if err := c.BodyParser(&event); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid input format",
		})
	}

	if err := h.eventService.CreateEvent(&event); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to create event",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"message": "Event created successfully",
	})
}

func (h *EventHandler) ListEvents(c *fiber.Ctx) error {
	events, err := h.eventService.ListEvents()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to retrieve events",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"data":  events,
	})
}

func (h *EventHandler) DeleteEvent(c *fiber.Ctx) error {
	id := c.Params("id")
	eventID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid event ID",
		})
	}
	if err := h.eventService.DeleteEvent(uint(eventID)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to delete event",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error":   false,
		"message": "Event deleted successfully",
	})
}
