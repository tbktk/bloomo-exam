package vo

import "errors"

// 按分割合（0より大きく100以下の整数）
type Ratio struct{ value int }

func NewRatio(v int) (Ratio, error) {
    if v <= 0 || v > 100 {
        return Ratio{}, errors.New("ratio must be between 1 and 100")
    }
    return Ratio{value: v}, nil
}

func (r Ratio) Value() int { return r.value }