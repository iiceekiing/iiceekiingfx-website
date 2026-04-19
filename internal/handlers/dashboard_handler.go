package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"iiceekiingfx.com/internal/models"
	"iiceekiingfx.com/internal/services"
)

type DashboardHandler struct {
	dashboardService *services.DashboardService
}

func NewDashboardHandler(dashboardService *services.DashboardService) *DashboardHandler {
	return &DashboardHandler{
		dashboardService: dashboardService,
	}
}

func (h *DashboardHandler) GetOverview(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	overview, err := h.dashboardService.GetOverview(user.ID.String())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch dashboard overview",
		})
	}

	return c.JSON(overview)
}

func (h *DashboardHandler) GetEquityCurve(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	// Get days from query parameter, default to 30
	daysStr := c.Query("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 {
		days = 30
	}

	equityData, err := h.dashboardService.GetEquityCurve(user.ID.String(), days)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch equity curve data",
		})
	}

	return c.JSON(equityData)
}

func (h *DashboardHandler) GetActivity(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	// Get limit from query parameter, default to 10
	limitStr := c.Query("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	activities, err := h.dashboardService.GetActivityFeed(user.ID.String(), limit)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch activity feed",
		})
	}

	return c.JSON(activities)
}
