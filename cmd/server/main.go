package main

import (
	"log"
	"os"

	"iiceekiingfx.com/config"
	"iiceekiingfx.com/internal/handlers"
	"iiceekiingfx.com/internal/middleware"
	"iiceekiingfx.com/internal/repositories"
	"iiceekiingfx.com/internal/services"
	"iiceekiingfx.com/pkg/crypto"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	db, err := repositories.NewDatabase(
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
	)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	portfolioRepo := repositories.NewPortfolioRepository(db)
	journalRepo := repositories.NewJournalRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	dashboardService := services.NewDashboardService(portfolioRepo, journalRepo)
	journalService := services.NewJournalService(journalRepo)

	// Initialize encryption service
	encryptionSvc := crypto.NewEncryptionService(cfg.JWTSecret)
	portfolioService := services.NewPortfolioService(portfolioRepo, encryptionSvc)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	dashboardHandler := handlers.NewDashboardHandler(dashboardService)
	journalHandler := handlers.NewJournalHandler(journalService)
	portfolioHandler := handlers.NewPortfolioHandler(portfolioService)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// Root route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to iiceekiingfx.com API",
			"version": "1.0.0",
			"service": "iiceekiingfx-api",
			"endpoints": fiber.Map{
				"health":    "/health",
				"auth":      "/api/auth",
				"dashboard": "/api/dashboard",
				"portfolio": "/api/portfolio",
				"journal":   "/api/journal",
				"courses":   "/api/courses",
				"signals":   "/api/signals",
			},
		})
	})

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "iiceekiingfx-api",
		})
	})

	// Routes
	api := app.Group("/api")

	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Get("/me", middleware.AuthMiddleware(authService), authHandler.GetMe)

	// Protected routes
	protected := api.Group("/", middleware.AuthMiddleware(authService))

	// Dashboard routes
	dashboard := protected.Group("/dashboard")
	dashboard.Get("/overview", dashboardHandler.GetOverview)
	dashboard.Get("/equity-curve", dashboardHandler.GetEquityCurve)
	dashboard.Get("/activity", dashboardHandler.GetActivity)

	// Portfolio routes
	portfolio := protected.Group("/portfolio")
	portfolio.Post("/connect", portfolioHandler.ConnectAccount)
	portfolio.Get("/accounts", portfolioHandler.GetAccounts)
	portfolio.Get("/accounts/:id", portfolioHandler.GetAccountDetails)
	portfolio.Put("/accounts/:id/status", portfolioHandler.UpdateAccountStatus)
	portfolio.Post("/accounts/:id/sync", portfolioHandler.SyncAccount)
	portfolio.Get("/accounts/:id/performance", portfolioHandler.GetAccountPerformance)
	portfolio.Get("/history", portfolioHandler.GetTradingHistory)

	// Journal routes
	journal := protected.Group("/journal")
	journal.Get("/", journalHandler.GetJournalEntries)
	journal.Post("/", journalHandler.CreateJournalEntry)
	journal.Put("/:id", journalHandler.UpdateJournalEntry)
	journal.Delete("/:id", journalHandler.DeleteJournalEntry)

	// Courses routes
	courses := protected.Group("/courses")
	courses.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Get courses endpoint",
		})
	})
	courses.Get("/:id", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Get course details endpoint",
		})
	})

	// Signals routes
	signals := protected.Group("/signals")
	signals.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Get signals endpoint",
		})
	})

	// Start server
	port := ":" + cfg.ServerPort
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = ":" + envPort
	}

	log.Printf("Server starting on port %s", port)
	if err := app.Listen(port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
