// Package logger provides functionality to handle logging.
package logger

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	LevelDebug = "DEBUG"
	LevelInfo  = "INFO"
	LevelWarn  = "WARN"
	LevelError = "ERROR"
)

type Config struct {
	Level     string   `default:"debug"`
	FilePath  string   `default:"sys.log"`
	Topic     string   `default:"logs"`
	Partition int32    `default:"0"`
	Brokers   []string `default:"localhost:9092,kafka:9092"`
	Offset    string   `default:"oldest"`
}

// Logger is wrapper around *zap.SugaredLogger that will handle all logging behavior.
type Logger struct{ *zap.SugaredLogger }

var _ zapcore.Core = (*kafkaCore)(nil)

type kafkaCore struct {
	levelEnabler zapcore.LevelEnabler
	encoder      zapcore.Encoder
	producer     sarama.SyncProducer
	topic        string
}

// New creates a new Logger instance accepting logger level and path for log file.
func New(cores ...zapcore.Core) *Logger {
	return &Logger{zap.New(zapcore.NewTee(cores...)).Sugar()}
}

// ToStandard returns a standard logger.
func (l *Logger) ToStandard() *log.Logger {
	return zap.NewStdLog(l.Desugar())
}

func WithConsoleCore(lvl string) zapcore.Core {
	cfg := getEncoderConfig()
	cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	return zapcore.NewCore(
		zapcore.NewConsoleEncoder(cfg),
		zapcore.Lock(os.Stderr),
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

func WithJSONCore(lvl string, paths ...string) zapcore.Core {
	if paths == nil {
		paths = []string{"sys.log"}
	}

	const perm = 0644

	file, err := os.OpenFile(path.Join(paths...), os.O_APPEND|os.O_CREATE|os.O_WRONLY, perm)
	if err != nil {
		log.Fatalf("creating logger file: %v\n", err)
	}

	return zapcore.NewCore(
		zapcore.NewJSONEncoder(getEncoderConfig()),
		zapcore.Lock(file),
		getLevel(lvl),
	)
}

func WithKafkaCore(lvl, topic string, brokers ...string) zapcore.Core {
	prod, err := sarama.NewSyncProducer(brokers, getKafkaConfig())
	if err != nil {
		log.Fatalf("creating kafka producer: %v\n", err)
	}

	return &kafkaCore{
		levelEnabler: zap.LevelEnablerFunc(func(level zapcore.Level) bool { return level >= getLevel(lvl) }),
		encoder:      zapcore.NewJSONEncoder(getEncoderConfig()),
		producer:     prod,
		topic:        topic,
	}
}

func getKafkaConfig() *sarama.Config {
	const maxRetry = 5

	cfg := sarama.NewConfig()
	cfg.Version = sarama.V2_8_0_0
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Retry.Max = maxRetry
	cfg.Producer.Timeout = maxRetry * time.Second
	cfg.Producer.Compression = sarama.CompressionSnappy
	cfg.Producer.Return.Successes = true

	return cfg
}

// Enabled returns whether the Core will accept entries that have the specified level.
func (k *kafkaCore) Enabled(level zapcore.Level) bool {
	return k.levelEnabler.Enabled(level)
}

// With adds structured context to the Core.
func (k *kafkaCore) With(fields []zapcore.Field) zapcore.Core {
	clone := *k
	for i := range fields {
		fields[i].AddTo(clone.encoder)
	}

	return &clone
}

// Check determines whether the supplied Entry should be logged.
func (k *kafkaCore) Check(entry zapcore.Entry, checked *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if k.Enabled(entry.Level) {
		return checked.AddCore(entry, k)
	}

	return checked
}

// Write serializes the Entry and any Fields supplied at the log site and writes.
func (k *kafkaCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	buf, err := k.encoder.EncodeEntry(entry, fields)
	if err != nil {
		return err
	}

	message := &sarama.ProducerMessage{
		Topic: k.topic,
		Value: sarama.ByteEncoder(buf.Bytes()),
	}

	part, offset, err := k.producer.SendMessage(message)
	if err != nil {
		return fmt.Errorf("storing message in kafka partition(%d) & offset(%d): %w", part, offset, err)
	}

	return nil
}

// Sync flushes buffered logs (if any).
func (k *kafkaCore) Sync() error {
	// In this case, there's no need to do anything on Sync, because
	// sarama.SyncProducer.SendMessage is a synchronous operation.
	return nil
}
