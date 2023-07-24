package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/GenesisEducationKyiv/main-project-delveper/sys/env"
	"github.com/GenesisEducationKyiv/main-project-delveper/sys/logger"
	"github.com/IBM/sarama"
)

type config struct {
	Log struct {
		Level     string   `default:"debug"`
		Topic     string   `default:"logs"`
		Partition int32    `default:"0"`
		Brokers   []string `default:"localhost:9092,kafka:9092"`
		Offset    string   `default:"newest"`
	}
}

func main() {
	var cfg config
	if err := env.ParseTo(&cfg); err != nil {
		log.Fatalf("failed to parse env: %v", err)
	}

	log := logger.New(logger.WithConsoleCore(cfg.Log.Level))
	defer log.Sync()

	if err := run(log, &cfg); err != nil {
		log.Errorw("startup error", "error", err)
		os.Exit(1)
	}
}

func run(log *logger.Logger, cfg *config) error {
	log.Infow("starting kafka log consumer")

	cs, csPart, err := createKafkaConsumer(cfg)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	go consumeMessages(ctx, log, csPart)

	select {
	case <-shutdown:
		log.Infow("shutting down")
		cancel()
	case <-ctx.Done():
		log.Infow("context cancelled")
	}

	return closeKafkaConsumer(cs, csPart, log)
}

func createKafkaConsumer(cfg *config) (sarama.Consumer, sarama.PartitionConsumer, error) {
	cs, err := sarama.NewConsumer(cfg.Log.Brokers, sarama.NewConfig())
	if err != nil {
		return nil, nil, fmt.Errorf("creating kafka cs: %w", err)
	}

	csPart, err := cs.ConsumePartition(cfg.Log.Topic, cfg.Log.Partition, parseOffset(cfg.Log.Offset))
	if err != nil {
		return nil, nil, fmt.Errorf("creating kafka partition cs: %w", err)
	}

	return cs, csPart, nil
}

func closeKafkaConsumer(cs sarama.Consumer, csPart sarama.PartitionConsumer, log *logger.Logger) error {
	if err := csPart.Close(); err != nil {
		log.Errorw("closing kafka partition consumer", "error", err)
	}

	if err := cs.Close(); err != nil {
		log.Errorw("closing kafka consumer", "error", err)
	}

	return nil
}

func consumeMessages(ctx context.Context, log *logger.Logger, csPart sarama.PartitionConsumer) {
	for {
		select {
		case msg := <-csPart.Messages():
			log.Infow("kafka log",
				"topic", msg.Topic,
				"partition", msg.Partition,
				"offset", msg.Offset,
				"key", msg.Key,
				"value", sarama.StringEncoder(msg.Value))

		case <-ctx.Done():
			return
		}
	}
}

func parseOffset(str string) int64 {
	switch strings.ToLower(str) {
	case "oldest":
		return sarama.OffsetOldest
	case "newest":
		return sarama.OffsetNewest
	default:
		return sarama.OffsetNewest
	}
}
