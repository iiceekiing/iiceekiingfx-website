package models

import (
	"time"

	"github.com/google/uuid"
)

type DashboardOverview struct {
	UserID         uuid.UUID `json:"user_id"`
	TotalProfit    float64   `json:"total_profit"`
	WinRate        float64   `json:"win_rate"`
	ActiveAccounts int       `json:"active_accounts"`
	CourseProgress float64   `json:"course_progress"`
	TotalTrades    int       `json:"total_trades"`
	ProfitFactor   float64   `json:"profit_factor"`
	Drawdown       float64   `json:"drawdown"`
	Expectancy     float64   `json:"expectancy"`
}

type EquityCurveData struct {
	Date   string  `json:"date"`
	Equity float64 `json:"equity"`
}

type ActivityFeed struct {
	ID        uuid.UUID `json:"id"`
	Type      string    `json:"type"` // "trade", "course", "signal"
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Data      interface{} `json:"data,omitempty"`
}

type TradeJournal struct {
	ID           uuid.UUID `json:"id" db:"id"`
	UserID       uuid.UUID `json:"user_id" db:"user_id"`
	Pair         string    `json:"pair" db:"pair"`
	EntryPrice   float64   `json:"entry_price" db:"entry_price"`
	ExitPrice    float64   `json:"exit_price" db:"exit_price"`
	LotSize      float64   `json:"lot_size" db:"lot_size"`
	Result       string    `json:"result" db:"result"` // "win", "loss", "breakeven"
	RMultiple    float64   `json:"r_multiple" db:"r_multiple"`
	Notes        string    `json:"notes" db:"notes"`
	StrategyTag  string    `json:"strategy_tag" db:"strategy_tag"`
	TradeDate    time.Time `json:"trade_date" db:"trade_date"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type JournalRequest struct {
	Pair        string  `json:"pair" validate:"required"`
	EntryPrice  float64 `json:"entry_price" validate:"required,gt=0"`
	ExitPrice   float64 `json:"exit_price" validate:"required,gt=0"`
	LotSize     float64 `json:"lot_size" validate:"required,gt=0"`
	Result      string  `json:"result" validate:"required,oneof=win loss breakeven"`
	RMultiple   float64 `json:"r_multiple" validate:"required"`
	Notes       string  `json:"notes"`
	StrategyTag string  `json:"strategy_tag"`
	TradeDate   string  `json:"trade_date" validate:"required"`
}
