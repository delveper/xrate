package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
)

// Logger is wrapper around *zap.SugaredLogger that will handle all logging behavior.
type Logger struct{ *zap.SugaredLogger }

func New() *Logger {
	core := zapcore.NewTee(getConsoleCore())
	return &Logger{zap.New(core).Sugar()}
}

func (l *Logger) ToStandard() *log.Logger {
	return zap.NewStdLog(l.Desugar())
}

func getConsoleCore() zapcore.Core {
	cfg := getEncoderConfig()
	cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	return zapcore.NewCore(
		zapcore.NewConsoleEncoder(cfg),
		zapcore.Lock(os.Stderr),
		zapcore.DebugLevel,
	)
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
