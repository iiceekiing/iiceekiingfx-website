package main

import (
	"log"
	"os"

	"iiceekiingfx.com/config"
	"iiceekiingfx.com/internal/handlers"
	"iiceekiingfx.com/internal/middleware"
	"iiceekiingfx.com/internal/repositories"
	"iiceekiingfx.com/internal/services"

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

	// Initialize services
	authService := services.NewAuthService(userRepo, cfg.JWTSecret)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)

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
	dashboard.Get("/overview", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Dashboard overview endpoint",
		})
	})
	dashboard.Get("/equity-curve", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Equity curve endpoint",
		})
	})
	dashboard.Get("/activity", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Activity feed endpoint",
		})
	})

	// Portfolio routes
	portfolio := protected.Group("/portfolio")
	portfolio.Post("/connect", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Connect MT5 account endpoint",
		})
	})
	portfolio.Get("/accounts", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Get trading accounts endpoint",
		})
	})
	portfolio.Get("/history", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Get trading history endpoint",
		})
	})

	// Journal routes
	journal := protected.Group("/journal")
	journal.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Get journal entries endpoint",
		})
	})
	journal.Post("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Create journal entry endpoint",
		})
	})
	journal.Put("/:id", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Update journal entry endpoint",
		})
	})
	journal.Delete("/:id", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Delete journal entry endpoint",
		})
	})

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
