package memory

import (
	"bloomo-exam-api/domain/portfolio"
	"bloomo-exam-api/domain/shared/vo"
	"fmt"
)

type portfolioData struct {
	UserID          int
	TargetPortfolio map[string]int
}

var portfolioDataList = []portfolioData{
	{
		UserID:          1,
		TargetPortfolio: map[string]int{"A": 40, "B": 60},
	},
	{
		UserID:          2,
		TargetPortfolio: map[string]int{"E": 100},
	},
	{
		UserID:          3,
		TargetPortfolio: map[string]int{"A": 31, "B": 40, "E": 29},
	},
	{
		UserID:          4,
		TargetPortfolio: map[string]int{"B": 50, "C": 49, "D": 1},
	},
}

type PortfolioRepository struct{}

func NewPortfolioRepository() *PortfolioRepository {
	return &PortfolioRepository{}
}

var _ portfolio.Repository = (*PortfolioRepository)(nil)

func (r *PortfolioRepository) FindByUserID(userID vo.UserID) (*portfolio.Portfolio, error) {
	for _, d := range portfolioDataList {
		if d.UserID == userID.Value() {
			return toPortfolioModel(d)
		}
	}
	return nil, portfolio.ErrNotFound
}

func toPortfolioModel(d portfolioData) (*portfolio.Portfolio, error) {
	userID, err := vo.NewUserID(d.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id: %w", err)
	}

	targetPortfolio := make(map[vo.Ticker]vo.Ratio, len(d.TargetPortfolio))
	for tickerStr, ratioInt := range d.TargetPortfolio {
		ticker, err := vo.NewTicker(tickerStr)
		if err != nil {
			return nil, fmt.Errorf("invalid ticker: %w", err)
		}
		ratio, err := vo.NewRatio(ratioInt)
		if err != nil {
			return nil, fmt.Errorf("invalid ratio: %w", err)
		}
		targetPortfolio[ticker] = ratio
	}

	return &portfolio.Portfolio{
		UserID:          userID,
		TargetPortfolio: targetPortfolio,
	}, nil
}