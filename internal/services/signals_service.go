package services

import (
	"errors"
	"time"

	"iiceekiingfx.com/internal/models"
	"iiceekiingfx.com/internal/repositories"
	"github.com/google/uuid"
)

type SignalsService struct {
	signalsRepo *repositories.SignalsRepository
}

func NewSignalsService(signalsRepo *repositories.SignalsRepository) *SignalsService {
	return &SignalsService{
		signalsRepo: signalsRepo,
	}
}

func (s *SignalsService) GetAllSignals(page, limit int) ([]*models.Signal, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	signals, err := s.signalsRepo.GetAllSignals(limit, offset)
	if err != nil {
		return nil, 0, err
	}

	totalCount, err := s.signalsRepo.GetSignalsCount()
	if err != nil {
		return nil, 0, err
	}

	return signals, totalCount, nil
}

func (s *SignalsService) GetSignalByID(signalID string) (*models.Signal, error) {
	if signalID == "" {
		return nil, errors.New("signal ID is required")
	}

	signal, err := s.signalsRepo.GetSignalByID(signalID)
	if err != nil {
		return nil, err
	}
	if signal == nil {
		return nil, errors.New("signal not found")
	}

	return signal, nil
}

func (s *SignalsService) CreateSignal(userID string, req *models.SignalRequest) (*models.Signal, error) {
	if req.Pair == "" || req.Action == "" {
		return nil, errors.New("pair and action are required")
	}

	if req.EntryPrice <= 0 || req.StopLoss <= 0 || req.TakeProfit <= 0 {
		return nil, errors.New("prices must be greater than 0")
	}

	if req.RiskReward <= 0 {
		return nil, errors.New("risk/reward ratio must be greater than 0")
	}

	signal := &models.Signal{
		ID:         uuid.New(),
		Pair:       req.Pair,
		Action:     req.Action,
		EntryPrice: req.EntryPrice,
		StopLoss:   req.StopLoss,
		TakeProfit: req.TakeProfit,
		RiskReward: req.RiskReward,
		Status:     "active",
		Notes:      req.Notes,
		CreatedBy:  uuid.MustParse(userID),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err := s.signalsRepo.CreateSignal(signal)
	if err != nil {
		return nil, err
	}

	return signal, nil
}

func (s *SignalsService) UpdateSignal(userID, signalID string, req *models.SignalRequest) (*models.Signal, error) {
	if signalID == "" {
		return nil, errors.New("signal ID is required")
	}

	// Get existing signal
	existingSignal, err := s.signalsRepo.GetSignalByID(signalID)
	if err != nil {
		return nil, err
	}
	if existingSignal == nil {
		return nil, errors.New("signal not found")
	}

	// Update signal
	existingSignal.Pair = req.Pair
	existingSignal.Action = req.Action
	existingSignal.EntryPrice = req.EntryPrice
	existingSignal.StopLoss = req.StopLoss
	existingSignal.TakeProfit = req.TakeProfit
	existingSignal.RiskReward = req.RiskReward
	existingSignal.Notes = req.Notes
	existingSignal.UpdatedAt = time.Now()

	err = s.signalsRepo.UpdateSignal(existingSignal)
	if err != nil {
		return nil, err
	}

	return existingSignal, nil
}

func (s *SignalsService) DeleteSignal(userID, signalID string) error {
	if signalID == "" {
		return errors.New("signal ID is required")
	}

	// Get existing signal to verify ownership
	existingSignal, err := s.signalsRepo.GetSignalByID(signalID)
	if err != nil {
		return err
	}
	if existingSignal == nil {
		return errors.New("signal not found")
	}

	// Check if user is the creator (in a real system, you'd check user role)
	if existingSignal.CreatedBy.String() != userID {
		return errors.New("not authorized to delete this signal")
	}

	err = s.signalsRepo.DeleteSignal(signalID)
	if err != nil {
		return err
	}

	return nil
}

func (s *SignalsService) GetSignalsByPair(pair string) ([]*models.Signal, error) {
	if pair == "" {
		return nil, errors.New("pair is required")
	}

	signals, err := s.signalsRepo.GetSignalsByPair(pair)
	if err != nil {
		return nil, err
	}

	return signals, nil
}

func (s *SignalsService) CloseSignal(signalID string) error {
	if signalID == "" {
		return errors.New("signal ID is required")
	}

	// Get existing signal
	existingSignal, err := s.signalsRepo.GetSignalByID(signalID)
	if err != nil {
		return err
	}
	if existingSignal == nil {
		return errors.New("signal not found")
	}

	// Update status to closed
	existingSignal.Status = "closed"
	existingSignal.UpdatedAt = time.Now()

	err = s.signalsRepo.UpdateSignal(existingSignal)
	if err != nil {
		return err
	}

	return nil
}
