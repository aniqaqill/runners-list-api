package handlers

import (
	"github.com/aniqaqill/runners-list/database"
	"github.com/aniqaqill/runners-list/models"

	"github.com/gofiber/fiber/v2"
)

// func for retrieving running event list
func ListEvents(c *fiber.Ctx) error {
	var events []models.RunningEvents
	result := database.DB.Db.Find(&events)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Could not retrieve events",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"data":  events,
	})
}

// func for event creation
func CreateEvents(c *fiber.Ctx) error {
	RunningEvents := new(models.RunningEvents)
	if err := c.BodyParser(RunningEvents); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Could not parse JSON",
		})
	}
	database.DB.Db.Create(&RunningEvents)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"message": "Event created successfully",
	})
}

// func for event deletion
func DeleteEvents(c *fiber.Ctx) error {
	id := c.Params("id")
	var event models.RunningEvents

	// Find the event by ID
	if err := database.DB.Db.First(&event, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Event not found",
		})
	}

	// Delete the event
	if err := database.DB.Db.Delete(&event).Error; err != nil {
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
