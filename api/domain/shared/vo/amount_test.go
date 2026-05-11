package vo

import "testing"

func TestNewAmount(t *testing.T) {
	tests := []struct {
		name    string
		value   int
		wantErr bool
	}{
		{
			name:    "valid amount at minimum",
			value:   1000,
			wantErr: false,
		},
		{
			name:    "valid amount above minimum",
			value:   10000,
			wantErr: false,
		},
		{
			name:    "invalid amount below minimum",
			value:   999,
			wantErr: true,
		},
		{
			name:    "invalid amount at zero",
			value:   0,
			wantErr: true,
		},
		{
			name:    "invalid amount negative",
			value:   -1000,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAmount(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAmount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.Value() != tt.value {
				t.Errorf("NewAmount() got value %d, want %d", got.Value(), tt.value)
			}
		})
	}
}

func TestNewOrderAmount(t *testing.T) {
	tests := []struct {
		name    string
		value   int
		wantErr bool
	}{
		{
			name:    "valid order amount at minimum",
			value:   200,
			wantErr: false,
		},
		{
			name:    "valid order amount above minimum",
			value:   1000,
			wantErr: false,
		},
		{
			name:    "invalid order amount below minimum",
			value:   199,
			wantErr: true,
		},
		{
			name:    "invalid order amount at zero",
			value:   0,
			wantErr: true,
		},
		{
			name:    "invalid order amount negative",
			value:   -200,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewOrderAmount(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewOrderAmount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.Value() != tt.value {
				t.Errorf("NewOrderAmount() got value %d, want %d", got.Value(), tt.value)
			}
		})
	}
}

func TestAmount_Value(t *testing.T) {
	t.Run("returns correct value", func(t *testing.T) {
		amount, _ := NewAmount(5000)
		if amount.Value() != 5000 {
			t.Errorf("Value() = %d, want 5000", amount.Value())
		}
	})
}
