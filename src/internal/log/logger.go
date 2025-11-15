package log

import (
	"context"

	"go.uber.org/zap"
)

type ctxLoggerKey struct{}

func Initialize(ctx context.Context) context.Context {
	logger, err := zap.NewProduction(
		zap.AddCallerSkip(1),
	)

	if err != nil {
		panic(err)
	}

	return loggerToContext(ctx, logger)
}

func loggerToContext(ctx context.Context, logger *zap.Logger) context.Context {
	ctx = context.WithValue(ctx, ctxLoggerKey{}, logger)
	return ctx
}

func loggerFromContext(ctx context.Context) *zap.Logger {
	return ctx.Value(ctxLoggerKey{}).(*zap.Logger)
}

func WithField(ctx context.Context, field zap.Field) context.Context {
	logger := loggerFromContext(ctx)
	newLogger := logger.With(field)
	return loggerToContext(ctx, newLogger)
}

func WithFields(ctx context.Context, fields ...zap.Field) context.Context {
	logger := loggerFromContext(ctx)
	newLogger := logger.With(fields...)
	return loggerToContext(ctx, newLogger)
}

func Info(ctx context.Context, msg string, fields ...zap.Field) {
	logger := loggerFromContext(ctx)
	logger.Info(msg, fields...)
}

func Error(ctx context.Context, msg string, fields ...zap.Field) {
	logger := loggerFromContext(ctx)
	logger.Error(msg, fields...)
}

func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	logger := loggerFromContext(ctx)
	logger.Warn(msg, fields...)
}

func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	logger := loggerFromContext(ctx)
	logger.Debug(msg, fields...)
}

func Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	logger := loggerFromContext(ctx)
	logger.Fatal(msg, fields...)
}

func Panic(ctx context.Context, msg string, fields ...zap.Field) {
	logger := loggerFromContext(ctx)
	logger.Panic(msg, fields...)
}
