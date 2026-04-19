package repositories

import (
	"iiceekiingfx.com/internal/models"
)

type JournalRepository struct {
	db *Database
}

func NewJournalRepository(db *Database) *JournalRepository {
	return &JournalRepository{db: db}
}

func (r *JournalRepository) CreateJournalEntry(entry *models.TradeJournal) error {
	query := `
		INSERT INTO trade_journal 
		(user_id, pair, entry_price, exit_price, lot_size, result, r_multiple, notes, strategy_tag, trade_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW())
	`
	
	_, err := r.db.DB.Exec(query,
		entry.UserID, entry.Pair, entry.EntryPrice, entry.ExitPrice,
		entry.LotSize, entry.Result, entry.RMultiple, entry.Notes,
		entry.StrategyTag, entry.TradeDate,
	)
	
	return err
}

func (r *JournalRepository) GetUserJournalEntries(userID string, limit, offset int) ([]*models.TradeJournal, error) {
	query := `
		SELECT id, user_id, pair, entry_price, exit_price, lot_size, result, r_multiple, 
			   notes, strategy_tag, trade_date, created_at, updated_at
		FROM trade_journal 
		WHERE user_id = $1
		ORDER BY trade_date DESC, created_at DESC
		LIMIT $2 OFFSET $3
	`
	
	rows, err := r.db.DB.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var entries []*models.TradeJournal
	for rows.Next() {
		entry := &models.TradeJournal{}
		err := rows.Scan(
			&entry.ID, &entry.UserID, &entry.Pair, &entry.EntryPrice, &entry.ExitPrice,
			&entry.LotSize, &entry.Result, &entry.RMultiple, &entry.Notes,
			&entry.StrategyTag, &entry.TradeDate, &entry.CreatedAt, &entry.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	
	return entries, nil
}

func (r *JournalRepository) GetJournalEntryByID(id, userID string) (*models.TradeJournal, error) {
	query := `
		SELECT id, user_id, pair, entry_price, exit_price, lot_size, result, r_multiple, 
			   notes, strategy_tag, trade_date, created_at, updated_at
		FROM trade_journal 
		WHERE id = $1 AND user_id = $2
	`
	
	entry := &models.TradeJournal{}
	err := r.db.DB.QueryRow(query, id, userID).Scan(
		&entry.ID, &entry.UserID, &entry.Pair, &entry.EntryPrice, &entry.ExitPrice,
		&entry.LotSize, &entry.Result, &entry.RMultiple, &entry.Notes,
		&entry.StrategyTag, &entry.TradeDate, &entry.CreatedAt, &entry.UpdatedAt,
	)
	
	return entry, err
}

func (r *JournalRepository) UpdateJournalEntry(entry *models.TradeJournal) error {
	query := `
		UPDATE trade_journal 
		SET pair = $2, entry_price = $3, exit_price = $4, lot_size = $5, 
			result = $6, r_multiple = $7, notes = $8, strategy_tag = $9, 
			trade_date = $10, updated_at = NOW()
		WHERE id = $1 AND user_id = $11
	`
	
	_, err := r.db.DB.Exec(query,
		entry.ID, entry.Pair, entry.EntryPrice, entry.ExitPrice, entry.LotSize,
		entry.Result, entry.RMultiple, entry.Notes, entry.StrategyTag,
		entry.TradeDate, entry.UserID,
	)
	
	return err
}

func (r *JournalRepository) DeleteJournalEntry(id, userID string) error {
	query := `DELETE FROM trade_journal WHERE id = $1 AND user_id = $2`
	
	_, err := r.db.DB.Exec(query, id, userID)
	return err
}

func (r *JournalRepository) GetJournalStats(userID string) (int, int, float64, float64, error) {
	query := `
		SELECT 
			COUNT(*) as total_trades,
			COUNT(CASE WHEN result = 'win' THEN 1 END) as winning_trades,
			COALESCE(AVG(CASE WHEN result = 'win' THEN r_multiple ELSE NULL END), 0) as avg_r_multiple,
			COALESCE(SUM(CASE WHEN result = 'win' THEN 1 ELSE 0 END)::DECIMAL / NULLIF(COUNT(*), 0), 0) as win_rate
		FROM trade_journal 
		WHERE user_id = $1
	`
	
	var totalTrades, winningTrades int
	var avgRMultiple, winRate float64
	err := r.db.DB.QueryRow(query, userID).Scan(&totalTrades, &winningTrades, &avgRMultiple, &winRate)
	
	return totalTrades, winningTrades, avgRMultiple, winRate, err
}
