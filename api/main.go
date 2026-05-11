package main

import (
	"bloomo-exam-api/domain/trade"
	"bloomo-exam-api/infrastructure/logger"
	"bloomo-exam-api/infrastructure/repository/memory"
	"bloomo-exam-api/interface/http/handler"
	"bloomo-exam-api/usecase"
	"log"
	"net/http"
)

func main() {
	// ロガーの初期化
	logDir := "$HOME/.bloomo-exam/logs"
	if err := logger.Init(logDir, logger.INFO, false); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	
	logger.Info("Application starting")

	stockRepo := memory.NewStockRepository()
	portfolioRepo := memory.NewPortfolioRepository()
	calculator := &trade.OrderCalculator{}

	tradeUsecase := usecase.NewTradeUsecase(stockRepo, portfolioRepo, calculator)
	tradeHandler := handler.NewTradeHandler(tradeUsecase)

	// ルーティング
	mux := http.NewServeMux()
	mux.HandleFunc("POST /users/{user_id}/trades", tradeHandler.Handle)

	// ミドルウェアを適用
	var handlerWithMiddleware http.Handler = mux
	handlerWithMiddleware = handler.RecoveryMiddleware(handlerWithMiddleware)
	handlerWithMiddleware = handler.LoggingMiddleware(handlerWithMiddleware)

	// サーバー起動
	addr := ":1111"
	logger.Info("Server starting", map[string]interface{}{"addr": addr})
	if err := http.ListenAndServe(addr, handlerWithMiddleware); err != nil {
		logger.Fatal("Server failed", map[string]interface{}{"error": err.Error()})
	}
}