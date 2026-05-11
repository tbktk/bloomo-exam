package vo

import (
	"errors"
	"math"
)

type Quantity struct{ value float64 }

// 数量 * 株価 <= 注文金額を満たす最大値（小数点以下3桁）
func NewQuantity(orderAmount int, price int) (Quantity, error) {
    if price <= 0 {
        return Quantity{}, errors.New("price must be positive")
    }
    raw := float64(orderAmount) / float64(price)
    // 小数点以下3桁切り捨て
    truncated := math.Floor(raw*1000) / 1000
    return Quantity{value: truncated}, nil
}

func (q Quantity) Value() float64 { return q.value }