package handlers

import (
	"strconv"

	"iiceekiingfx.com/internal/models"
	"iiceekiingfx.com/internal/services"
	"github.com/gofiber/fiber/v2"
)

type SignalsHandler struct {
	signalsService *services.SignalsService
}

func NewSignalsHandler(signalsService *services.SignalsService) *SignalsHandler {
	return &SignalsHandler{
		signalsService: signalsService,
	}
}

func (h *SignalsHandler) GetAllSignals(c *fiber.Ctx) error {
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

	signals, total, err := h.signalsService.GetAllSignals(page, limit)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch signals",
		})
	}

	return c.JSON(fiber.Map{
		"signals": signals,
		"page":    page,
		"limit":   limit,
		"total":   total,
	})
}

func (h *SignalsHandler) GetSignalByID(c *fiber.Ctx) error {
	signalID := c.Params("id")

	signal, err := h.signalsService.GetSignalByID(signalID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(signal)
}

func (h *SignalsHandler) GetSignalsByPair(c *fiber.Ctx) error {
	pair := c.Params("pair")

	signals, err := h.signalsService.GetSignalsByPair(pair)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch signals",
		})
	}

	return c.JSON(fiber.Map{
		"signals": signals,
		"pair":    pair,
		"count":   len(signals),
	})
}

// Admin only endpoints
func (h *SignalsHandler) CreateSignal(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	
	// Check if user is admin
	if user.Role != models.RoleAdmin {
		return c.Status(403).JSON(fiber.Map{
			"error": "Admin access required",
		})
	}

	var req models.SignalRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	signal, err := h.signalsService.CreateSignal(user.ID.String(), &req)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Signal created successfully",
		"signal":  signal,
	})
}

func (h *SignalsHandler) UpdateSignal(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	
	// Check if user is admin
	if user.Role != models.RoleAdmin {
		return c.Status(403).JSON(fiber.Map{
			"error": "Admin access required",
		})
	}

	signalID := c.Params("id")
	var req models.SignalRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	signal, err := h.signalsService.UpdateSignal(user.ID.String(), signalID, &req)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Signal updated successfully",
		"signal":  signal,
	})
}

func (h *SignalsHandler) DeleteSignal(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	
	// Check if user is admin
	if user.Role != models.RoleAdmin {
		return c.Status(403).JSON(fiber.Map{
			"error": "Admin access required",
		})
	}

	signalID := c.Params("id")

	err := h.signalsService.DeleteSignal(user.ID.String(), signalID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Signal deleted successfully",
	})
}

func (h *SignalsHandler) CloseSignal(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	
	// Check if user is admin
	if user.Role != models.RoleAdmin {
		return c.Status(403).JSON(fiber.Map{
			"error": "Admin access required",
		})
	}

	signalID := c.Params("id")

	err := h.signalsService.CloseSignal(signalID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Signal closed successfully",
	})
}
