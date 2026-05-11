package trade_test

import (
	"bloomo-exam-api/domain/portfolio"
	"bloomo-exam-api/domain/shared/vo"
	"bloomo-exam-api/domain/stock"
	"bloomo-exam-api/domain/trade"
	"testing"
)

// テスト用ヘルパー

func mustTicker(s string) vo.Ticker {
	t, err := vo.NewTicker(s)
	if err != nil {
		panic(err)
	}
	return t
}

func mustPrice(v int) vo.Price {
	p, err := vo.NewPrice(v)
	if err != nil {
		panic(err)
	}
	return p
}

func mustRatio(v int) vo.Ratio {
	r, err := vo.NewRatio(v)
	if err != nil {
		panic(err)
	}
	return r
}

func mustAmount(v int) vo.Amount {
	a, err := vo.NewAmount(v)
	if err != nil {
		panic(err)
	}
	return a
}

func makeStock(ticker string, price int, tradable bool) stock.Stock {
	return stock.Stock{
		Ticker:   mustTicker(ticker),
		Price:    mustPrice(price),
		Tradable: tradable,
	}
}

func makePortfolio(ratios map[string]int) portfolio.Portfolio {
	m := make(map[vo.Ticker]vo.Ratio, len(ratios))
	for k, v := range ratios {
		m[mustTicker(k)] = mustRatio(v)
	}
	return portfolio.Portfolio{TargetPortfolio: m}
}

// テスト本体

func TestOrderCalculator_Calculate(t *testing.T) {
	calc := &trade.OrderCalculator{}

	t.Run("正常系：全銘柄tradable、全て200円以上", func(t *testing.T) {
		// A:40%, B:60%, amount:10000
		// A: 10000*40/100=4000円, 数量=4000/1000=4.0
		// B: 10000*60/100=6000円, 数量=6000/155=38.709...→38.709
		stockA := makeStock("A", 1000, true)
		stockB := makeStock("B", 155, true)
		stocks := []*stock.Stock{&stockA, &stockB}
		pf := makePortfolio(map[string]int{"A": 40, "B": 60})
		amount := mustAmount(10000)

		orders, err := calc.Calculate(amount, &pf, stocks)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(orders) != 2 {
			t.Fatalf("expected 2 orders, got %d", len(orders))
		}

		assertOrder(t, orders, "A", 4000, 4.0)
		assertOrder(t, orders, "B", 6000, 38.709)
	})

	t.Run("tradable=falseの銘柄を除外して按分計算", func(t *testing.T) {
		// A:40%, B:60%, Bがtradable=false
		// Bを除外→Aのみ100%相当で按分
		// A: 10000*40/100=4000円, 数量=4000/1000=4.0
		// ※除外後に再按分はしない（仕様解釈）
		stockA := makeStock("A", 1000, true)
		stockB := makeStock("B", 155, false)
		stocks := []*stock.Stock{&stockA, &stockB}
		pf := makePortfolio(map[string]int{"A": 40, "B": 60})
		amount := mustAmount(10000)

		orders, err := calc.Calculate(amount, &pf, stocks)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(orders) != 1 {
			t.Fatalf("expected 1 order, got %d", len(orders))
		}
		assertOrder(t, orders, "A", 4000, 4.0)
	})

	t.Run("按分後に200円未満になる銘柄を除外", func(t *testing.T) {
		// A:99%, B:1%, amount:10000
		// A: 9900円, B: 100円→200円未満なので除外
		stockA := makeStock("A", 1000, true)
		stockB := makeStock("B", 155, true)
		stocks := []*stock.Stock{&stockA, &stockB}
		pf := makePortfolio(map[string]int{"A": 99, "B": 1})
		amount := mustAmount(10000)

		orders, err := calc.Calculate(amount, &pf, stocks)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(orders) != 1 {
			t.Fatalf("expected 1 order, got %d", len(orders))
		}
		assertOrder(t, orders, "A", 9900, 9.9)
	})

	t.Run("全銘柄がtradable=false", func(t *testing.T) {
		stockA := makeStock("A", 1000, false)
		stockB := makeStock("B", 155, false)
		stocks := []*stock.Stock{&stockA, &stockB}
		pf := makePortfolio(map[string]int{"A": 40, "B": 60})
		amount := mustAmount(10000)

		orders, err := calc.Calculate(amount, &pf, stocks)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(orders) != 0 {
			t.Fatalf("expected 0 orders, got %d", len(orders))
		}
	})

	t.Run("全銘柄が按分後200円未満", func(t *testing.T) {
		// amount:1000, A:50%, B:50%
		// A: 500円/株価10000→0.05株→数量0.05*10000=500円 → OK?
		// ここでは株価を高くして注文金額を200円未満にする
		// A:1%, B:1%, 残り98%は存在しない銘柄
		// → A:10000*1/100=100円、B:100円、両方200円未満
		stockA := makeStock("A", 1000, true)
		stockB := makeStock("B", 1000, true)
		stocks := []*stock.Stock{&stockA, &stockB}
		pf := makePortfolio(map[string]int{"A": 1, "B": 1})
		amount := mustAmount(10000)

		orders, err := calc.Calculate(amount, &pf, stocks)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(orders) != 0 {
			t.Fatalf("expected 0 orders, got %d", len(orders))
		}
	})

	t.Run("数量の小数点以下3桁切り捨て確認", func(t *testing.T) {
		// amount:10000, A:100%, 株価:3333
		// 注文金額:10000円, 数量=10000/3333=3.0003000...→3.0
		stockA := makeStock("A", 3333, true)
		stocks := []*stock.Stock{&stockA}
		pf := makePortfolio(map[string]int{"A": 100})
		amount := mustAmount(10000)

		orders, err := calc.Calculate(amount, &pf, stocks)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(orders) != 1 {
			t.Fatalf("expected 1 order, got %d", len(orders))
		}
		assertOrder(t, orders, "A", 10000, 3.0)
	})

	t.Run("ポートフォリオにない銘柄はstocksに存在しても無視", func(t *testing.T) {
		// portfolioにAのみ、stocksにA,B両方存在
		stockA := makeStock("A", 1000, true)
		stockB := makeStock("B", 155, true)
		stocks := []*stock.Stock{&stockA, &stockB}
		pf := makePortfolio(map[string]int{"A": 100})
		amount := mustAmount(10000)

		orders, err := calc.Calculate(amount, &pf, stocks)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(orders) != 1 {
			t.Fatalf("expected 1 order, got %d", len(orders))
		}
		assertOrder(t, orders, "A", 10000, 10.0)
	})
}

// アサーションヘルパー

func assertOrder(t *testing.T, orders []*trade.Order, symbol string, expectedAmount int, expectedQty float64) {
	t.Helper()
	for _, o := range orders {
		if o.Symbol.Value() == symbol {
			if o.Amount.Value() != expectedAmount {
				t.Errorf("order %s: expected amount %d, got %d", symbol, expectedAmount, o.Amount.Value())
			}
			if o.Quantity.Value() != expectedQty {
				t.Errorf("order %s: expected quantity %f, got %f", symbol, expectedQty, o.Quantity.Value())
			}
			return
		}
	}
	t.Errorf("order %s not found", symbol)
}