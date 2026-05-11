package main

import (
	"bloomo-exam-api/domain/trade"
	"bloomo-exam-api/infrastructure/memory"
	"bloomo-exam-api/interface/http/handler"
	"bloomo-exam-api/usecase"
	"log"
	"net/http"
)

func main() {
	stockRepo := memory.NewStockRepository()
	portfolioRepo := memory.NewPortfolioRepository()
	calculator := &trade.OrderCalculator{}

	tradeUsecase := usecase.NewTradeUsecase(stockRepo, portfolioRepo, calculator)
	tradeHandler := handler.NewTradeHandler(tradeUsecase)

	// ルーティング
	mux := http.NewServeMux()
	mux.HandleFunc("POST /users/{user_id}/trades", tradeHandler.Handle)

	// サーバー起動
	addr := ":1111"
	log.Printf("server starting on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}