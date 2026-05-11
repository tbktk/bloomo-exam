package vo

import "testing"

func TestNewPrice(t *testing.T) {
	tests := []struct {
		name    string
		value   int
		wantErr bool
	}{
		{
			name:    "valid positive price",
			value:   1000,
			wantErr: false,
		},
		{
			name:    "valid small positive price",
			value:   1,
			wantErr: false,
		},
		{
			name:    "invalid zero price",
			value:   0,
			wantErr: true,
		},
		{
			name:    "invalid negative price",
			value:   -100,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPrice(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPrice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.Value() != tt.value {
				t.Errorf("NewPrice() got value %d, want %d", got.Value(), tt.value)
			}
		})
	}
}

func TestPrice_Value(t *testing.T) {
	t.Run("returns correct value", func(t *testing.T) {
		price, _ := NewPrice(2500)
		if price.Value() != 2500 {
			t.Errorf("Value() = %d, want 2500", price.Value())
		}
	})
}
