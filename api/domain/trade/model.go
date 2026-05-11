package trade

import "bloomo-exam-api/domain/shared/vo"

type Order struct {
    Symbol   vo.Ticker
    Amount   vo.Amount
    Quantity vo.Quantity
}

type Trade struct {
    Amount          vo.Amount
    TargetPortfolio map[vo.Ticker]vo.Ratio
    Orders          []Order
}