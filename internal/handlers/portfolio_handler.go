package handlers

import (
	"github.com/gofiber/fiber/v2"
	"iiceekiingfx.com/internal/models"
	"iiceekiingfx.com/internal/services"
)

type PortfolioHandler struct {
	portfolioService *services.PortfolioService
}

func NewPortfolioHandler(portfolioService *services.PortfolioService) *PortfolioHandler {
	return &PortfolioHandler{
		portfolioService: portfolioService,
	}
}

func (h *PortfolioHandler) ConnectAccount(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	var req models.ConnectAccountRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	account, err := h.portfolioService.ConnectAccount(user.ID.String(), &req)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Account connected successfully",
		"account": account,
	})
}

func (h *PortfolioHandler) GetAccounts(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	accounts, err := h.portfolioService.GetUserAccounts(user.ID.String())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch accounts",
		})
	}

	return c.JSON(fiber.Map{
		"accounts": accounts,
		"count":    len(accounts),
	})
}

func (h *PortfolioHandler) GetAccountDetails(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	accountID := c.Params("id")

	account, err := h.portfolioService.GetAccountByID(user.ID.String(), accountID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Account not found",
		})
	}

	return c.JSON(account)
}

func (h *PortfolioHandler) UpdateAccountStatus(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	accountID := c.Params("id")

	var req struct {
		IsActive bool `json:"is_active"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	err := h.portfolioService.UpdateAccountStatus(user.ID.String(), accountID, req.IsActive)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Account status updated successfully",
	})
}

func (h *PortfolioHandler) SyncAccount(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	accountID := c.Params("id")

	err := h.portfolioService.SyncAccount(user.ID.String(), accountID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Account synchronization started",
	})
}

func (h *PortfolioHandler) GetAccountPerformance(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	accountID := c.Params("id")

	performance, err := h.portfolioService.GetAccountPerformance(user.ID.String(), accountID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(performance)
}

func (h *PortfolioHandler) GetTradingHistory(c *fiber.Ctx) error {
	// This would get trading history for all user accounts
	// For now, return placeholder
	return c.JSON(fiber.Map{
		"message": "Trading history endpoint",
		"trades":  []interface{}{},
		"total":   0,
	})
}
