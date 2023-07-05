// Package logger provides functionality to handle logging.
package logger

import (
	"log"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	LevelDebug = "DEBUG"
	LevelInfo  = "INFO"
	LevelWarn  = "WARN"
	LevelError = "ERROR"
)

// Logger is wrapper around *zap.SugaredLogger that will handle all logging behavior.
type Logger struct{ *zap.SugaredLogger }

// New creates a new Logger instance accepting logger level and path for log file.
func New(lvl, pth string) *Logger {
	core := zapcore.NewTee(getConsoleCore(lvl), getJSONCore(pth, lvl))
	return &Logger{zap.New(core).Sugar()}
}

func (l *Logger) ToStandard() *log.Logger {
	return zap.NewStdLog(l.Desugar())
}

func getConsoleCore(lvl string) zapcore.Core {
	cfg := getEncoderConfig()
	cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	return zapcore.NewCore(
		zapcore.NewConsoleEncoder(cfg),
		zapcore.Lock(os.Stderr),
		getLevel(lvl),
	)
}

func getJSONCore(pth, lvl string) zapcore.Core {
	cfg := getEncoderConfig()

	const perm = 0644

	file, err := os.OpenFile(pth, os.O_APPEND|os.O_CREATE|os.O_WRONLY, perm)
	if err != nil {
		log.Fatalf("creating logger file: %v\n", err)
	}

	return zapcore.NewCore(
		zapcore.NewJSONEncoder(cfg),
		zapcore.Lock(file),
		getLevel(lvl),
	)
}

func getLevel(lvl string) zapcore.Level {
	switch strings.ToLower(lvl) {
	case "debug":
		return zap.DebugLevel
	case "error":
		return zap.ErrorLevel
	case "info":
		return zap.InfoLevel
	default:
		return zap.DebugLevel
	}
}

func getEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:          "message",
		LevelKey:            "level",
		TimeKey:             "time",
		NameKey:             "name",
		CallerKey:           "caller",
		FunctionKey:         "",
		StacktraceKey:       "stacktrace",
		SkipLineEnding:      false,
		LineEnding:          "\n",
		EncodeLevel:         zapcore.CapitalLevelEncoder,
		EncodeTime:          zapcore.ISO8601TimeEncoder,
		EncodeDuration:      zapcore.NanosDurationEncoder,
		EncodeCaller:        zapcore.ShortCallerEncoder,
		NewReflectedEncoder: nil,
		ConsoleSeparator:    "\t",
	}
}
