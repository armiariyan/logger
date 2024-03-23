package logger

import (
	"fmt"
)

type OptionsQueue struct {
	Type     string          `json:"type"`
	Topic    string          `json:"topic"`
	Producer ProducerOptions `json:"producer"`
	Mask     bool            `json:"mask"`
	Level    Level           `json:"level"`
}

type ProducerOptions struct {
	Address         string `json:"address"`
	RetryMax        int    `json:"retryMax"`
	ReturnSuccesses bool   `json:"returnSuccesses"`
}

// SetupLoggerQueue will return legacy Logger using Queue interface with new logic using Logger
func SetupLoggerQueue(serviceName string, config *OptionsQueue) Logger {
	fmt.Println("Try newLogger Queue...")

	if config == nil {
		panic("legacy logger queue config is nil")
	}

	if config.Type != QueueTypeKafka {
		panic(fmt.Errorf("legacy logger queue unsupported queue type %s", config.Type))
	}

	var opt = make([]Option, 0)

	if config.Mask {
		opt = append(opt, MaskEnabled())
	}

	opt = append(opt, WithKafkaOutput(config))
	opt = append(opt, WithLevel(config.Level))

	log, err := newLogger(opt...)
	if err != nil {
		panic(fmt.Errorf("init legacy logger with mode %s error: %w", Queue, err))
	}

	return log
}
