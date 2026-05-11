package vo

import (
	"math"
	"testing"
)

func TestNewQuantity(t *testing.T) {
	tests := []struct {
		name           string
		orderAmount    int
		price          int
		expectedValue  float64
		wantErr        bool
		description    string
	}{
		{
			name:          "simple division",
			orderAmount:   4000,
			price:         1000,
			expectedValue: 4.0,
			wantErr:       false,
			description:   "4000 / 1000 = 4.0",
		},
		{
			name:          "truncate to 3 decimal places",
			orderAmount:   6000,
			price:         155,
			expectedValue: 38.709,
			wantErr:       false,
			description:   "6000 / 155 = 38.7096... → 38.709",
		},
		{
			name:          "very small quantity",
			orderAmount:   200,
			price:         1000,
			expectedValue: 0.2,
			wantErr:       false,
			description:   "200 / 1000 = 0.2",
		},
		{
			name:          "truncate to 3 decimal places (complex)",
			orderAmount:   1000,
			price:         333,
			expectedValue: 3.003,
			wantErr:       false,
			description:   "1000 / 333 = 3.003003... → 3.003",
		},
		{
			name:        "invalid zero price",
			orderAmount: 1000,
			price:       0,
			wantErr:     true,
			description: "price must be positive",
		},
		{
			name:        "invalid negative price",
			orderAmount: 1000,
			price:       -100,
			wantErr:     true,
			description: "price must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewQuantity(tt.orderAmount, tt.price)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewQuantity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// 浮動小数点数の比較（誤差を許容）
				if math.Abs(got.Value()-tt.expectedValue) > 1e-9 {
					t.Errorf("NewQuantity() got value %v, want %v (%s)",
						got.Value(), tt.expectedValue, tt.description)
				}
			}
		})
	}
}

func TestQuantity_Value(t *testing.T) {
	t.Run("returns correct value", func(t *testing.T) {
		quantity, _ := NewQuantity(4000, 1000)
		if quantity.Value() != 4.0 {
			t.Errorf("Value() = %v, want 4.0", quantity.Value())
		}
	})
}
