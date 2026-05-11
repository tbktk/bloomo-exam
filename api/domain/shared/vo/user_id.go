package vo

import "errors"

type UserID struct{ value int }

func NewUserID(v int) (UserID, error) {
    if v <= 0 {
        return UserID{}, errors.New("user_id must be positive")
    }
    return UserID{value: v}, nil
}

func (u UserID) Value() int { return u.value }

func (u UserID) Equals(other UserID) bool {
	return u.value == other.value
}