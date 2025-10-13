package config

import (
	"context"
	"log/slog"
	"os"

	"github.com/spf13/viper"
	"github.com/yuanyu90221/uniplile-authentication-with-linkedin/internal/logger"
)

// Config - App config.
type Config struct {
	Port               string `mapstructure:"PORT"`
	GinMode            string `mapstructure:"GIN_MODE"`
	DBURL              string `mapstructure:"DB_URL"`
	UnipileBaseURL     string `mapstructure:"UNIPILE_BASE_URL"`
	UnipileAccessToken string `mapstructure:"UNIPILE_ACCESS_TOKEN"`
}

// AppConfig - global config.
var AppConfig *Config

func Init(ctx context.Context) {
	log := logger.FromContext(ctx)
	v := viper.New()
	v.AddConfigPath(".")
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AutomaticEnv()

	failOnError(v.BindEnv("PORT"), "Failed on Bind PORT", log)
	failOnError(v.BindEnv("GIN_MODE"), "Failed on Bind GIN_MODE", log)
	failOnError(v.BindEnv("DB_URL"), "Failed on Bind GIN_MODE", log)
	failOnError(v.BindEnv("UNIPILE_BASE_URL"), "Failed on Bind UNIPILE_BASE_URL", log)
	failOnError(v.BindEnv("UNIPILE_ACCESS_TOKEN"), "Failed on Bind UNIPILE_ACCESS_TOKEN", log)
	err := v.ReadInConfig()
	if err != nil {
		log.WarnContext(ctx, "Load from environment variable")
	}
	err = v.Unmarshal(&AppConfig)
	if err != nil {
		failOnError(err, "Failed to read enivronment", log)
	}
}

func failOnError(err error, msg string, log *slog.Logger) {
	if err != nil {
		log.Error(msg)
		os.Exit(1)
	}
}
