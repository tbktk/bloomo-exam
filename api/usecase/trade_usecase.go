package usecase

import (
	"bloomo-exam-api/domain/portfolio"
	"bloomo-exam-api/domain/shared/vo"
	"bloomo-exam-api/domain/stock"
	"bloomo-exam-api/domain/trade"
	"fmt"
)

type TradeUsecase struct {
	stockRepo     stock.Repository
	portfolioRepo portfolio.Repository
	calculator    *trade.OrderCalculator
}

func NewTradeUsecase(
	stockRepo stock.Repository,
	portfolioRepo portfolio.Repository,
	calculator *trade.OrderCalculator,
) *TradeUsecase {
	return &TradeUsecase{
		stockRepo:     stockRepo,
		portfolioRepo: portfolioRepo,
		calculator:    calculator,
	}
}

type TradeInput struct {
	UserID int
	Amount int
}

type TradeOutput struct {
	Amount          int
	TargetPortfolio map[string]int
	Orders          []OrderOutput
}

type OrderOutput struct {
	Symbol   string
	Amount   int
	Quantity float64
}

func (u *TradeUsecase) Execute(input TradeInput) (*TradeOutput, error) {
	// UserID値オブジェクト生成
	userID, err := vo.NewUserID(input.UserID)
	if err != nil {
		return nil, ValidationError{
			Field:   "user_id",
			Message: err.Error(),
		}
	}

	// Amount値オブジェクト生成（最低取引金額チェック）
	amount, err := vo.NewAmount(input.Amount)
	if err != nil {
		return nil, ValidationError{
			Field:   "amount",
			Message: err.Error(),
		}
	}

	// ポートフォリオ取得
	pf, err := u.portfolioRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	// 銘柄一覧取得
	stocks, err := u.stockRepo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch stocks: %w", err)
	}

	// 注文計算（ドメインサービス委譲）
	orders, err := u.calculator.Calculate(amount, pf, stocks)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate orders: %w", err)
	}

	// 出力変換
	return toTradeOutput(input.Amount, pf, orders), nil
}

func toTradeOutput(amount int, pf *portfolio.Portfolio, orders []*trade.Order) *TradeOutput {
	targetPortfolio := make(map[string]int, len(pf.TargetPortfolio))
	for ticker, ratio := range pf.TargetPortfolio {
		targetPortfolio[ticker.Value()] = ratio.Value()
	}

	orderOutputs := make([]OrderOutput, 0, len(orders))
	for _, o := range orders {
		orderOutputs = append(orderOutputs, OrderOutput{
			Symbol:   o.Symbol.Value(),
			Amount:   o.Amount.Value(),
			Quantity: o.Quantity.Value(),
		})
	}

	return &TradeOutput{
		Amount:          amount,
		TargetPortfolio: targetPortfolio,
		Orders:          orderOutputs,
	}
}