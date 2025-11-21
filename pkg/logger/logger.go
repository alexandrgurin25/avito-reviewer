package logger

import (
	"context"

	"go.uber.org/zap"
)

const (
	RequestId = "request_id"
)

// Обертка для работы с контекстом
type Logger struct {
	zap *zap.Logger
}

type ctxKey struct{}

func New(ctx context.Context) (context.Context, error) {
	if GetLoggerFromCtx(ctx) != nil {
		return ctx, nil
	}

	logger, err := zap.NewDevelopment(zap.AddCaller())
	if err != nil {
		return nil, err
	}

	ctx = context.WithValue(ctx, ctxKey{}, &Logger{logger})
	return ctx, nil
}

func GetLoggerFromCtx(ctx context.Context) *Logger {
	if logger := loggerFromContext(ctx); logger != nil {
		return logger
	}

	zapLogger, _ := zap.NewDevelopment(zap.AddCaller())
	return &Logger{zap: zapLogger}
}

func loggerFromContext(ctx context.Context) *Logger {
	if ctx == nil {
		return nil
	}

	if logger, ok := ctx.Value(ctxKey{}).(*Logger); ok && logger != nil {
		return logger
	}
	return nil
}

func (l *Logger) addContextFields(ctx context.Context, fields []zap.Field) []zap.Field {
	if ctx == nil {
		return fields
	}

	if requestId, ok := ctx.Value(RequestId).(string); ok && requestId != "" {
		fields = append(fields, zap.String("request_id", requestId))
	}

	return fields
}

func (l *Logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	fields = l.addContextFields(ctx, fields)
	l.zap.Info(msg, fields...)
}

func (l *Logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	fields = l.addContextFields(ctx, fields)
	l.zap.Error(msg, fields...)
}

func (l *Logger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	fields = l.addContextFields(ctx, fields)
	l.zap.Fatal(msg, fields...)
}

func (l *Logger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	fields = l.addContextFields(ctx, fields)
	l.zap.Debug(msg, fields...)
}

func (l *Logger) With(fields ...zap.Field) *Logger {
	return &Logger{zap: l.zap.With(fields...)}
}
