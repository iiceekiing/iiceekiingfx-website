package repositories

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"iiceekiingfx.com/internal/models"
)

type CourseRepository struct {
	db *Database
}

func NewCourseRepository(db *Database) *CourseRepository {
	return &CourseRepository{db: db}
}

func (r *CourseRepository) CreateCourse(course *models.Course) error {
	query := `
		INSERT INTO courses (title, description, price, image_url, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.DB.Exec(query,
		course.Title, course.Description, course.Price,
		course.ImageURL, course.IsActive, time.Now(), time.Now(),
	)

	return err
}

func (r *CourseRepository) GetAllCourses() ([]*models.Course, error) {
	query := `
		SELECT id, title, description, price, image_url, is_active, created_at, updated_at
		FROM courses 
		WHERE is_active = true
		ORDER BY created_at DESC
	`

	rows, err := r.db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []*models.Course
	for rows.Next() {
		course := &models.Course{}
		var idStr string
		var imageURL sql.NullString
		err := rows.Scan(
			&idStr, &course.Title, &course.Description, &course.Price,
			&imageURL, &course.IsActive, &course.CreatedAt, &course.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Parse UUID from string
		course.ID, err = uuid.Parse(idStr)
		if err != nil {
			return nil, err
		}

		// Handle nullable image_url
		if imageURL.Valid {
			course.ImageURL = imageURL.String
		} else {
			course.ImageURL = ""
		}

		courses = append(courses, course)
	}

	return courses, nil
}

func (r *CourseRepository) GetCourseByID(id string) (*models.Course, error) {
	query := `
		SELECT id, title, description, price, image_url, is_active, created_at, updated_at
		FROM courses 
		WHERE id = $1 AND is_active = true
	`

	course := &models.Course{}
	err := r.db.DB.QueryRow(query, id).Scan(
		&course.ID, &course.Title, &course.Description, &course.Price,
		&course.ImageURL, &course.IsActive, &course.CreatedAt, &course.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return course, err
}

func (r *CourseRepository) CreateLesson(lesson *models.Lesson) error {
	query := `
		INSERT INTO lessons (course_id, title, content, video_url, "order", is_locked, duration, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.DB.Exec(query,
		lesson.CourseID, lesson.Title, lesson.Content, lesson.VideoURL,
		lesson.Order, lesson.IsLocked, lesson.Duration, time.Now(), time.Now(),
	)

	return err
}

func (r *CourseRepository) GetLessonsByCourseID(courseID string) ([]*models.Lesson, error) {
	query := `
		SELECT id, course_id, title, content, video_url, "order", is_locked, duration, created_at, updated_at
		FROM lessons 
		WHERE course_id = $1
		ORDER BY "order" ASC
	`

	rows, err := r.db.DB.Query(query, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lessons []*models.Lesson
	for rows.Next() {
		lesson := &models.Lesson{}
		err := rows.Scan(
			&lesson.ID, &lesson.CourseID, &lesson.Title, &lesson.Content,
			&lesson.VideoURL, &lesson.Order, &lesson.IsLocked, &lesson.Duration,
			&lesson.CreatedAt, &lesson.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		lessons = append(lessons, lesson)
	}

	return lessons, nil
}

func (r *CourseRepository) GetLessonByID(lessonID string) (*models.Lesson, error) {
	query := `
		SELECT id, course_id, title, content, video_url, "order", is_locked, duration, created_at, updated_at
		FROM lessons 
		WHERE id = $1
	`

	lesson := &models.Lesson{}
	err := r.db.DB.QueryRow(query, lessonID).Scan(
		&lesson.ID, &lesson.CourseID, &lesson.Title, &lesson.Content,
		&lesson.VideoURL, &lesson.Order, &lesson.IsLocked, &lesson.Duration,
		&lesson.CreatedAt, &lesson.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return lesson, err
}

func (r *CourseRepository) CreateOrUpdateProgress(progress *models.CourseProgress) error {
	query := `
		INSERT INTO course_progress (user_id, course_id, progress, completed, started_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (user_id, course_id) DO UPDATE SET
			progress = EXCLUDED.progress,
			completed = EXCLUDED.completed,
			updated_at = NOW()
	`

	_, err := r.db.DB.Exec(query,
		progress.UserID, progress.CourseID, progress.Progress,
		progress.Completed, progress.StartedAt, time.Now(),
	)

	return err
}

func (r *CourseRepository) GetUserProgress(userID, courseID string) (*models.CourseProgress, error) {
	query := `
		SELECT id, user_id, course_id, progress, completed, started_at, updated_at
		FROM course_progress 
		WHERE user_id = $1 AND course_id = $2
	`

	progress := &models.CourseProgress{}
	err := r.db.DB.QueryRow(query, userID, courseID).Scan(
		&progress.ID, &progress.UserID, &progress.CourseID, &progress.Progress,
		&progress.Completed, &progress.StartedAt, &progress.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return progress, err
}

func (r *CourseRepository) GetUserCourses(userID string) ([]*models.CourseProgress, error) {
	query := `
		SELECT cp.id, cp.user_id, cp.course_id, cp.progress, cp.completed, cp.started_at, cp.updated_at,
			   c.title, c.description, c.price, c.image_url
		FROM course_progress cp
		JOIN courses c ON cp.course_id = c.id
		WHERE cp.user_id = $1
		ORDER BY cp.started_at DESC
	`

	rows, err := r.db.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var progresses []*models.CourseProgress
	for rows.Next() {
		progress := &models.CourseProgress{}
		err := rows.Scan(
			&progress.ID, &progress.UserID, &progress.CourseID, &progress.Progress,
			&progress.Completed, &progress.StartedAt, &progress.UpdatedAt,
			new(string), new(string), new(float64), new(string), // Skip course fields for now
		)
		if err != nil {
			return nil, err
		}
		progresses = append(progresses, progress)
	}

	return progresses, nil
}
