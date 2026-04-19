package services

import (
	"errors"
	"time"

	"iiceekiingfx.com/internal/models"
	"iiceekiingfx.com/internal/repositories"
	"github.com/google/uuid"
)

type CourseService struct {
	courseRepo *repositories.CourseRepository
}

func NewCourseService(courseRepo *repositories.CourseRepository) *CourseService {
	return &CourseService{
		courseRepo: courseRepo,
	}
}

func (s *CourseService) GetAllCourses() ([]*models.Course, error) {
	return s.courseRepo.GetAllCourses()
}

func (s *CourseService) GetCourseByID(courseID string) (*models.Course, error) {
	if courseID == "" {
		return nil, errors.New("course ID is required")
	}

	course, err := s.courseRepo.GetCourseByID(courseID)
	if err != nil {
		return nil, err
	}
	if course == nil {
		return nil, errors.New("course not found")
	}

	return course, nil
}

func (s *CourseService) GetCourseLessons(courseID string) ([]*models.Lesson, error) {
	if courseID == "" {
		return nil, errors.New("course ID is required")
	}

	lessons, err := s.courseRepo.GetLessonsByCourseID(courseID)
	if err != nil {
		return nil, err
	}

	return lessons, nil
}

func (s *CourseService) GetLessonByID(lessonID string) (*models.Lesson, error) {
	if lessonID == "" {
		return nil, errors.New("lesson ID is required")
	}

	lesson, err := s.courseRepo.GetLessonByID(lessonID)
	if err != nil {
		return nil, err
	}
	if lesson == nil {
		return nil, errors.New("lesson not found")
	}

	return lesson, nil
}

func (s *CourseService) EnrollUser(userID, courseID string) (*models.CourseProgress, error) {
	if userID == "" || courseID == "" {
		return nil, errors.New("user ID and course ID are required")
	}

	// Check if course exists
	course, err := s.courseRepo.GetCourseByID(courseID)
	if err != nil {
		return nil, err
	}
	if course == nil {
		return nil, errors.New("course not found")
	}

	// Check if already enrolled
	existingProgress, err := s.courseRepo.GetUserProgress(userID, courseID)
	if err != nil {
		return nil, err
	}
	if existingProgress != nil {
		return existingProgress, nil
	}

	// Create new enrollment
	progress := &models.CourseProgress{
		ID:        uuid.New(),
		UserID:    uuid.MustParse(userID),
		CourseID:  uuid.MustParse(courseID),
		Progress:  0.0,
		Completed: false,
		StartedAt: time.Now(),
		UpdatedAt:  time.Now(),
	}

	err = s.courseRepo.CreateOrUpdateProgress(progress)
	if err != nil {
		return nil, err
	}

	return progress, nil
}

func (s *CourseService) GetUserProgress(userID, courseID string) (*models.CourseProgress, error) {
	if userID == "" || courseID == "" {
		return nil, errors.New("user ID and course ID are required")
	}

	progress, err := s.courseRepo.GetUserProgress(userID, courseID)
	if err != nil {
		return nil, err
	}

	return progress, nil
}

func (s *CourseService) UpdateProgress(userID, courseID string, progressPercentage float64) error {
	if userID == "" || courseID == "" {
		return errors.New("user ID and course ID are required")
	}

	if progressPercentage < 0 || progressPercentage > 100 {
		return errors.New("progress must be between 0 and 100")
	}

	// Check if enrolled
	existingProgress, err := s.courseRepo.GetUserProgress(userID, courseID)
	if err != nil {
		return err
	}
	if existingProgress == nil {
		return errors.New("user not enrolled in this course")
	}

	// Update progress
	updatedProgress := &models.CourseProgress{
		UserID:    uuid.MustParse(userID),
		CourseID:  uuid.MustParse(courseID),
		Progress:  progressPercentage,
		Completed: progressPercentage >= 100.0,
		StartedAt: existingProgress.StartedAt,
		UpdatedAt:  time.Now(),
	}

	err = s.courseRepo.CreateOrUpdateProgress(updatedProgress)
	if err != nil {
		return err
	}

	return nil
}

func (s *CourseService) GetUserCourses(userID string) ([]*models.CourseProgress, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	return s.courseRepo.GetUserCourses(userID)
}

func (s *CourseService) CreateCourse(course *models.Course) error {
	if course.Title == "" {
		return errors.New("course title is required")
	}

	// Set default values
	if course.IsActive == false {
		course.IsActive = true
	}

	course.ID = uuid.New()
	course.CreatedAt = time.Now()
	course.UpdatedAt = time.Now()

	return s.courseRepo.CreateCourse(course)
}

func (s *CourseService) CreateLesson(lesson *models.Lesson) error {
	if lesson.Title == "" || lesson.CourseID.String() == "" {
		return errors.New("lesson title and course ID are required")
	}

	lesson.ID = uuid.New()
	lesson.CreatedAt = time.Now()
	lesson.UpdatedAt = time.Now()

	return s.courseRepo.CreateLesson(lesson)
}
