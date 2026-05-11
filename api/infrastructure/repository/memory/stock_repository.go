package memory

import (
	"bloomo-exam-api/domain/shared/vo"
	"bloomo-exam-api/domain/stock"
	"fmt"
)

type stockData struct {
	Ticker   string
	Price    int
	Tradable bool
}

var stockDataList = []stockData{
	{Ticker: "A", Price: 1000, Tradable: true},
	{Ticker: "B", Price: 155, Tradable: true},
	{Ticker: "C", Price: 2222, Tradable: true},
	{Ticker: "D", Price: 467, Tradable: true},
	{Ticker: "E", Price: 888, Tradable: false},
}

type StockRepository struct{}

func NewStockRepository() *StockRepository {
	return &StockRepository{}
}

var _ stock.Repository = (*StockRepository)(nil)

func (r *StockRepository) FindAll() ([]*stock.Stock, error) {
	stocks := make([]*stock.Stock, 0, len(stockDataList))
	for _, d := range stockDataList {
		s, err := toStockModel(d)
		if err != nil {
			return nil, fmt.Errorf("failed to convert stock data: %w", err)
		}
		stocks = append(stocks, s)
	}
	return stocks, nil
}

func (r *StockRepository) FindByTicker(ticker vo.Ticker) (*stock.Stock, error) {
	for _, d := range stockDataList {
		if d.Ticker == ticker.Value() {
			s, err := toStockModel(d)
			if err != nil {
				return nil, fmt.Errorf("failed to convert stock data: %w", err)
			}
			return s, nil
		}
	}
	return nil, stock.ErrNotFound
}

func toStockModel(d stockData) (*stock.Stock, error) {
	ticker, err := vo.NewTicker(d.Ticker)
	if err != nil {
		return nil, err
	}
	price, err := vo.NewPrice(d.Price)
	if err != nil {
		return nil, err
	}
	return &stock.Stock{
		Ticker:   ticker,
		Price:    price,
		Tradable: d.Tradable,
	}, nil
}