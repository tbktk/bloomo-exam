package handler

import (
	"bloomo-exam-api/domain/portfolio"
	"bloomo-exam-api/infrastructure/logger"
	"bloomo-exam-api/usecase"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

type TradeHandler struct {
	tradeUsecase *usecase.TradeUsecase
}

func NewTradeHandler(tradeUsecase *usecase.TradeUsecase) *TradeHandler {
	return &TradeHandler{tradeUsecase: tradeUsecase}
}

type tradeRequest struct {
	Amount int `json:"amount"`
}

type orderResponse struct {
	Symbol   string  `json:"symbol"`
	Amount   int     `json:"amount"`
	Quantity float64 `json:"quantity"`
}

type tradeResponse struct {
	Amount          int               `json:"amount"`
	TargetPortfolio map[string]int    `json:"target_portfolio"`
	Orders          []orderResponse   `json:"orders"`
}

func (h *TradeHandler) Handle(w http.ResponseWriter, r *http.Request) {
	// HTTPメソッドのチェック
	if r.Method != http.MethodPost {
		respondWithError(w, r, http.StatusMethodNotAllowed, ErrCodeMethodNotAllowed,
			"only POST method is allowed")
		return
	}

	// user_idをパスパラメータから取得
	userIDStr := r.PathValue("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		respondWithError(w, r, http.StatusBadRequest, ErrCodeInvalidInput,
			"invalid user_id parameter", "user_id must be a positive integer")
		return
	}

	if userID <= 0 {
		respondWithError(w, r, http.StatusBadRequest, ErrCodeInvalidInput,
			"user_id must be positive", "user_id="+userIDStr+" is not allowed")
		return
	}

	// リクエストボディのデコード
	var req tradeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, r, http.StatusBadRequest, ErrCodeInvalidInput,
			"invalid request body", err.Error())
		logger.Warn("Failed to decode request body", map[string]interface{}{
			"user_id": userID,
			"error":   err.Error(),
		})
		return
	}

	// リクエストボディのバリデーション
	if req.Amount <= 0 {
		respondWithError(w, r, http.StatusBadRequest, ErrCodeInvalidInput,
			"invalid amount", "amount must be a positive integer")
		return
	}

	// ユースケース実行
	ctx := logger.NewContext().
		Add("user_id", userID).
		Add("amount", req.Amount)
	logger.LogContext(logger.DEBUG, "Executing trade usecase", ctx)
	
	output, err := h.tradeUsecase.Execute(usecase.TradeInput{
		UserID: userID,
		Amount: req.Amount,
	})
	if err != nil {
		if errors.Is(err, portfolio.ErrNotFound) {
			respondWithError(w, r, http.StatusNotFound, ErrCodeNotFound,
				"user portfolio not found", "user_id="+userIDStr+" does not have a portfolio")
			return
		}

		// その他のエラーを詳細に分類
		if strings.Contains(err.Error(), "invalid") {
			respondWithError(w, r, http.StatusBadRequest, ErrCodeInvalidInput,
				"invalid input parameters", err.Error())
		} else {
			respondWithError(w, r, http.StatusInternalServerError, ErrCodeInternalError,
				"failed to execute trade", err.Error())
			logger.Error("Trade execution failed", map[string]interface{}{
				"user_id": userID,
				"amount":  req.Amount,
				"error":   err.Error(),
			})
		}
		return
	}

	// レスポンス変換
	orders := make([]orderResponse, 0, len(output.Orders))
	for _, o := range output.Orders {
		orders = append(orders, orderResponse{
			Symbol:   o.Symbol,
			Amount:   o.Amount,
			Quantity: o.Quantity,
		})
	}

	response := tradeResponse{
		Amount:          output.Amount,
		TargetPortfolio: output.TargetPortfolio,
		Orders:          orders,
	}

	logger.LogContext(logger.INFO, "Trade executed successfully", logger.NewContext().
		Add("user_id", userID).
		Add("amount", output.Amount).
		Add("orders_count", len(orders)))
	
	respondJSON(w, http.StatusOK, response)
}