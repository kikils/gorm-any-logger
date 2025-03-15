# gorm-any-logger

`gorm-any-logger` is a logger for Gorm that supports any logger.

## Installation

```bash
go get github.com/kikils/gorm-any-logger
```

## Usage

```go
import (
	"github.com/kikils/gorm-any-logger"
	"gorm.io/gorm"
)

logger := gormanylogger.New(
	gormanylogger.WithLogLevel{LogLevel: gormanylogger.LogLevelInfo},
	gormanylogger.WithLogFunc{LogFunc: func(ctx context.Context, level gormanylogger.LogLevel, msg string, params *gormanylogger.QueryLogParams) {
		// Use your logger here
	}},
)

db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
	Logger: logger,
})
```

## Examples

### Using with logrus

```go
import (
	"github.com/kikils/gorm-any-logger"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

log := logrus.New()

logger := gormanylogger.New(
	gormanylogger.WithLogLevel{LogLevel: gormanylogger.LogLevelInfo},
	gormanylogger.WithLogFunc{LogFunc: func(ctx context.Context, level gormanylogger.LogLevel, msg string, params *gormanylogger.QueryLogParams) {
		logEntry := log.WithContext(ctx)
		
		if params != nil {
			logEntry = logEntry.WithFields(logrus.Fields{
				"elapsed":     params.Elapsed.String(),
				"affected":    params.Affected,
				"sql":         params.SQL,
				"is_slow":     params.IsSlowQuery,
			})

			if params.Err != nil {
				logEntry = logEntry.WithError(params.Err)
			}
		}
		
		switch level {
		case gormanylogger.LogLevelError:
			logEntry.Error(msg)
		case gormanylogger.LogLevelWarn:
			logEntry.Warn(msg)
		case gormanylogger.LogLevelInfo:
			logEntry.Info(msg)
		}
	}},
)

db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
	Logger: logger,
})
```

### Using with zap

```go
import (
	"github.com/kikils/gorm-any-logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

zapLogger, _ := zap.NewProduction()

logger := gormanylogger.New(
	gormanylogger.WithLogLevel{LogLevel: gormanylogger.LogLevelInfo},
	gormanylogger.WithLogFunc{LogFunc: func(ctx context.Context, level gormanylogger.LogLevel, msg string, params *gormanylogger.QueryLogParams) {
		var fields []zap.Field
		
		if params != nil {
			fields = []zap.Field{
				zap.String("elapsed", params.Elapsed.String()),
				zap.Int64("affected", params.Affected),
				zap.String("sql", params.SQL),
				zap.Bool("is_slow", params.IsSlowQuery),
			}

			if params.Err != nil {
				fields = append(fields, zap.Error(params.Err))
			}
		}
		
		switch level {
		case gormanylogger.LogLevelError:
			zapLogger.Error(msg, fields...)
		case gormanylogger.LogLevelWarn:
			zapLogger.Warn(msg, fields...)
		case gormanylogger.LogLevelInfo:
			zapLogger.Info(msg, fields...)
		}
	}},
)

db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
	Logger: logger,
})
```

### Using with Go 1.21 slog

```go
import (
	"github.com/kikils/gorm-any-logger"
	"log/slog"
	"gorm.io/gorm"
)

slogger := slog.Default()

logger := gormanylogger.New(
	gormanylogger.WithLogLevel{LogLevel: gormanylogger.LogLevelInfo},
	gormanylogger.WithLogFunc{LogFunc: func(ctx context.Context, level gormanylogger.LogLevel, msg string, params *gormanylogger.QueryLogParams) {
		var attrs []slog.Attr
		
		if params != nil {
			attrs = []slog.Attr{
				slog.String("elapsed", params.Elapsed.String()),
				slog.Int64("affected", params.Affected),
				slog.String("sql", params.SQL),
				slog.Bool("is_slow", params.IsSlowQuery),
			}

			if params.Err != nil {
				attrs = append(attrs, slog.Any("error", params.Err))
			}
		}
		
		switch level {
		case gormanylogger.LogLevelError:
			slogger.LogAttrs(ctx, slog.LevelError, msg, attrs...)
		case gormanylogger.LogLevelWarn:
			slogger.LogAttrs(ctx, slog.LevelWarn, msg, attrs...)
		case gormanylogger.LogLevelInfo:
			slogger.LogAttrs(ctx, slog.LevelInfo, msg, attrs...)
		}
	}},
)

db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
	Logger: logger,
})
```

## Options

### WithLogLevel

```go
gormanylogger.WithLogLevel{gormanylogger.LogLevelInfo}
```

### WithLogFunc

```go
gormanylogger.WithLogFunc{func(ctx context.Context, level gormanylogger.LogLevel, msg string, params *gormanylogger.QueryLogParams) {
	// Use your logger here
	// Check if params is nil before using it as it may be nil in some cases
}}
```

### WithIgnoreRecordNotFoundError

```go
gormanylogger.WithIgnoreRecordNotFoundError{true}
```

### WithSlowThreshold

```go
gormanylogger.WithSlowThreshold{100 * time.Millisecond}
```

## License

`gorm-any-logger` is released under the MIT License. See the [LICENSE](LICENSE) file for details.
