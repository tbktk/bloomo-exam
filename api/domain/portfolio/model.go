package portfolio

import (
	"errors"

	"bloomo-exam-api/domain/shared/vo"
)

var ErrNotFound = errors.New("portfolio not found")

type Portfolio struct {
    UserID          vo.UserID
    TargetPortfolio map[vo.Ticker]vo.Ratio
}