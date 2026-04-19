package handlers

import (
	"strconv"

	"iiceekiingfx.com/internal/models"
	"iiceekiingfx.com/internal/services"

	"github.com/gofiber/fiber/v2"
)

type JournalHandler struct {
	journalService *services.JournalService
}

func NewJournalHandler(journalService *services.JournalService) *JournalHandler {
	return &JournalHandler{
		journalService: journalService,
	}
}

func (h *JournalHandler) GetJournalEntries(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	
	// Get pagination parameters
	pageStr := c.Query("page", "1")
	limitStr := c.Query("limit", "20")
	
	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	
	offset := (page - 1) * limit

	entries, err := h.journalService.GetUserEntries(user.ID.String(), limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch journal entries",
		})
	}

	return c.JSON(fiber.Map{
		"entries": entries,
		"page": page,
		"limit": limit,
	})
}

func (h *JournalHandler) CreateJournalEntry(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	
	var req models.JournalRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	entry, err := h.journalService.CreateEntry(user.ID.String(), &req)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(201).JSON(entry)
}

func (h *JournalHandler) UpdateJournalEntry(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	
	id := c.Params("id")
	var req models.JournalRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	entry, err := h.journalService.UpdateEntry(user.ID.String(), id, &req)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(entry)
}

func (h *JournalHandler) DeleteJournalEntry(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	
	id := c.Params("id")
	
	err := h.journalService.DeleteEntry(user.ID.String(), id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Journal entry deleted successfully",
	})
}
