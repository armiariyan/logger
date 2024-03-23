package logger

import (
	"context"
	"fmt"
)

type Options struct {
	Name       string        `json:"name"`
	SysOptions OptionsLogger `json:"sysOptions"`
	TdrOptions OptionsLogger `json:"tdrOptions"`
}

type OptionsLogger struct {
	Type         string       `json:"type"`
	OptionsFile  OptionsFile  `json:"optionsFile"`
	OptionsQueue OptionsQueue `json:"optionsQueue"`
}

type combineLogger struct {
	sysLog Logger
	tdrLog Logger
}

func (c *combineLogger) Debug(ctx context.Context, message string, fields ...Field) {
	c.sysLog.Debug(ctx, message, fields...)
}

func (c *combineLogger) Info(ctx context.Context, message string, fields ...Field) {
	c.sysLog.Info(ctx, message, fields...)
}

func (c *combineLogger) Warn(ctx context.Context, message string, fields ...Field) {
	c.sysLog.Warn(ctx, message, fields...)
}

func (c *combineLogger) Error(ctx context.Context, message string, fields ...Field) {
	c.sysLog.Error(ctx, message, fields...)
}

func (c *combineLogger) Fatal(ctx context.Context, message string, fields ...Field) {
	c.sysLog.Fatal(ctx, message, fields...)
}

func (c *combineLogger) Panic(ctx context.Context, message string, fields ...Field) {
	c.sysLog.Panic(ctx, message, fields...)
}

func (c *combineLogger) TDR(ctx context.Context, tdr LogTdrModel) {
	c.tdrLog.TDR(ctx, tdr)
}

func (c *combineLogger) Close() error {
	var err error

	if _err := c.sysLog.Close(); _err != nil {
		err = fmt.Errorf("error close syslog: %w", _err)
	}

	if _err := c.tdrLog.Close(); _err != nil {
		err = fmt.Errorf("error close tdrlog: %w", _err)
	}

	return err
}

func SetupLoggerCombine(options Options) Logger {
	fmt.Println("Try newLogger ...")

	var sysLog Logger
	switch options.SysOptions.Type {
	case Queue:
		sysLog = SetupLoggerQueue(options.Name, &options.SysOptions.OptionsQueue)
	case File:
		sysLog = SetupLoggerFile(options.Name, &options.SysOptions.OptionsFile)
	default:
		panic("syslog not found")
	}

	var tdrLog Logger
	switch options.TdrOptions.Type {
	case Queue:
		tdrLog = SetupLoggerQueue(options.Name, &options.TdrOptions.OptionsQueue)
	case File:
		tdrLog = SetupLoggerFile(options.Name, &options.TdrOptions.OptionsFile)
	default:
		panic("tdrLog not found")
	}

	return &combineLogger{
		sysLog: sysLog,
		tdrLog: tdrLog,
	}
}
