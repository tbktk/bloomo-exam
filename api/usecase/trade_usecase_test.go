package usecase_test

import (
	"bloomo-exam-api/domain/trade"
	"bloomo-exam-api/infrastructure/memory"
	"bloomo-exam-api/usecase"
	"testing"
)

// テスト用ヘルパー関数

func createTestUsecase() *usecase.TradeUsecase {
	stockRepo := memory.NewStockRepository()
	portfolioRepo := memory.NewPortfolioRepository()
	calculator := &trade.OrderCalculator{}
	return usecase.NewTradeUsecase(stockRepo, portfolioRepo, calculator)
}

func TestTradeUsecase_Execute(t *testing.T) {
	t.Run("successful trade with user_id 1 (A:40%, B:60%)", func(t *testing.T) {
		// Arrange
		uc := createTestUsecase()

		// Act
		output, err := uc.Execute(usecase.TradeInput{
			UserID: 1,
			Amount: 10000,
		})

		// Assert
		if err != nil {
			t.Fatalf("Execute() error = %v, want nil", err)
		}
		if output == nil {
			t.Fatalf("Execute() output is nil, want non-nil")
		}
		if output.Amount != 10000 {
			t.Errorf("Execute() amount = %d, want 10000", output.Amount)
		}
		// A:40%, B:60% で両方取引可能なので2注文
		if len(output.Orders) != 2 {
			t.Errorf("Execute() orders count = %d, want 2", len(output.Orders))
		}
		if output.TargetPortfolio["A"] != 40 || output.TargetPortfolio["B"] != 60 {
			t.Errorf("Execute() target portfolio mismatch")
		}
	})

	t.Run("error when amount is below minimum", func(t *testing.T) {
		// Arrange
		uc := createTestUsecase()

		// Act
		output, err := uc.Execute(usecase.TradeInput{
			UserID: 1,
			Amount: 999,
		})

		// Assert
		if err == nil {
			t.Errorf("Execute() error = nil, want error for invalid amount")
		}
		if output != nil {
			t.Errorf("Execute() output should be nil on error")
		}
	})

	t.Run("error when user_id is invalid (negative)", func(t *testing.T) {
		// Arrange
		uc := createTestUsecase()

		// Act
		output, err := uc.Execute(usecase.TradeInput{
			UserID: -1,
			Amount: 10000,
		})

		// Assert
		if err == nil {
			t.Errorf("Execute() error = nil, want error for invalid user_id")
		}
		if output != nil {
			t.Errorf("Execute() output should be nil on error")
		}
	})

	t.Run("error when portfolio not found (user_id 999)", func(t *testing.T) {
		// Arrange
		uc := createTestUsecase()

		// Act
		output, err := uc.Execute(usecase.TradeInput{
			UserID: 999,
			Amount: 10000,
		})

		// Assert
		if err == nil {
			t.Errorf("Execute() error = nil, want error for portfolio not found")
		}
		if output != nil {
			t.Errorf("Execute() output should be nil on error")
		}
	})

	t.Run("user_id 2 with E:100% (tradable=false)", func(t *testing.T) {
		// Arrange
		uc := createTestUsecase()

		// Act - ユーザー2はE銘柄のみで、E銘柄はtradable=false
		output, err := uc.Execute(usecase.TradeInput{
			UserID: 2,
			Amount: 10000,
		})

		// Assert
		if err != nil {
			t.Fatalf("Execute() error = %v, want nil", err)
		}
		// Eはtradable=falseなので注文なし
		if len(output.Orders) != 0 {
			t.Errorf("Execute() orders count = %d, want 0", len(output.Orders))
		}
	})

	t.Run("user_id 3 with mixed tradability (A:31%, B:40%, E:29%)", func(t *testing.T) {
		// Arrange
		uc := createTestUsecase()

		// Act - ユーザー3はA(tradable=true), B(tradable=true), E(tradable=false)
		output, err := uc.Execute(usecase.TradeInput{
			UserID: 3,
			Amount: 10000,
		})

		// Assert
		if err != nil {
			t.Fatalf("Execute() error = %v, want nil", err)
		}
		// AとBだけが取引可能
		if len(output.Orders) != 2 {
			t.Errorf("Execute() orders count = %d, want 2", len(output.Orders))
		}
	})

	t.Run("user_id 4 with multiple tradable stocks", func(t *testing.T) {
		// Arrange
		uc := createTestUsecase()

		// Act - ユーザー4はB:50%, C:49%, D:1%で全て取引可能
		output, err := uc.Execute(usecase.TradeInput{
			UserID: 4,
			Amount: 10000,
		})

		// Assert
		if err != nil {
			t.Fatalf("Execute() error = %v, want nil", err)
		}
		// B, C, Dで3注文、ただしDの1%は10000*1/100=100円で200円未満なので除外
		if len(output.Orders) != 2 {
			t.Errorf("Execute() orders count = %d, want 2 (D is excluded due to minimum)", len(output.Orders))
		}
	})

	t.Run("successful trade with large amount", func(t *testing.T) {
		// Arrange
		uc := createTestUsecase()

		// Act
		output, err := uc.Execute(usecase.TradeInput{
			UserID: 1,
			Amount: 100000,
		})

		// Assert
		if err != nil {
			t.Fatalf("Execute() error = %v, want nil", err)
		}
		if output.Amount != 100000 {
			t.Errorf("Execute() amount = %d, want 100000", output.Amount)
		}
		if len(output.Orders) != 2 {
			t.Errorf("Execute() orders count = %d, want 2", len(output.Orders))
		}
	})
}
