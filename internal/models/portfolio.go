package models

import (
	"time"

	"github.com/google/uuid"
)

type TradingAccount struct {
	ID           uuid.UUID `json:"id" db:"id"`
	UserID       uuid.UUID `json:"user_id" db:"user_id"`
	BrokerName   string    `json:"broker_name" db:"broker_name"`
	AccountID    string    `json:"account_id" db:"account_id"`
	AccountType  string    `json:"account_type" db:"account_type"`
	Currency     string    `json:"currency" db:"currency"`
	Leverage     int       `json:"leverage" db:"leverage"`
	Balance      float64   `json:"balance" db:"balance"`
	Equity       float64   `json:"equity" db:"equity"`
	Margin       float64   `json:"margin" db:"margin"`
	FreeMargin   float64   `json:"free_margin" db:"free_margin"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	LastSync     time.Time `json:"last_sync" db:"last_sync"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type Trade struct {
	ID           uuid.UUID `json:"id" db:"id"`
	AccountID    uuid.UUID `json:"account_id" db:"account_id"`
	Ticket       int64     `json:"ticket" db:"ticket"`
	Symbol       string    `json:"symbol" db:"symbol"`
	Type         string    `json:"type" db:"type"` // BUY, SELL
	Volume       float64   `json:"volume" db:"volume"`
	OpenPrice    float64   `json:"open_price" db:"open_price"`
	ClosePrice   float64   `json:"close_price" db:"close_price"`
	OpenTime     time.Time `json:"open_time" db:"open_time"`
	CloseTime    time.Time `json:"close_time" db:"close_time"`
	Profit       float64   `json:"profit" db:"profit"`
	Commission   float64   `json:"commission" db:"commission"`
	Swap         float64   `json:"swap" db:"swap"`
	Comment      string    `json:"comment" db:"comment"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

type ConnectAccountRequest struct {
	BrokerName  string `json:"broker_name" validate:"required"`
	AccountID   string `json:"account_id" validate:"required"`
	Login       string `json:"login" validate:"required"`
	Password    string `json:"password" validate:"required"`
	Server      string `json:"server" validate:"required"`
	AccountType string `json:"account_type" validate:"required"`
	Currency    string `json:"currency" validate:"required"`
	Leverage    int    `json:"leverage" validate:"required,min=1,max=1000"`
}
