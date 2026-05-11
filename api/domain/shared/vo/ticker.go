package vo

import "errors"

type Ticker struct{ value string }

func NewTicker(v string) (Ticker, error) {
    if v == "" {
        return Ticker{}, errors.New("ticker must not be empty")
    }
    return Ticker{value: v}, nil
}

func (t Ticker) Value() string { return t.value }
func (t Ticker) Equals(other Ticker) bool { return t.value == other.value }