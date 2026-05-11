package portfolio

import "bloomo-exam-api/domain/shared/vo"

type Repository interface {
    FindByUserID(userID vo.UserID) (*Portfolio, error)
}