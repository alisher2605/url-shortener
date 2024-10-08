package main

import (
	"github.com/alisher2605/url-shortener/config"
	"github.com/alisher2605/url-shortener/internal/database"
	"github.com/alisher2605/url-shortener/internal/http"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
)

// @title           url-shortener
// @version         1.0
// @description    A simple URL shortener running on AWS
func main() {
	err := setupLogger(os.Getenv("DEBUG") != "")
	if err != nil {
		log.Fatalf("Can't initialize logger: %v", err)
	}

	configuration := config.OpenConfig()

	db := database.NewDatabase(configuration.Database)
	db.Connect()

	http.NewServer(configuration.AppPort, configuration.MaxAge).Run()
}

func encoderConf() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "severity",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func initConfiguredLogger(externalConfig zap.Config) error {
	logger, err := externalConfig.Build()
	if err != nil {
		return err
	}

	zap.ReplaceGlobals(logger)
	zap.RedirectStdLog(logger)

	return nil
}

func setupLogger(debug bool) error {
	level := zap.InfoLevel
	if debug {
		level = zap.DebugLevel
	}

	const samplingCongValue = 100

	var defaultJsonConfig = zap.Config{
		Level:         zap.NewAtomicLevelAt(level),
		Development:   debug,
		DisableCaller: true,
		Sampling: &zap.SamplingConfig{
			Initial:    samplingCongValue,
			Thereafter: samplingCongValue,
		},
		Encoding:         "json",
		EncoderConfig:    encoderConf(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	zap.AddStacktrace(zap.ErrorLevel)

	return initConfiguredLogger(defaultJsonConfig)
}
