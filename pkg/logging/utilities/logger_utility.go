package utilities

import (
	"errors"
	"github.com/ereb-or-od/kenobi/pkg/logging/enumeration"
	"github.com/xiam/to"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

func ToZapFields(parameters ...map[string]interface{}) *[]zap.Field {
	var zapFields []zap.Field
	if parameters != nil && len(parameters) > 0 {
		for _, parameterItem := range parameters {
			for key, value := range parameterItem {
				if field, err := toZapField(key, value); err == nil {
					if field != nil {
						zapFields = append(zapFields, *field)
					}
				}

			}
		}
		return &zapFields
	}
	return &zapFields
}

func toZapField(key string, value interface{}) (*zap.Field, error) {
	if len(key) == 0 {
		return nil, errors.New("key cannot be nil")
	}
	var zapField zap.Field
	switch value.(type) {
	case int:
		zapField = zap.Int(key, to.Int(value))
	case float64:
		zapField = zap.Float64(key, to.Float64(value))
	case string:
		zapField = zap.String(key, to.String(value))
	case time.Time:
		zapField = zap.Time(key, to.Time(value))
	case bool:
		zapField = zap.Bool(key, to.Bool(value))
	default:
		zapField = zap.String(key, to.String(value))
	}
	return &zapField, nil
}

func ToZapLogLevel(level string) zapcore.Level {
	switch level {
	case "DEBUG":
		return zapcore.DebugLevel
	case "INFO":
		return zapcore.InfoLevel
	case "WARN":
		return zapcore.WarnLevel
	case "ERROR":
		return zapcore.ErrorLevel
	case "PANIC":
		return zapcore.PanicLevel
	case "FATAL":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

func ToZapLevelEncoder(encoder enumeration.EncodeLevel) zapcore.LevelEncoder {
	switch encoder {
	case enumeration.Lowercase:
		return zapcore.LowercaseLevelEncoder
	case enumeration.Camelcase:
		return zapcore.CapitalLevelEncoder
	default:
		return zapcore.LowercaseLevelEncoder
	}
}

func ToZapTimeEncoder(encoder enumeration.EncodeTime) zapcore.TimeEncoder {
	switch encoder {
	case enumeration.RFC3339Nano:
		return zapcore.RFC3339NanoTimeEncoder
	case enumeration.RFC3339:
		return zapcore.RFC3339TimeEncoder
	case enumeration.ISO8601:
		return zapcore.ISO8601TimeEncoder
	case enumeration.Milliseconds:
		return zapcore.EpochMillisTimeEncoder
	case enumeration.Nanoseconds:
		return zapcore.EpochNanosTimeEncoder
	default:
		return zapcore.EpochTimeEncoder
	}
}

func ToZapDurationEncoder(encoder enumeration.EncodeDuration) zapcore.DurationEncoder {
	switch encoder {
	case enumeration.StringDuration:
		return zapcore.StringDurationEncoder
	case enumeration.NanosecondDuration:
		return zapcore.NanosDurationEncoder
	case enumeration.MillisecondDuration:
		return zapcore.NanosDurationEncoder
	default:
		return zapcore.SecondsDurationEncoder
	}
}

func ToZapCallerEncoder(encoder enumeration.EncodeCaller) zapcore.CallerEncoder {
	switch encoder {
	case enumeration.LongestFunctionName:
		return zapcore.FullCallerEncoder
	case enumeration.ShortestFunctionName:
		return zapcore.ShortCallerEncoder
	default:
		return zapcore.FullCallerEncoder
	}
}
