package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log  *zap.Logger
	once sync.Once
)

// Init inicializa el logger global
func Init(level, format string) error {
	var err error
	once.Do(func() {
		var cfg zap.Config

		logLevel := mapLevel(level)

		if format == "json" {
			cfg = zap.NewProductionConfig()
		} else {
			cfg = zap.NewDevelopmentConfig()
		}

		cfg.Level = zap.NewAtomicLevelAt(logLevel)
		cfg.OutputPaths = []string{"stdout"}
		cfg.ErrorOutputPaths = []string{"stderr"}

		log, err = cfg.Build(
			zap.AddStacktrace(zapcore.ErrorLevel),
		)
	})

	return err
}

// GetLogger retorna el logger global
func GetLogger() *zap.Logger {
	if log == nil {
		log, _ = zap.NewProduction()
	}
	return log
}

// Info log de información
func Info(msg string, fields ...zap.Field) {
	GetLogger().Info(msg, fields...)
}

// Error log de error
func Error(msg string, fields ...zap.Field) {
	GetLogger().Error(msg, fields...)
}

// Warn log de advertencia
func Warn(msg string, fields ...zap.Field) {
	GetLogger().Warn(msg, fields...)
}

// Debug log de debug
func Debug(msg string, fields ...zap.Field) {
	GetLogger().Debug(msg, fields...)
}

// Fatal log fatal y salida
func Fatal(msg string, fields ...zap.Field) {
	GetLogger().Fatal(msg, fields...)
	os.Exit(1)
}

// WithRequest añade información de request al contexto
func WithRequest(method, path string, statusCode int, durationMs int64) []zap.Field {
	return []zap.Field{
		zap.String("method", method),
		zap.String("path", path),
		zap.Int("status", statusCode),
		zap.Int64("duration_ms", durationMs),
	}
}

// WithUser añade información de usuario al contexto
func WithUser(userID, email string) []zap.Field {
	return []zap.Field{
		zap.String("user_id", userID),
		zap.String("email", email),
	}
}

// WithError añade información de error al contexto
func WithError(err error) zap.Field {
	return zap.Error(err)
}

// mapLevel mapea string de nivel a zap level
func mapLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}
