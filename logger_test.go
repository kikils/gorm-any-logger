package gormanylogger_test

import (
	"context"
	"errors"
	"testing"
	"time"

	gormanylogger "github.com/kikils/gorm-any-logger"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestNew(t *testing.T) {
	t.Run("should return a new logger with default options", func(t *testing.T) {
		logger := gormanylogger.New()
		require.NotNil(t, logger)
		require.Equal(t, gormanylogger.LogLevelInfo, logger.LogLevel)
		require.Equal(t, 200*time.Millisecond, logger.SlowThreshold)
		require.Equal(t, gormanylogger.LogLevelInfo, logger.LogLevel)
	})
}

func TestLogger_LogMode(t *testing.T) {
	t.Run("should return a new logger with the given log level", func(t *testing.T) {
		logger := gormanylogger.New()
		newLogger := logger.LogMode(gormanylogger.LogLevelWarn)
		require.NotNil(t, newLogger)
		require.Equal(t, gormanylogger.LogLevelWarn, newLogger.(*gormanylogger.Logger).LogLevel)
	})
}

func TestLogger_Info(t *testing.T) {
	t.Run("under log level", func(t *testing.T) {
		logger := gormanylogger.New(
			gormanylogger.WithLogLevel{gormanylogger.LogLevelSilent},
			gormanylogger.WithLogFunc{func(ctx context.Context, level gormanylogger.LogLevel, msg string, params *gormanylogger.QueryLogParams) {
				t.Errorf("should not be called")
			}},
		)
		logger.Info(t.Context(), "test")
	})

	t.Run("without data", func(t *testing.T) {
		logger := gormanylogger.New(
			gormanylogger.WithLogLevel{gormanylogger.LogLevelInfo},
			gormanylogger.WithLogFunc{func(ctx context.Context, level gormanylogger.LogLevel, msg string, params *gormanylogger.QueryLogParams) {
				require.Equal(t, t.Context(), ctx)
				require.Equal(t, gormanylogger.LogLevelInfo, level)
				require.Equal(t, "test", msg)
				require.Nil(t, params)
			}},
		)
		logger.Info(t.Context(), "test")
	})

	t.Run("with data", func(t *testing.T) {
		logger := gormanylogger.New(
			gormanylogger.WithLogLevel{gormanylogger.LogLevelInfo},
			gormanylogger.WithLogFunc{func(ctx context.Context, level gormanylogger.LogLevel, msg string, params *gormanylogger.QueryLogParams) {
				require.Equal(t, t.Context(), ctx)
				require.Equal(t, gormanylogger.LogLevelInfo, level)
				require.Equal(t, "test: data", msg)
			}},
		)
		logger.Info(t.Context(), "test: %s", "data")
	})
}

func TestLogger_Warn(t *testing.T) {
	t.Run("under log level", func(t *testing.T) {
		logger := gormanylogger.New(
			gormanylogger.WithLogLevel{gormanylogger.LogLevelSilent},
			gormanylogger.WithLogFunc{func(ctx context.Context, level gormanylogger.LogLevel, msg string, params *gormanylogger.QueryLogParams) {
				t.Errorf("should not be called")
			}},
		)
		logger.Warn(t.Context(), "test")
	})

	t.Run("without data", func(t *testing.T) {
		logger := gormanylogger.New(
			gormanylogger.WithLogLevel{gormanylogger.LogLevelWarn},
			gormanylogger.WithLogFunc{func(ctx context.Context, level gormanylogger.LogLevel, msg string, params *gormanylogger.QueryLogParams) {
				require.Equal(t, t.Context(), ctx)
				require.Equal(t, gormanylogger.LogLevelWarn, level)
				require.Equal(t, "test", msg)
				require.Nil(t, params)
			}},
		)
		logger.Warn(t.Context(), "test")
	})

	t.Run("with data", func(t *testing.T) {
		logger := gormanylogger.New(
			gormanylogger.WithLogLevel{gormanylogger.LogLevelWarn},
			gormanylogger.WithLogFunc{func(ctx context.Context, level gormanylogger.LogLevel, msg string, params *gormanylogger.QueryLogParams) {
				require.Equal(t, t.Context(), ctx)
				require.Equal(t, gormanylogger.LogLevelWarn, level)
				require.Equal(t, "test: data", msg)
			}},
		)
		logger.Warn(t.Context(), "test: %s", "data")
	})
}

func TestLogger_Error(t *testing.T) {
	t.Run("under log level", func(t *testing.T) {
		logger := gormanylogger.New(
			gormanylogger.WithLogLevel{gormanylogger.LogLevelSilent},
			gormanylogger.WithLogFunc{func(ctx context.Context, level gormanylogger.LogLevel, msg string, params *gormanylogger.QueryLogParams) {
				t.Errorf("should not be called")
			}},
		)
		logger.Error(t.Context(), "test")
	})

	t.Run("without data", func(t *testing.T) {
		logger := gormanylogger.New(
			gormanylogger.WithLogLevel{gormanylogger.LogLevelError},
			gormanylogger.WithLogFunc{func(ctx context.Context, level gormanylogger.LogLevel, msg string, params *gormanylogger.QueryLogParams) {
				require.Equal(t, t.Context(), ctx)
				require.Equal(t, gormanylogger.LogLevelError, level)
				require.Equal(t, "test", msg)
				require.Nil(t, params)
			}},
		)

		logger.Error(t.Context(), "test")
	})

	t.Run("with data", func(t *testing.T) {
		logger := gormanylogger.New(
			gormanylogger.WithLogLevel{gormanylogger.LogLevelError},
			gormanylogger.WithLogFunc{func(ctx context.Context, level gormanylogger.LogLevel, msg string, params *gormanylogger.QueryLogParams) {
				require.Equal(t, t.Context(), ctx)
				require.Equal(t, gormanylogger.LogLevelError, level)
				require.Equal(t, "test: data", msg)
			}},
		)

		logger.Error(t.Context(), "test: %s", "data")
	})
}

func TestLogger_Trace(t *testing.T) {
	type args struct {
		ctx   context.Context
		begin time.Time
		fc    func() (string, int64)
		err   error
	}

	successQueryArgs := args{
		ctx:   t.Context(),
		begin: time.Now().Add(-1 * time.Second),
		fc: func() (string, int64) {
			return "SELECT * FROM users", 1
		},
		err: nil,
	}

	failureQueryArgs := args{
		ctx:   t.Context(),
		begin: time.Now(),
		fc: func() (string, int64) {
			return "SELECT * FROM users", 1
		},
		err: errors.New("error"),
	}

	notFoundQueryArgs := args{
		ctx:   t.Context(),
		begin: time.Now(),
		fc: func() (string, int64) {
			return "SELECT * FROM users", 0
		},
		err: gorm.ErrRecordNotFound,
	}

	tests := []struct {
		name string
		args args
		opts []gormanylogger.Option
	}{
		{
			name: "should not log anything if log level is silent",
			args: successQueryArgs,
			opts: []gormanylogger.Option{
				gormanylogger.WithLogLevel{gormanylogger.LogLevelSilent},
				gormanylogger.WithLogFunc{func(ctx context.Context, level gormanylogger.LogLevel, msg string, params *gormanylogger.QueryLogParams) {
					t.Errorf("should not be called")
				}},
			},
		},
		{
			name: "should log error if query failed",
			args: failureQueryArgs,
			opts: []gormanylogger.Option{
				gormanylogger.WithLogLevel{gormanylogger.LogLevelError},
				gormanylogger.WithLogFunc{func(ctx context.Context, level gormanylogger.LogLevel, msg string, params *gormanylogger.QueryLogParams) {
					require.Equal(t, t.Context(), ctx)
					require.Equal(t, gormanylogger.LogLevelError, level)
					require.Equal(t, "SELECT * FROM users", msg)
					require.Equal(t, int64(1), params.Affected)
					require.Equal(t, "error", params.Err.Error())
				}},
			},
		},
		{
			name: "should log warning if query is slow",
			args: successQueryArgs,
			opts: []gormanylogger.Option{
				gormanylogger.WithLogLevel{gormanylogger.LogLevelWarn},
				gormanylogger.WithSlowThreshold{100 * time.Millisecond},
				gormanylogger.WithLogFunc{func(ctx context.Context, level gormanylogger.LogLevel, msg string, params *gormanylogger.QueryLogParams) {
					require.Equal(t, t.Context(), ctx)
					require.Equal(t, gormanylogger.LogLevelWarn, level)
					require.Equal(t, "SELECT * FROM users", msg)
					require.Equal(t, int64(1), params.Affected)
					require.True(t, params.IsSlowQuery)
				}},
			},
		},
		{
			name: "should log info if query is successful",
			args: successQueryArgs,
			opts: []gormanylogger.Option{
				gormanylogger.WithLogLevel{gormanylogger.LogLevelInfo},
				gormanylogger.WithSlowThreshold{2 * time.Second},
				gormanylogger.WithLogFunc{func(ctx context.Context, level gormanylogger.LogLevel, msg string, params *gormanylogger.QueryLogParams) {
					require.Equal(t, t.Context(), ctx)
					require.Equal(t, gormanylogger.LogLevelInfo, level)
					require.Equal(t, "SELECT * FROM users", msg)
					require.Equal(t, int64(1), params.Affected)
					require.False(t, params.IsSlowQuery)
				}},
			},
		},
		{
			name: "should log info if query is not found",
			args: notFoundQueryArgs,
			opts: []gormanylogger.Option{
				gormanylogger.WithLogLevel{gormanylogger.LogLevelInfo},
				gormanylogger.WithLogFunc{func(ctx context.Context, level gormanylogger.LogLevel, msg string, params *gormanylogger.QueryLogParams) {
					require.Equal(t, t.Context(), ctx)
					require.Equal(t, gormanylogger.LogLevelError, level)
					require.Equal(t, "SELECT * FROM users", msg)
					require.Equal(t, int64(0), params.Affected)
					require.Equal(t, gorm.ErrRecordNotFound.Error(), params.Err.Error())
				}},
			},
		},
		{
			name: "should log waning if query is slow and log level is info",
			args: successQueryArgs,
			opts: []gormanylogger.Option{
				gormanylogger.WithLogLevel{gormanylogger.LogLevelInfo},
				gormanylogger.WithSlowThreshold{100 * time.Millisecond},
				gormanylogger.WithLogFunc{func(ctx context.Context, level gormanylogger.LogLevel, msg string, params *gormanylogger.QueryLogParams) {
					require.Equal(t, t.Context(), ctx)
					require.Equal(t, gormanylogger.LogLevelWarn, level)
					require.Equal(t, "SELECT * FROM users", msg)
					require.Equal(t, int64(1), params.Affected)
					require.True(t, params.IsSlowQuery)
				}},
			},
		},
		{
			name: "ignore not found error if log level is error",
			args: notFoundQueryArgs,
			opts: []gormanylogger.Option{
				gormanylogger.WithLogLevel{gormanylogger.LogLevelError},
				gormanylogger.WithIgnoreRecordNotFoundError{true},
				gormanylogger.WithLogFunc{func(ctx context.Context, level gormanylogger.LogLevel, msg string, params *gormanylogger.QueryLogParams) {
					t.Errorf("should not be called")
				}},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			logger := gormanylogger.New(test.opts...)
			logger.Trace(test.args.ctx, test.args.begin, test.args.fc, test.args.err)
		})
	}
}
