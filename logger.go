package gormanylogger

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type LogLevel = gormlogger.LogLevel

const (
	// Silent silent log level
	LogLevelSilent = gormlogger.Silent
	// Error log level
	LogLevelError = gormlogger.Error
	// Warn log level
	LogLevelWarn = gormlogger.Warn
	// Info log level
	LogLevelInfo = gormlogger.Info
)

type QueryLogParams struct {
	SQL         string
	Affected    int64
	Elapsed     time.Duration
	Err         error
	IsSlowQuery bool
}

type Logger struct {
	IgnoreRecordNotFoundError bool
	SlowThreshold             time.Duration
	LogLevel                  LogLevel
	LogFunc                   func(ctx context.Context, level LogLevel, msg string, params *QueryLogParams)
}

func New(opts ...Option) *Logger {
	allOpts := []Option{
		WithSlowThreshold{200 * time.Millisecond},
		WithLogLevel{LogLevelInfo},
		WithLogFunc{func(ctx context.Context, level LogLevel, msg string, params *QueryLogParams) {
			switch level {
			case LogLevelSilent:
				return
			case LogLevelError:
				fmt.Printf("[ERROR] %s\n", msg)
			case LogLevelWarn:
				fmt.Printf("[WARN] %s\n", msg)
			case LogLevelInfo:
				fmt.Printf("[INFO] %s\n", msg)
			}
		}},
	}
	allOpts = append(allOpts, opts...)

	options := options{}
	for _, opt := range allOpts {
		opt.apply(&options)
	}

	return &Logger{
		IgnoreRecordNotFoundError: options.IgnoreRecordNotFoundError,
		SlowThreshold:             options.SlowThreshold,
		LogLevel:                  options.LogLevel,
		LogFunc:                   options.LogFunc,
	}
}

func (l *Logger) LogMode(level LogLevel) gormlogger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

func (l *Logger) Info(ctx context.Context, msg string, data ...any) {
	if l.LogLevel < LogLevelInfo {
		return
	}
	if len(data) <= 0 {
		l.LogFunc(ctx, LogLevelInfo, msg, nil)
		return
	}
	q, ok := data[0].(*QueryLogParams)
	if !ok {
		l.LogFunc(ctx, LogLevelInfo, fmt.Sprintf(msg, data...), nil)
		return
	}
	l.LogFunc(ctx, LogLevelInfo, msg, q)
}

func (l *Logger) Warn(ctx context.Context, msg string, data ...any) {
	if l.LogLevel < LogLevelWarn {
		return
	}
	if len(data) <= 0 {
		l.LogFunc(ctx, LogLevelWarn, msg, nil)
		return
	}
	q, ok := data[0].(*QueryLogParams)
	if !ok {
		l.LogFunc(ctx, LogLevelWarn, fmt.Sprintf(msg, data...), nil)
		return
	}
	l.LogFunc(ctx, LogLevelWarn, msg, q)
}

func (l *Logger) Error(ctx context.Context, msg string, data ...any) {
	if l.LogLevel < LogLevelError {
		return
	}
	if len(data) <= 0 {
		l.LogFunc(ctx, LogLevelError, msg, nil)
		return
	}
	q, ok := data[0].(*QueryLogParams)
	if !ok {
		l.LogFunc(ctx, LogLevelError, fmt.Sprintf(msg, data...), nil)
		return
	}
	l.LogFunc(ctx, LogLevelError, msg, q)
}

func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= LogLevelSilent {
		return
	}
	elapsed := time.Since(begin)
	sql, affected := fc()

	switch {
	case err != nil && l.LogLevel >= LogLevelError && (!l.IgnoreRecordNotFoundError || !errors.Is(err, gorm.ErrRecordNotFound)):
		l.Error(ctx, sql, &QueryLogParams{
			SQL:         sql,
			Affected:    affected,
			Elapsed:     elapsed,
			Err:         err,
			IsSlowQuery: false,
		})
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= LogLevelWarn:
		l.Warn(ctx, sql, &QueryLogParams{
			SQL:         sql,
			Affected:    affected,
			Elapsed:     elapsed,
			IsSlowQuery: true,
		})
	case l.LogLevel >= LogLevelInfo:
		l.Info(ctx, sql, &QueryLogParams{
			SQL:         sql,
			Affected:    affected,
			Elapsed:     elapsed,
			IsSlowQuery: false,
		})
	}
}
