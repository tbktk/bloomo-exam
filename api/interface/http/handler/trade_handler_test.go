package handler_test

import (
	"bloomo-exam-api/domain/trade"
	"bloomo-exam-api/infrastructure/memory"
	"bloomo-exam-api/interface/http/handler"
	"bloomo-exam-api/usecase"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupTradeHandler() *handler.TradeHandler {
	stockRepo := memory.NewStockRepository()
	portfolioRepo := memory.NewPortfolioRepository()

	calculator := &trade.OrderCalculator{}
	uc := usecase.NewTradeUsecase(stockRepo, portfolioRepo, calculator)

	return handler.NewTradeHandler(uc)
}

func TestTradeHandler_Handle_Success(t *testing.T) {
	// Arrange
	h := setupTradeHandler()
	body := map[string]int{"amount": 10000}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/users/1/trades", bytes.NewBuffer(bodyBytes))
	req.SetPathValue("user_id", "1")
	w := httptest.NewRecorder()

	// Act
	h.Handle(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("Handle() status = %d, want %d", w.Code, http.StatusOK)
	}

	var response map[string]interface{}
	json.NewDecoder(w.Body).Decode(&response)

	if response["amount"] != float64(10000) {
		t.Errorf("Handle() amount = %v, want 10000", response["amount"])
	}
	if orders, ok := response["orders"].([]interface{}); ok {
		if len(orders) != 2 {
			t.Errorf("Handle() orders count = %d, want 2", len(orders))
		}
	}
}

func TestTradeHandler_Handle_MethodNotAllowed(t *testing.T) {
	// Arrange
	h := setupTradeHandler()
	req := httptest.NewRequest("GET", "/users/1/trades", nil)
	req.SetPathValue("user_id", "1")
	w := httptest.NewRecorder()

	// Act
	h.Handle(w, req)

	// Assert
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Handle() status = %d, want %d", w.Code, http.StatusMethodNotAllowed)
	}

	var errResponse map[string]interface{}
	json.NewDecoder(w.Body).Decode(&errResponse)
	if code, ok := errResponse["code"].(string); !ok || code != "METHOD_NOT_ALLOWED" {
		t.Errorf("Handle() error code = %v, want METHOD_NOT_ALLOWED", errResponse["code"])
	}
}

func TestTradeHandler_Handle_InvalidUserID(t *testing.T) {
	// Arrange
	h := setupTradeHandler()
	body := map[string]int{"amount": 10000}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/users/invalid/trades", bytes.NewBuffer(bodyBytes))
	req.SetPathValue("user_id", "invalid")
	w := httptest.NewRecorder()

	// Act
	h.Handle(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("Handle() status = %d, want %d", w.Code, http.StatusBadRequest)
	}

	var errResponse map[string]interface{}
	json.NewDecoder(w.Body).Decode(&errResponse)
	if code, ok := errResponse["code"].(string); !ok || code != "INVALID_INPUT" {
		t.Errorf("Handle() error code = %v, want INVALID_INPUT", errResponse["code"])
	}
}

func TestTradeHandler_Handle_InvalidRequestBody(t *testing.T) {
	// Arrange
	h := setupTradeHandler()
	req := httptest.NewRequest("POST", "/users/1/trades", bytes.NewBuffer([]byte("invalid json")))
	req.SetPathValue("user_id", "1")
	w := httptest.NewRecorder()

	// Act
	h.Handle(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("Handle() status = %d, want %d", w.Code, http.StatusBadRequest)
	}

	var errResponse map[string]interface{}
	json.NewDecoder(w.Body).Decode(&errResponse)
	if code, ok := errResponse["code"].(string); !ok || code != "INVALID_INPUT" {
		t.Errorf("Handle() error code = %v, want INVALID_INPUT", errResponse["code"])
	}
}

func TestTradeHandler_Handle_InvalidAmount(t *testing.T) {
	// Arrange
	h := setupTradeHandler()
	body := map[string]int{"amount": 500}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/users/1/trades", bytes.NewBuffer(bodyBytes))
	req.SetPathValue("user_id", "1")
	w := httptest.NewRecorder()

	// Act
	h.Handle(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("Handle() status = %d, want %d", w.Code, http.StatusBadRequest)
	}

	var errResponse map[string]interface{}
	json.NewDecoder(w.Body).Decode(&errResponse)
	if code, ok := errResponse["code"].(string); !ok || code != "INVALID_INPUT" {
		t.Errorf("Handle() error code = %v, want INVALID_INPUT", errResponse["code"])
	}
}

func TestTradeHandler_Handle_UserNotFound(t *testing.T) {
	// Arrange
	h := setupTradeHandler()
	body := map[string]int{"amount": 10000}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/users/999/trades", bytes.NewBuffer(bodyBytes))
	req.SetPathValue("user_id", "999")
	w := httptest.NewRecorder()

	// Act
	h.Handle(w, req)

	// Assert
	if w.Code != http.StatusNotFound {
		t.Errorf("Handle() status = %d, want %d", w.Code, http.StatusNotFound)
	}

	var errResponse map[string]interface{}
	json.NewDecoder(w.Body).Decode(&errResponse)
	if code, ok := errResponse["code"].(string); !ok || code != "NOT_FOUND" {
		t.Errorf("Handle() error code = %v, want NOT_FOUND", errResponse["code"])
	}
}

func TestTradeHandler_Handle_ResponseFormat(t *testing.T) {
	// Arrange
	h := setupTradeHandler()
	body := map[string]int{"amount": 10000}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/users/1/trades", bytes.NewBuffer(bodyBytes))
	req.SetPathValue("user_id", "1")
	w := httptest.NewRecorder()

	// Act
	h.Handle(w, req)

	// Assert
	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Handle() Content-Type = %s, want application/json",
			w.Header().Get("Content-Type"))
	}

	var response map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Handle() response is not valid JSON: %v", err)
	}

	// Check required fields in success response
	if _, ok := response["amount"]; !ok {
		t.Errorf("Handle() response missing 'amount' field")
	}
	if _, ok := response["target_portfolio"]; !ok {
		t.Errorf("Handle() response missing 'target_portfolio' field")
	}
	if _, ok := response["orders"]; !ok {
		t.Errorf("Handle() response missing 'orders' field")
	}
}

func TestTradeHandler_Handle_ErrorResponseFormat(t *testing.T) {
	// Arrange - invalid user_id to trigger error response
	h := setupTradeHandler()
	body := map[string]int{"amount": 10000}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/users/invalid/trades", bytes.NewBuffer(bodyBytes))
	req.SetPathValue("user_id", "invalid")
	w := httptest.NewRecorder()

	// Act
	h.Handle(w, req)

	// Assert
	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Handle() Content-Type = %s, want application/json",
			w.Header().Get("Content-Type"))
	}

	var errResponse map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&errResponse); err != nil {
		t.Fatalf("Handle() error response is not valid JSON: %v", err)
	}

	// Check required fields in error response
	if _, ok := errResponse["code"]; !ok {
		t.Errorf("Handle() error response missing 'code' field")
	}
	if _, ok := errResponse["message"]; !ok {
		t.Errorf("Handle() error response missing 'message' field")
	}
	if _, ok := errResponse["timestamp"]; !ok {
		t.Errorf("Handle() error response missing 'timestamp' field")
	}
	if _, ok := errResponse["path"]; !ok {
		t.Errorf("Handle() error response missing 'path' field")
	}
}
