package trade

import (
	"bloomo-exam-api/domain/portfolio"
	"bloomo-exam-api/domain/shared/vo"
	"bloomo-exam-api/domain/stock"
)

type OrderCalculator struct{}

func (c *OrderCalculator) Calculate(
    amount vo.Amount,
    portfolio *portfolio.Portfolio,
    stocks []*stock.Stock,
) ([]*Order, error) {
	// stocksをtickerでマップ化（検索用）
	stockMap := makeStockMap(stocks)

	// 1. ポートフォリオからtradable=falseを除外
	tradableTargets := filterTradable(portfolio.TargetPortfolio, stockMap)
	if len(tradableTargets) == 0 {
		return []*Order{}, nil
	}

	// 2. 按分計算（切り捨て）・200円未満を除外
	orders := calcOrders(amount, tradableTargets, stockMap)

	return orders, nil
}

// stocksをticker keyのmapに変換
func makeStockMap(stocks []*stock.Stock) map[vo.Ticker]*stock.Stock {
	m := make(map[vo.Ticker]*stock.Stock, len(stocks))
	for _, s := range stocks {
		m[s.Ticker] = s
	}
	return m
}

// tradable=trueかつstocksに存在する銘柄のみ残す
func filterTradable(
	targetPortfolio map[vo.Ticker]vo.Ratio,
	stockMap map[vo.Ticker]*stock.Stock,
) map[vo.Ticker]vo.Ratio {
	result := make(map[vo.Ticker]vo.Ratio)
	for ticker, ratio := range targetPortfolio {
		s, exists := stockMap[ticker]
		if !exists {
			continue
		}
		if !s.IsTradable() {
			continue
		}
		result[ticker] = ratio
	}
	return result
}

// 按分計算・200円未満除外・数量計算
func calcOrders(
	amount vo.Amount,
	targets map[vo.Ticker]vo.Ratio,
	stockMap map[vo.Ticker]*stock.Stock,
) []*Order {
	orders := make([]*Order, 0, len(targets))

	for ticker, ratio := range targets {
		// 按分金額計算（切り捨て）
		orderAmountValue := (amount.Value() * ratio.Value()) / 100

		// 200円未満は除外
		orderAmount, err := vo.NewOrderAmount(orderAmountValue)
		if err != nil {
			continue
		}

		// 数量計算
		s := stockMap[ticker]
		quantity, err := vo.NewQuantity(orderAmountValue, s.Price.Value())
		if err != nil {
			continue
		}

		orders = append(orders, &Order{
			Symbol:   ticker,
			Amount:   orderAmount,
			Quantity: quantity,
		})
	}

	return orders
}