package repositories

import (
	"database/sql"
	"time"

	"iiceekiingfx.com/internal/models"
)

type SignalsRepository struct {
	db *Database
}

func NewSignalsRepository(db *Database) *SignalsRepository {
	return &SignalsRepository{db: db}
}

func (r *SignalsRepository) CreateSignal(signal *models.Signal) error {
	query := `
		INSERT INTO signals (pair, action, entry_price, stop_loss, take_profit, risk_reward, status, notes, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	
	_, err := r.db.DB.Exec(query,
		signal.Pair, signal.Action, signal.EntryPrice, signal.StopLoss,
		signal.TakeProfit, signal.RiskReward, signal.Status, signal.Notes,
		signal.CreatedBy, time.Now(), time.Now(),
	)
	
	return err
}

func (r *SignalsRepository) GetAllSignals(limit, offset int) ([]*models.Signal, error) {
	query := `
		SELECT id, pair, action, entry_price, stop_loss, take_profit, risk_reward, status, notes, created_by, created_at, updated_at
		FROM signals 
		WHERE status = 'active'
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	
	rows, err := r.db.DB.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var signals []*models.Signal
	for rows.Next() {
		signal := &models.Signal{}
		err := rows.Scan(
			&signal.ID, &signal.Pair, &signal.Action, &signal.EntryPrice,
			&signal.StopLoss, &signal.TakeProfit, &signal.RiskReward, &signal.Status,
			&signal.Notes, &signal.CreatedBy, &signal.CreatedAt, &signal.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		signals = append(signals, signal)
	}
	
	return signals, nil
}

func (r *SignalsRepository) GetSignalByID(id string) (*models.Signal, error) {
	query := `
		SELECT id, pair, action, entry_price, stop_loss, take_profit, risk_reward, status, notes, created_by, created_at, updated_at
		FROM signals 
		WHERE id = $1
	`
	
	signal := &models.Signal{}
	err := r.db.DB.QueryRow(query, id).Scan(
		&signal.ID, &signal.Pair, &signal.Action, &signal.EntryPrice,
		&signal.StopLoss, &signal.TakeProfit, &signal.RiskReward, &signal.Status,
		&signal.Notes, &signal.CreatedBy, &signal.CreatedAt, &signal.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	
	return signal, err
}

func (r *SignalsRepository) UpdateSignal(signal *models.Signal) error {
	query := `
		UPDATE signals 
		SET pair = $2, action = $3, entry_price = $4, stop_loss = $5, 
			take_profit = $6, risk_reward = $7, status = $8, notes = $9, updated_at = $10
		WHERE id = $1
	`
	
	_, err := r.db.DB.Exec(query,
		signal.ID, signal.Pair, signal.Action, signal.EntryPrice, signal.StopLoss,
		signal.TakeProfit, signal.RiskReward, signal.Status, signal.Notes, time.Now(),
	)
	
	return err
}

func (r *SignalsRepository) DeleteSignal(id string) error {
	query := `DELETE FROM signals WHERE id = $1`
	
	_, err := r.db.DB.Exec(query, id)
	return err
}

func (r *SignalsRepository) GetSignalsCount() (int, error) {
	query := `SELECT COUNT(*) FROM signals WHERE status = 'active'`
	
	var count int
	err := r.db.DB.QueryRow(query).Scan(&count)
	return count, err
}

func (r *SignalsRepository) GetSignalsByPair(pair string) ([]*models.Signal, error) {
	query := `
		SELECT id, pair, action, entry_price, stop_loss, take_profit, risk_reward, status, notes, created_by, created_at, updated_at
		FROM signals 
		WHERE pair = $1 AND status = 'active'
		ORDER BY created_at DESC
	`
	
	rows, err := r.db.DB.Query(query, pair)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var signals []*models.Signal
	for rows.Next() {
		signal := &models.Signal{}
		err := rows.Scan(
			&signal.ID, &signal.Pair, &signal.Action, &signal.EntryPrice,
			&signal.StopLoss, &signal.TakeProfit, &signal.RiskReward, &signal.Status,
			&signal.Notes, &signal.CreatedBy, &signal.CreatedAt, &signal.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		signals = append(signals, signal)
	}
	
	return signals, nil
}
