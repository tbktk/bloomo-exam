package stock

import "bloomo-exam-api/domain/shared/vo"

type Repository interface {
    FindAll() ([]*Stock, error)
    FindByTicker(ticker vo.Ticker) (*Stock, error)
}