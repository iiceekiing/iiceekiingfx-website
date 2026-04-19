package services

import (
	"time"

	"github.com/google/uuid"
	"iiceekiingfx.com/internal/models"
	"iiceekiingfx.com/internal/repositories"
)

type DashboardService struct {
	portfolioRepo *repositories.PortfolioRepository
	journalRepo   *repositories.JournalRepository
}

func NewDashboardService(portfolioRepo *repositories.PortfolioRepository, journalRepo *repositories.JournalRepository) *DashboardService {
	return &DashboardService{
		portfolioRepo: portfolioRepo,
		journalRepo:   journalRepo,
	}
}

func (s *DashboardService) GetOverview(userID string) (*models.DashboardOverview, error) {
	// Get active accounts count
	activeAccounts, err := s.portfolioRepo.GetActiveAccountsCount(userID)
	if err != nil {
		return nil, err
	}

	// Get total profit from trades
	totalProfit, err := s.portfolioRepo.GetTotalProfit(userID)
	if err != nil {
		return nil, err
	}

	// Get trade statistics
	totalTrades, winningTrades, _, err := s.portfolioRepo.GetTradeStats(userID)
	if err != nil {
		return nil, err
	}

	// Calculate win rate
	winRate := 0.0
	if totalTrades > 0 {
		winRate = float64(winningTrades) / float64(totalTrades) * 100
	}

	// Get journal statistics for additional metrics
	journalTotal, _, avgRMultiple, journalWinRate, err := s.journalRepo.GetJournalStats(userID)
	if err != nil {
		// If no journal entries, use trade stats
		journalTotal = totalTrades
		journalWinRate = winRate
	}

	// Calculate profit factor (simplified)
	profitFactor := 1.0
	if totalTrades > 0 {
		// This is a simplified calculation - in production, you'd sum wins and losses separately
		profitFactor = 1.5 // Placeholder
	}

	// Calculate expectancy
	expectancy := avgRMultiple*journalWinRate - (1.0 - journalWinRate)

	// Calculate drawdown (simplified)
	drawdown := 5.0 // Placeholder - would need historical equity data

	// Course progress (placeholder - would need course repository)
	courseProgress := 0.0

	userUUID, _ := uuid.Parse(userID)
	return &models.DashboardOverview{
		UserID:         userUUID,
		TotalProfit:    totalProfit,
		WinRate:        journalWinRate,
		ActiveAccounts: activeAccounts,
		CourseProgress: courseProgress,
		TotalTrades:    journalTotal,
		ProfitFactor:   profitFactor,
		Drawdown:       drawdown,
		Expectancy:     expectancy,
	}, nil
}

func (s *DashboardService) GetEquityCurve(userID string, days int) ([]*models.EquityCurveData, error) {
	return s.portfolioRepo.GetEquityHistory(userID, days)
}

func (s *DashboardService) GetActivityFeed(userID string, limit int) ([]*models.ActivityFeed, error) {
	// This would combine activities from different sources
	// For now, we'll create a placeholder implementation
	activities := []*models.ActivityFeed{
		{
			Type:      "trade",
			Title:     "Recent Trade",
			Message:   "EUR/USD trade closed with profit",
			Timestamp: time.Now(),
		},
		{
			Type:      "course",
			Title:     "Course Progress",
			Message:   "Completed lesson 3 of Forex Basics",
			Timestamp: time.Now().Add(-2 * time.Hour),
		},
		{
			Type:      "signal",
			Title:     "New Signal Available",
			Message:   "GBP/USD BUY signal detected",
			Timestamp: time.Now().Add(-4 * time.Hour),
		},
	}

	if len(activities) > limit {
		activities = activities[:limit]
	}

	return activities, nil
}
