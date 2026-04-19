package models

import (
	"time"

	"github.com/google/uuid"
)

type Course struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Price       float64   `json:"price" db:"price"`
	ImageURL    string    `json:"image_url" db:"image_url"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type Lesson struct {
	ID         uuid.UUID `json:"id" db:"id"`
	CourseID   uuid.UUID `json:"course_id" db:"course_id"`
	Title      string    `json:"title" db:"title"`
	Content    string    `json:"content" db:"content"`
	VideoURL   string    `json:"video_url" db:"video_url"`
	Order      int       `json:"order" db:"order"`
	IsLocked   bool      `json:"is_locked" db:"is_locked"`
	Duration   int       `json:"duration" db:"duration"` // in minutes
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type CourseProgress struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	CourseID  uuid.UUID `json:"course_id" db:"course_id"`
	Progress  float64   `json:"progress" db:"progress"` // 0-100
	Completed bool      `json:"completed" db:"completed"`
	StartedAt time.Time `json:"started_at" db:"started_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Signal struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Pair        string    `json:"pair" db:"pair"`
	Action      string    `json:"action" db:"action"` // "BUY", "SELL"
	EntryPrice  float64   `json:"entry_price" db:"entry_price"`
	StopLoss    float64   `json:"stop_loss" db:"stop_loss"`
	TakeProfit  float64   `json:"take_profit" db:"take_profit"`
	RiskReward  float64   `json:"risk_reward" db:"risk_reward"`
	Status      string    `json:"status" db:"status"` // "active", "closed", "cancelled"
	Notes       string    `json:"notes" db:"notes"`
	CreatedBy   uuid.UUID `json:"created_by" db:"created_by"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type SignalRequest struct {
	Pair       string  `json:"pair" validate:"required"`
	Action     string  `json:"action" validate:"required,oneof=BUY SELL"`
	EntryPrice float64 `json:"entry_price" validate:"required,gt=0"`
	StopLoss   float64 `json:"stop_loss" validate:"required,gt=0"`
	TakeProfit float64 `json:"take_profit" validate:"required,gt=0"`
	RiskReward float64 `json:"risk_reward" validate:"required,gt=0"`
	Notes      string  `json:"notes"`
}
