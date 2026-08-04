package database

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	gormLogger "gorm.io/gorm/logger"
)

var _ gormLogger.Interface = (*GormLogger)(nil)

type GormLogger struct {
	logger *slog.Logger
}

func NewGormLogger() *GormLogger {
	return &GormLogger{
		logger: slog.Default().With("component", "database"),
	}
}

func (l *GormLogger) LogMode(_ gormLogger.LogLevel) gormLogger.Interface {
	return l
}

func (l *GormLogger) Info(ctx context.Context, msg string, data ...any) {
	l.logf(ctx, slog.LevelInfo, msg, data...)
}

func (l *GormLogger) Warn(ctx context.Context, msg string, data ...any) {
	l.logf(ctx, slog.LevelWarn, msg, data...)
}

func (l *GormLogger) Error(ctx context.Context, msg string, data ...any) {
	l.logf(ctx, slog.LevelError, msg, data...)
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	ctx = ensureContext(ctx)
	level := slog.LevelDebug
	if err != nil && !errors.Is(err, gormLogger.ErrRecordNotFound) {
		level = slog.LevelError
	}
	if !l.logger.Enabled(ctx, level) {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()
	if err == nil {
		l.logger.DebugContext(ctx, "Trace", "elapsed", elapsed, "rows", rows, "sql", sql)
		return
	}

	if level == slog.LevelError {
		l.logger.ErrorContext(ctx, "Trace", "error", err, "elapsed", elapsed, "rows", rows, "sql", sql)
		return
	}
	l.logger.DebugContext(ctx, "Trace", "error", err, "elapsed", elapsed, "rows", rows, "sql", sql)
}

func (l *GormLogger) logf(ctx context.Context, level slog.Level, msg string, data ...any) {
	ctx = ensureContext(ctx)
	if !l.logger.Enabled(ctx, level) {
		return
	}
	l.logger.Log(ctx, level, fmt.Sprintf(msg, data...))
}

func ensureContext(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}
	return ctx
}
