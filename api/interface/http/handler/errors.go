package handler

import (
	"bloomo-exam-api/infrastructure/logger"
	"encoding/json"
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

	// ロギングコンテキストを作成
	ctx := logger.NewContext().
		Add("method", r.Method).
		Add("path", r.URL.Path).
		Add("status_code", status).
		Add("error_code", code).
		Add("details", details)

	// ログレベルを決定
	logLevel := logger.WARN
	if status >= 500 {
		logLevel = logger.ERROR
	}

	// ロギング
	logger.LogContext(logLevel, message, ctx)

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
				ctx := logger.NewContext().
					Add("method", r.Method).
					Add("path", r.URL.Path).
					Add("panic", err)
				logger.LogContext(logger.ERROR, "Panic recovered", ctx)

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

// responseWriterWrapper はhttp.ResponseWriterをラップしてステータスコードを記録する
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (w *responseWriterWrapper) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *responseWriterWrapper) Write(b []byte) (int, error) {
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}
	return w.ResponseWriter.Write(b)
}

// LoggingMiddleware はリクエストとレスポンスをログ出力する
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		ctx := logger.NewContext().
			Add("method", r.Method).
			Add("path", r.URL.Path)
		logger.LogContext(logger.DEBUG, "Request received", ctx)

		wrapped := &responseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)
		ctxResponse := logger.NewContext().
			Add("method", r.Method).
			Add("path", r.URL.Path).
			Add("status_code", wrapped.statusCode).
			Add("duration", duration.String())
		logger.LogContext(logger.DEBUG, "Request completed", ctxResponse)
	})
}
