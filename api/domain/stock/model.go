package stock

import (
	"errors"

	"bloomo-exam-api/domain/shared/vo"
)

var ErrNotFound = errors.New("stock not found")

type Stock struct {
    Ticker   vo.Ticker
    Price    vo.Price
    Tradable bool
}

func (s Stock) IsTradable() bool { return s.Tradable }