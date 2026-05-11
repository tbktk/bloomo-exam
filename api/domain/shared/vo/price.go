package vo

import "errors"

type Price struct{ value int }

func NewPrice(v int) (Price, error) {
    if v <= 0 {
        return Price{}, errors.New("price must be positive")
    }
    return Price{value: v}, nil
}

func (p Price) Value() int { return p.value }