package bootstrap

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger(cfg LogConfig) (*zap.Logger, error) {
	level := zapcore.InfoLevel
	if err := level.Set(strings.ToLower(cfg.Level)); err != nil {
		level = zapcore.InfoLevel
	}

	encoding := cfg.Encoding
	if encoding == "" {
		encoding = "json"
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "ts"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	writeSyncer, err := buildWriteSyncer(cfg.Output)
	if err != nil {
		return nil, err
	}

	var encoder zapcore.Encoder
	if strings.ToLower(encoding) == "console" {
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	}

	core := zapcore.NewCore(encoder, writeSyncer, level)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	return logger, nil
}

func buildWriteSyncer(outputs []string) (zapcore.WriteSyncer, error) {
	if len(outputs) == 0 {
		outputs = []string{"stdout"}
	}

	var syncers []zapcore.WriteSyncer
	for _, out := range outputs {
		switch strings.ToLower(strings.TrimSpace(out)) {
		case "stdout":
			syncers = append(syncers, zapcore.AddSync(os.Stdout))
		case "stderr":
			syncers = append(syncers, zapcore.AddSync(os.Stderr))
		case "":
			continue
		default:
			f, err := os.OpenFile(out, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return nil, err
			}
			syncers = append(syncers, zapcore.AddSync(f))
		}
	}

	if len(syncers) == 0 {
		syncers = append(syncers, zapcore.AddSync(os.Stdout))
	}

	return zapcore.NewMultiWriteSyncer(syncers...), nil
}
