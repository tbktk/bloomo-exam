package handler

import (
	"bloomo-exam-api/domain/portfolio"
	"bloomo-exam-api/usecase"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
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
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// user_idをパスパラメータから取得
	userIDStr := r.PathValue("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid user_id")
		return
	}

	// リクエストボディのデコード
	var req tradeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// ユースケース実行
	output, err := h.tradeUsecase.Execute(usecase.TradeInput{
		UserID: userID,
		Amount: req.Amount,
	})
	if err != nil {
		if errors.Is(err, portfolio.ErrNotFound) {
			respondError(w, http.StatusNotFound, "portfolio not found")
			return
		}
		respondError(w, http.StatusBadRequest, err.Error())
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

	respondJSON(w, http.StatusOK, tradeResponse{
		Amount:          output.Amount,
		TargetPortfolio: output.TargetPortfolio,
		Orders:          orders,
	})
}

func respondJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}