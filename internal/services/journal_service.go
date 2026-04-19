package services

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"iiceekiingfx.com/internal/models"
	"iiceekiingfx.com/internal/repositories"
)

type JournalService struct {
	journalRepo *repositories.JournalRepository
}

func NewJournalService(journalRepo *repositories.JournalRepository) *JournalService {
	return &JournalService{
		journalRepo: journalRepo,
	}
}

func (s *JournalService) GetUserEntries(userID string, limit, offset int) ([]*models.TradeJournal, error) {
	return s.journalRepo.GetUserJournalEntries(userID, limit, offset)
}

func (s *JournalService) CreateEntry(userID string, req *models.JournalRequest) (*models.TradeJournal, error) {
	// Parse trade date
	tradeDate, err := time.Parse("2006-01-02", req.TradeDate)
	if err != nil {
		return nil, err
	}

	entry := &models.TradeJournal{
		ID:          uuid.New(),
		UserID:      uuid.MustParse(userID),
		Pair:        req.Pair,
		EntryPrice:  req.EntryPrice,
		ExitPrice:   req.ExitPrice,
		LotSize:     req.LotSize,
		Result:      req.Result,
		RMultiple:   req.RMultiple,
		Notes:       req.Notes,
		StrategyTag: req.StrategyTag,
		TradeDate:   tradeDate,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.journalRepo.CreateJournalEntry(entry); err != nil {
		return nil, err
	}

	return entry, nil
}

func (s *JournalService) UpdateEntry(userID, entryID string, req *models.JournalRequest) (*models.TradeJournal, error) {
	// Get existing entry
	existingEntry, err := s.journalRepo.GetJournalEntryByID(entryID, userID)
	if err != nil {
		return nil, err
	}
	if existingEntry == nil {
		return nil, errors.New("journal entry not found")
	}

	// Parse trade date
	tradeDate, err := time.Parse("2006-01-02", req.TradeDate)
	if err != nil {
		return nil, err
	}

	// Update entry
	existingEntry.Pair = req.Pair
	existingEntry.EntryPrice = req.EntryPrice
	existingEntry.ExitPrice = req.ExitPrice
	existingEntry.LotSize = req.LotSize
	existingEntry.Result = req.Result
	existingEntry.RMultiple = req.RMultiple
	existingEntry.Notes = req.Notes
	existingEntry.StrategyTag = req.StrategyTag
	existingEntry.TradeDate = tradeDate
	existingEntry.UpdatedAt = time.Now()

	if err := s.journalRepo.UpdateJournalEntry(existingEntry); err != nil {
		return nil, err
	}

	return existingEntry, nil
}

func (s *JournalService) DeleteEntry(userID, entryID string) error {
	// Verify entry exists and belongs to user
	existingEntry, err := s.journalRepo.GetJournalEntryByID(entryID, userID)
	if err != nil {
		return err
	}
	if existingEntry == nil {
		return errors.New("journal entry not found")
	}

	return s.journalRepo.DeleteJournalEntry(entryID, userID)
}

func (s *JournalService) GetStats(userID string) (int, int, float64, float64, error) {
	return s.journalRepo.GetJournalStats(userID)
}
