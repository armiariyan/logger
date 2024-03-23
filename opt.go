package logger

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/go-playground/validator/v10"
	rotateLogs "github.com/lestrrat-go/file-rotatelogs"
)

type Option func(*defaultLogger) error

func OptNoop() Option {
	return func(logger *defaultLogger) error {
		logger.noopLogger = true
		return nil
	}
}

func MaskEnabled() Option {
	return func(logger *defaultLogger) error {
		logger.maskEnabled = true
		return nil
	}
}

func WithStdout() Option {
	return func(logger *defaultLogger) error {
		// Wire STD output for both type
		logger.writers = append(logger.writers, os.Stdout)
		return nil
	}
}

func WithFileOutput(conf *OptionsFile) Option {
	return func(logger *defaultLogger) error {
		err := validator.New().Struct(conf)
		if err != nil {
			return fmt.Errorf("config for file output error: %w", err)
		}

		outputSys, err := rotateLogs.New(
			conf.FileLocation+".%Y%m%d",
			rotateLogs.WithLinkName(conf.FileLocation),
			rotateLogs.WithMaxAge(conf.FileMaxAge*24*time.Hour),
			rotateLogs.WithRotationTime(time.Hour),
		)

		if err != nil {
			return fmt.Errorf("sys file error: %w", err)
		}

		// Wire SYS config only in sys
		logger.writers = append(logger.writers, outputSys)
		logger.closer = append(logger.closer, outputSys)
		return nil
	}
}

type wrapKafkaWriter struct {
	topic    string
	producer sarama.SyncProducer
}

func (w *wrapKafkaWriter) Write(p []byte) (n int, err error) {
	_, _, err = w.producer.SendMessage(&sarama.ProducerMessage{
		Topic: w.topic,
		Key:   sarama.StringEncoder(fmt.Sprint(time.Now().UTC())),
		Value: sarama.ByteEncoder(p),

		// Below this point are filled in by the producer as the message is processed
		Offset:    0,
		Partition: 0,
		Timestamp: time.Time{},
	})

	return
}

var _ io.Writer = (*wrapKafkaWriter)(nil)

// WithKafkaOutput can be called multiple times to add functionality
// where we want to broadcast log into different kafka cluster
// or topic.
func WithKafkaOutput(conf *OptionsQueue) Option {
	return func(logger *defaultLogger) error {
		err := validator.New().Struct(conf)
		if err != nil {
			return fmt.Errorf("config for kafka output error: %w", err)
		}

		saramaConf := sarama.NewConfig()
		saramaConf.Producer.Partitioner = sarama.NewRandomPartitioner
		saramaConf.Producer.RequiredAcks = sarama.WaitForAll
		saramaConf.Producer.Retry.Max = conf.Producer.RetryMax
		saramaConf.Producer.Return.Successes = conf.Producer.ReturnSuccesses

		kafkaProducer, err := sarama.NewSyncProducer(strings.Split(conf.Producer.Address, ","), saramaConf)
		if err != nil {
			return fmt.Errorf("kafka sync producer error: %w", err)
		}

		kafkaWriter := &wrapKafkaWriter{
			topic:    conf.Topic,
			producer: kafkaProducer,
		}

		// wire Kafka writer to log
		logger.writers = append(logger.writers, kafkaWriter)
		logger.closer = append(logger.closer, kafkaProducer)
		return nil
	}
}

// WithCustomWriter add custom writer, so you can write using any storage method
// without waiting this package to be updated.
func WithCustomWriter(writer io.WriteCloser) Option {
	return func(logger *defaultLogger) error {
		if writer == nil {
			return fmt.Errorf("writer is nil")
		}

		// wire custom writer to log
		logger.writers = append(logger.writers, writer)
		logger.closer = append(logger.closer, writer)
		return nil
	}
}

// WithLevel set level of logger
func WithLevel(level Level) Option {
	return func(logger *defaultLogger) error {
		logger.level = level
		return nil
	}
}
