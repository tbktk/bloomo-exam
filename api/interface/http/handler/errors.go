package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// ErrorCode はエラーの種別を示すコード
type ErrorCode string

const (
	ErrCodeInvalidInput    ErrorCode = "INVALID_INPUT"
	ErrCodeNotFound        ErrorCode = "NOT_FOUND"
	ErrCodeInternalError   ErrorCode = "INTERNAL_ERROR"
	ErrCodeMethodNotAllowed ErrorCode = "METHOD_NOT_ALLOWED"
)

// ErrorResponse は統一されたエラーレスポンス形式
type ErrorResponse struct {
	Code      ErrorCode `json:"code"`
	Message   string    `json:"message"`
	Timestamp string    `json:"timestamp"`
	Path      string    `json:"path,omitempty"`
	Details   []string  `json:"details,omitempty"`
}

// SuccessResponse は統一されたサクセスレスポンス形式
type SuccessResponse struct {
	Data      interface{} `json:"data"`
	Timestamp string      `json:"timestamp"`
}

// respondWithError はエラーレスポンスを返す
func respondWithError(w http.ResponseWriter, r *http.Request, status int, code ErrorCode, message string, details ...string) {
	errorResp := ErrorResponse{
		Code:      code,
		Message:   message,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Path:      r.URL.Path,
		Details:   details,
	}

	// ログ出力
	logLevel := "WARN"
	if status >= 500 {
		logLevel = "ERROR"
	}
	log.Printf("[%s] %s - status=%d code=%s message=%s path=%s",
		logLevel, time.Now().Format(time.RFC3339), status, code, message, r.URL.Path)
	if len(details) > 0 {
		log.Printf("[DETAIL] %v", details)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errorResp)
}

// respondJSON はJSONレスポンスを返す
func respondJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// RecoveryMiddleware はパニックをキャッチしてログし、エラーレスポンスを返す
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("[PANIC] %s %s - %v", r.Method, r.URL.Path, err)
				
				errorResp := ErrorResponse{
					Code:      ErrCodeInternalError,
					Message:   "internal server error",
					Timestamp: time.Now().UTC().Format(time.RFC3339),
					Path:      r.URL.Path,
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(errorResp)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// LoggingMiddleware はリクエストとレスポンスをログ出力する
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("[INFO] %s %s started", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)

		duration := time.Since(start)
		log.Printf("[INFO] %s %s completed in %v", r.Method, r.URL.Path, duration)
	})
}
