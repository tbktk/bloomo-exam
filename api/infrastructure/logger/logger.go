package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

// LogLevel はログレベルを表す
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

var levelNames = []string{
	"DEBUG",
	"INFO",
	"WARN",
	"ERROR",
	"FATAL",
}

// LogEntry はログエントリの構造体
type LogEntry struct {
	Timestamp string      `json:"timestamp"`
	Level     string      `json:"level"`
	Message   string      `json:"message"`
	Context   interface{} `json:"context,omitempty"`
	Caller    string      `json:"caller,omitempty"`
	Duration  string      `json:"duration,omitempty"`
}

// Logger はアプリケーション全体で使用するロガー
type Logger struct {
	level      LogLevel
	writers    []io.Writer
	mu         sync.Mutex
	jsonFormat bool
}

var (
	defaultLogger *Logger
	once          sync.Once
)

// Init はロガーを初期化する
func Init(logDir string, level LogLevel, jsonFormat bool) error {
	var err error
	once.Do(func() {
		writers := []io.Writer{os.Stdout}

		// ログファイルへの出力を設定
		if logDir != "" {
			logDir = os.ExpandEnv(logDir)
			if err := os.MkdirAll(logDir, 0755); err != nil {
				log.Printf("[WARN] Failed to create log directory: %v", err)
			} else {
				logFile := filepath.Join(logDir, fmt.Sprintf("app-%s.log", time.Now().Format("2006-01-02")))
				file, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
				if err != nil {
					log.Printf("[WARN] Failed to open log file: %v", err)
				} else {
					writers = append(writers, file)
				}
			}
		}

		defaultLogger = &Logger{
			level:      level,
			writers:    writers,
			jsonFormat: jsonFormat,
		}
	})

	return err
}

// GetLogger はデフォルトロガーを返す
func GetLogger() *Logger {
	if defaultLogger == nil {
		Init("", INFO, false)
	}
	return defaultLogger
}

// Debug はDEBUGレベルのログを出力する
func Debug(message string, ctx ...interface{}) {
	GetLogger().log(DEBUG, message, ctx...)
}

// Info はINFOレベルのログを出力する
func Info(message string, ctx ...interface{}) {
	GetLogger().log(INFO, message, ctx...)
}

// Warn はWARNレベルのログを出力する
func Warn(message string, ctx ...interface{}) {
	GetLogger().log(WARN, message, ctx...)
}

// Error はERRORレベルのログを出力する
func Error(message string, ctx ...interface{}) {
	GetLogger().log(ERROR, message, ctx...)
}

// Fatal はFATALレベルのログを出力して終了する
func Fatal(message string, ctx ...interface{}) {
	GetLogger().log(FATAL, message, ctx...)
	os.Exit(1)
}

// WithDuration は処理時間付きでログを出力する
func WithDuration(level LogLevel, message string, duration time.Duration, ctx ...interface{}) {
	GetLogger().logWithDuration(level, message, duration, ctx...)
}

// log はログを出力する（内部用）
func (l *Logger) log(level LogLevel, message string, ctx ...interface{}) {
	if level < l.level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	var context interface{}
	if len(ctx) > 0 {
		if len(ctx) == 1 {
			context = ctx[0]
		} else {
			context = ctx
		}
	}

	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     levelNames[level],
		Message:   message,
		Context:   context,
		Caller:    getCaller(3),
	}

	var output string
	if l.jsonFormat {
		data, _ := json.Marshal(entry)
		output = string(data)
	} else {
		output = formatPlainText(entry)
	}

	for _, w := range l.writers {
		fmt.Fprintln(w, output)
	}
}

// logWithDuration は処理時間付きでログを出力する（内部用）
func (l *Logger) logWithDuration(level LogLevel, message string, duration time.Duration, ctx ...interface{}) {
	if level < l.level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	var context interface{}
	if len(ctx) > 0 {
		if len(ctx) == 1 {
			context = ctx[0]
		} else {
			context = ctx
		}
	}

	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     levelNames[level],
		Message:   message,
		Context:   context,
		Caller:    getCaller(3),
		Duration:  duration.String(),
	}

	var output string
	if l.jsonFormat {
		data, _ := json.Marshal(entry)
		output = string(data)
	} else {
		output = formatPlainText(entry)
	}

	for _, w := range l.writers {
		fmt.Fprintln(w, output)
	}
}

// formatPlainText はプレーンテキスト形式でログエントリをフォーマットする
func formatPlainText(entry LogEntry) string {
	sb := strings.Builder{}
	sb.WriteString("[")
	sb.WriteString(entry.Level)
	sb.WriteString("] ")
	sb.WriteString(entry.Timestamp)
	sb.WriteString(" ")
	sb.WriteString(entry.Message)

	if entry.Caller != "" {
		sb.WriteString(" (")
		sb.WriteString(entry.Caller)
		sb.WriteString(")")
	}

	if entry.Duration != "" {
		sb.WriteString(" duration=")
		sb.WriteString(entry.Duration)
	}

	if entry.Context != nil {
		sb.WriteString(" context=")
		if contextBytes, err := json.Marshal(entry.Context); err == nil {
			sb.Write(contextBytes)
		}
	}

	return sb.String()
}

// getCaller は呼び出し元のファイル:行番号を取得する
func getCaller(skip int) string {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return ""
	}

	funcName := runtime.FuncForPC(pc).Name()
	funcName = funcName[strings.LastIndex(funcName, ".")+1:]
	file = file[strings.LastIndex(file, "/")+1:]

	return fmt.Sprintf("%s:%d:%s()", file, line, funcName)
}

// Context は構造化ログ用のコンテキスト情報
type Context struct {
	RequestID  string                 `json:"request_id,omitempty"`
	UserID     int                    `json:"user_id,omitempty"`
	Method     string                 `json:"method,omitempty"`
	Path       string                 `json:"path,omitempty"`
	StatusCode int                    `json:"status_code,omitempty"`
	Duration   string                 `json:"duration,omitempty"`
	Extra      map[string]interface{} `json:"extra,omitempty"`
}

// NewContext は新しいコンテキストを作成する
func NewContext() *Context {
	return &Context{
		Extra: make(map[string]interface{}),
	}
}

// Add はコンテキストに情報を追加する
func (c *Context) Add(key string, value interface{}) *Context {
	c.Extra[key] = value
	return c
}

// LogContext はコンテキスト情報をログに出力する
func LogContext(level LogLevel, message string, ctx *Context) {
	logger := GetLogger()
	if level >= logger.level {
		logger.log(level, message, ctx)
	}
}
