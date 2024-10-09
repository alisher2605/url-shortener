package config

import (
	"fmt"
	validation "github.com/alisher2605/url-shortener/util/validator"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	ConfigDirectory = "config"
	ConfigName      = "config"
	ConfigFormat    = "json"
)

type Configuration struct {
	MaxAge   int      `mapstructure:"max_age" validate:"required"`
	UrlTtl   int      `mapstructure:"url_ttl"  validate:"required"`
	AppPort  string   `mapstructure:"app_port" validate:"required"`
	Database Database `mapstructure:"database" validate:"required"`
}

type Database struct {
	RetryAttempts   int    `mapstructure:"retry_attempts" validate:"required"`
	Region          string `mapstructure:"region" validate:"required"`
	Endpoint        string `mapstructure:"endpoint" validate:"required"`
	AccessKeyId     string `mapstructure:"access_key_id" validate:"required"`
	AccessKeySecret string `mapstructure:"access_key_secret" validate:"required"`
}

func OpenConfig() *Configuration {
	v := viper.New()

	v.SetConfigName(ConfigName)
	v.AddConfigPath(fmt.Sprintf("./%s", ConfigDirectory))
	v.SetConfigType(ConfigFormat)
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		zap.S().Fatalf("can't open the config file: %v", err)
	}

	config := new(Configuration)

	if err := v.Unmarshal(config); err != nil {
		zap.S().Fatalf("can't u the config file: %v", err)
	}

	validator := validation.NewValidator()

	err := validator.Validator.Struct(config)
	if err != nil {
		zap.S().Fatal("Invalid config: %v", err)
	}

	return config
}
