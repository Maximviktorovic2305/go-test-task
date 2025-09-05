package utils

import (
	"log/slog"
	"os"
)

// Logger это обертка вокруг slog.Logger
type Logger struct {
	*slog.Logger
}

// NewLogger создает новый экземпляр логгера
func NewLogger() *Logger {
	// Создать новый логгер с JSON обработчиком
	handler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(handler)

	// Установить логгер по умолчанию
	slog.SetDefault(logger)

	return &Logger{logger}
}

// Info записывает информационное сообщение
func (l *Logger) Info(msg string, args ...any) {
	l.Logger.Info(msg, args...)
}

// Error записывает сообщение об ошибке
func (l *Logger) Error(msg string, args ...any) {
	l.Logger.Error(msg, args...)
}

// Debug записывает отладочное сообщение
func (l *Logger) Debug(msg string, args ...any) {
	l.Logger.Debug(msg, args...)
}

// Warn записывает предупреждающее сообщение
func (l *Logger) Warn(msg string, args ...any) {
	l.Logger.Warn(msg, args...)
}
