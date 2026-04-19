package repositories

import (
	"fmt"
	"time"

	"iiceekiingfx.com/internal/models"
)

type PortfolioRepository struct {
	db *Database
}

func NewPortfolioRepository(db *Database) *PortfolioRepository {
	return &PortfolioRepository{db: db}
}

func (r *PortfolioRepository) CreateTradingAccount(account *models.TradingAccount) error {
	query := `
		INSERT INTO trading_accounts 
		(user_id, broker_name, account_id, account_type, currency, leverage, 
		 balance, equity, margin, free_margin, is_active, last_sync, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		ON CONFLICT (user_id, account_id) DO UPDATE SET
			broker_name = EXCLUDED.broker_name,
			account_type = EXCLUDED.account_type,
			currency = EXCLUDED.currency,
			leverage = EXCLUDED.leverage,
			updated_at = NOW()
	`

	_, err := r.db.DB.Exec(query,
		account.UserID, account.BrokerName, account.AccountID, account.AccountType,
		account.Currency, account.Leverage, account.Balance, account.Equity,
		account.Margin, account.FreeMargin, account.IsActive, time.Now(),
		time.Now(), time.Now(),
	)

	return err
}

func (r *PortfolioRepository) GetUserAccounts(userID string) ([]*models.TradingAccount, error) {
	query := `
		SELECT id, user_id, broker_name, account_id, account_type, currency, leverage,
			   balance, equity, margin, free_margin, is_active, last_sync, created_at, updated_at
		FROM trading_accounts 
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []*models.TradingAccount
	for rows.Next() {
		account := &models.TradingAccount{}
		err := rows.Scan(
			&account.ID, &account.UserID, &account.BrokerName, &account.AccountID,
			&account.AccountType, &account.Currency, &account.Leverage,
			&account.Balance, &account.Equity, &account.Margin, &account.FreeMargin,
			&account.IsActive, &account.LastSync, &account.CreatedAt, &account.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (r *PortfolioRepository) GetActiveAccountsCount(userID string) (int, error) {
	query := `SELECT COUNT(*) FROM trading_accounts WHERE user_id = $1 AND is_active = true`

	var count int
	err := r.db.DB.QueryRow(query, userID).Scan(&count)
	return count, err
}

func (r *PortfolioRepository) GetTotalProfit(userID string) (float64, error) {
	query := `
		SELECT COALESCE(SUM(profit), 0) 
		FROM trades t
		JOIN trading_accounts ta ON t.account_id = ta.id
		WHERE ta.user_id = $1 AND t.close_time IS NOT NULL
	`

	var totalProfit float64
	err := r.db.DB.QueryRow(query, userID).Scan(&totalProfit)
	return totalProfit, err
}

func (r *PortfolioRepository) GetTradeStats(userID string) (int, int, float64, error) {
	query := `
		SELECT 
			COUNT(*) as total_trades,
			COUNT(CASE WHEN profit > 0 THEN 1 END) as winning_trades,
			COALESCE(SUM(profit), 0) as total_profit
		FROM trades t
		JOIN trading_accounts ta ON t.account_id = ta.id
		WHERE ta.user_id = $1 AND t.close_time IS NOT NULL
	`

	var totalTrades, winningTrades int
	var totalProfit float64
	err := r.db.DB.QueryRow(query, userID).Scan(&totalTrades, &winningTrades, &totalProfit)

	return totalTrades, winningTrades, totalProfit, err
}

func (r *PortfolioRepository) GetEquityHistory(userID string, days int) ([]*models.EquityCurveData, error) {
	query := `
		SELECT DISTINCT 
			DATE(t.close_time) as date,
			COALESCE(SUM(
				CASE 
					WHEN t.close_time IS NOT NULL THEN t.profit 
					ELSE 0 
				END
			) OVER (PARTITION BY DATE(t.close_time) ORDER BY t.close_time), 0) as equity
		FROM trades t
		JOIN trading_accounts ta ON t.account_id = ta.id
		WHERE ta.user_id = $1 
			AND t.close_time >= NOW() - INTERVAL '%d days'
		ORDER BY date ASC
	`

	rows, err := r.db.DB.Query(fmt.Sprintf(query, days), userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var equityData []*models.EquityCurveData
	for rows.Next() {
		data := &models.EquityCurveData{}
		err := rows.Scan(&data.Date, &data.Equity)
		if err != nil {
			return nil, err
		}
		equityData = append(equityData, data)
	}

	return equityData, nil
}
