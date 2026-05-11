package vo

import "errors"

const MinTradeAmount = 1000
const MinOrderAmount = 200

type Amount struct{ value int }

func NewAmount(v int) (Amount, error) {
    if v < MinTradeAmount {
        return Amount{}, errors.New("amount must be at least 1000")
    }
    return Amount{value: v}, nil
}

func NewOrderAmount(v int) (Amount, error) {
    if v < MinOrderAmount {
        return Amount{}, errors.New("order amount must be at least 200")
    }
    return Amount{value: v}, nil
}

func (a Amount) Value() int { return a.value }