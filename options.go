package gormanylogger

import (
	"context"
	"time"
)

type options struct {
	IgnoreRecordNotFoundError bool
	SlowThreshold             time.Duration
	LogLevel                  LogLevel
	LogFunc                   func(ctx context.Context, level LogLevel, msg string, params *QueryLogParams)
}

type Option interface{ apply(*options) }

type WithIgnoreRecordNotFoundError struct {
	IgnoreRecordNotFoundError bool
}

func (o WithIgnoreRecordNotFoundError) apply(opts *options) {
	opts.IgnoreRecordNotFoundError = o.IgnoreRecordNotFoundError
}

type WithSlowThreshold struct {
	SlowThreshold time.Duration
}

func (o WithSlowThreshold) apply(opts *options) {
	opts.SlowThreshold = o.SlowThreshold
}

type WithLogLevel struct {
	LogLevel LogLevel
}

func (o WithLogLevel) apply(opts *options) {
	opts.LogLevel = o.LogLevel
}

type WithLogFunc struct {
	LogFunc func(ctx context.Context, level LogLevel, msg string, params *QueryLogParams)
}

func (o WithLogFunc) apply(opts *options) {
	opts.LogFunc = o.LogFunc
}
