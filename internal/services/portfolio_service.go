package services

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"iiceekiingfx.com/internal/models"
	"iiceekiingfx.com/internal/repositories"
	"iiceekiingfx.com/pkg/crypto"
)

type PortfolioService struct {
	portfolioRepo *repositories.PortfolioRepository
	encryptionSvc *crypto.EncryptionService
}

func NewPortfolioService(portfolioRepo *repositories.PortfolioRepository, encryptionSvc *crypto.EncryptionService) *PortfolioService {
	return &PortfolioService{
		portfolioRepo: portfolioRepo,
		encryptionSvc: encryptionSvc,
	}
}

func (s *PortfolioService) ConnectAccount(userID string, req *models.ConnectAccountRequest) (*models.TradingAccount, error) {
	// Validate input
	if req.BrokerName == "" || req.AccountID == "" || req.Login == "" || req.Password == "" {
		return nil, errors.New("all required fields must be provided")
	}

	// Check if account already exists
	accounts, err := s.portfolioRepo.GetUserAccounts(userID)
	if err != nil {
		return nil, err
	}

	for _, account := range accounts {
		if account.AccountID == req.AccountID {
			return nil, errors.New("account with this ID already exists")
		}
	}

	// Encrypt credentials (for future use in secure storage)
	_, err = s.encryptionSvc.EncryptMT5Credentials(req.Login, req.Password, req.Server)
	if err != nil {
		return nil, err
	}

	// Create trading account record
	account := &models.TradingAccount{
		ID:          uuid.New(),
		UserID:      uuid.MustParse(userID),
		BrokerName:  req.BrokerName,
		AccountID:   req.AccountID,
		AccountType: req.AccountType,
		Currency:    req.Currency,
		Leverage:    req.Leverage,
		Balance:     0.0,
		Equity:      0.0,
		Margin:      0.0,
		FreeMargin:  0.0,
		IsActive:    true,
		LastSync:    time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Note: In a real implementation, you would store encrypted credentials separately
	// For now, we'll create the account without storing credentials directly

	if err := s.portfolioRepo.CreateTradingAccount(account); err != nil {
		return nil, err
	}

	return account, nil
}

func (s *PortfolioService) GetUserAccounts(userID string) ([]*models.TradingAccount, error) {
	return s.portfolioRepo.GetUserAccounts(userID)
}

func (s *PortfolioService) GetAccountByID(userID, accountID string) (*models.TradingAccount, error) {
	accounts, err := s.portfolioRepo.GetUserAccounts(userID)
	if err != nil {
		return nil, err
	}

	for _, account := range accounts {
		if account.ID.String() == accountID {
			return account, nil
		}
	}

	return nil, errors.New("account not found")
}

func (s *PortfolioService) UpdateAccountStatus(userID, accountID string, isActive bool) error {
	// This would update the account status in the database
	// For now, we'll return a placeholder implementation
	return nil
}

func (s *PortfolioService) SyncAccount(userID, accountID string) error {
	// This would trigger MT5 synchronization
	// For now, we'll return a placeholder implementation
	return nil
}

func (s *PortfolioService) GetAccountPerformance(userID, accountID string) (map[string]interface{}, error) {
	// This would calculate performance metrics for the account
	// For now, we'll return placeholder data
	performance := map[string]interface{}{
		"total_trades":    0,
		"winning_trades":  0,
		"total_profit":    0.0,
		"win_rate":        0.0,
		"profit_factor":   0.0,
		"max_drawdown":    0.0,
		"current_balance": 0.0,
		"current_equity":  0.0,
	}

	return performance, nil
}
