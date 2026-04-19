package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"iiceekiingfx.com/internal/models"
	"iiceekiingfx.com/internal/services"
)

type CourseHandler struct {
	courseService *services.CourseService
}

func NewCourseHandler(courseService *services.CourseService) *CourseHandler {
	return &CourseHandler{
		courseService: courseService,
	}
}

func (h *CourseHandler) GetAllCourses(c *fiber.Ctx) error {
	courses, err := h.courseService.GetAllCourses()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch courses",
		})
	}

	return c.JSON(fiber.Map{
		"courses": courses,
		"count":   len(courses),
	})
}

func (h *CourseHandler) GetCourseByID(c *fiber.Ctx) error {
	courseID := c.Params("id")

	course, err := h.courseService.GetCourseByID(courseID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(course)
}

func (h *CourseHandler) GetCourseLessons(c *fiber.Ctx) error {
	courseID := c.Params("id")

	lessons, err := h.courseService.GetCourseLessons(courseID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"lessons": lessons,
		"count":   len(lessons),
	})
}

func (h *CourseHandler) GetLessonByID(c *fiber.Ctx) error {
	lessonID := c.Params("lessonId")

	lesson, err := h.courseService.GetLessonByID(lessonID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(lesson)
}

func (h *CourseHandler) EnrollCourse(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	courseID := c.Params("id")

	progress, err := h.courseService.EnrollUser(user.ID.String(), courseID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message":  "Successfully enrolled in course",
		"progress": progress,
	})
}

func (h *CourseHandler) GetUserProgress(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	courseID := c.Params("id")

	progress, err := h.courseService.GetUserProgress(user.ID.String(), courseID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if progress == nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "No progress found for this course",
		})
	}

	return c.JSON(progress)
}

func (h *CourseHandler) UpdateProgress(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	courseID := c.Params("id")

	var req struct {
		Progress float64 `json:"progress"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	err := h.courseService.UpdateProgress(user.ID.String(), courseID, req.Progress)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Progress updated successfully",
	})
}

func (h *CourseHandler) GetUserCourses(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	// Get pagination parameters
	pageStr := c.Query("page", "1")
	limitStr := c.Query("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	courses, err := h.courseService.GetUserCourses(user.ID.String())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch user courses",
		})
	}

	return c.JSON(fiber.Map{
		"courses": courses,
		"page":    page,
		"limit":   limit,
		"total":   len(courses),
	})
}

// Admin only endpoints
func (h *CourseHandler) CreateCourse(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	// Check if user is admin
	if user.Role != models.RoleAdmin {
		return c.Status(403).JSON(fiber.Map{
			"error": "Admin access required",
		})
	}

	var req struct {
		Title       string  `json:"title"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		ImageURL    string  `json:"image_url"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	course := &models.Course{
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		ImageURL:    req.ImageURL,
	}

	err := h.courseService.CreateCourse(course)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Course created successfully",
		"course":  course,
	})
}

func (h *CourseHandler) CreateLesson(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	// Check if user is admin
	if user.Role != models.RoleAdmin {
		return c.Status(403).JSON(fiber.Map{
			"error": "Admin access required",
		})
	}

	courseID := c.Params("id")

	var req struct {
		Title    string `json:"title"`
		Content  string `json:"content"`
		VideoURL string `json:"video_url"`
		Order    int    `json:"order"`
		Duration int    `json:"duration"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	lesson := &models.Lesson{
		CourseID: uuid.MustParse(courseID),
		Title:    req.Title,
		Content:  req.Content,
		VideoURL: req.VideoURL,
		Order:    req.Order,
		Duration: req.Duration,
		IsLocked: true, // Default to locked
	}

	err := h.courseService.CreateLesson(lesson)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Lesson created successfully",
		"lesson":  lesson,
	})
}
