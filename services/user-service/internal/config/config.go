package config

import (
	"log"
	"sync"

	"github.com/spf13/viper"
)

var (
	config *Config
	once   sync.Once
)

type Config struct {
	Port        string `mapstructure:"PORT"`
	Environment string `mapstructure:"ENVIRONMENT"`
	AppName     string `mapstructure:"APP_NAME"`
	DebugMode   bool   `mapstructure:"DEBUG"`

	PostgresDSN string `mapstructure:"DB_POSTGRES_DSN"`

	JWTAccessKey            string `mapstructure:"JWT_ACCESS_KEY"`
	JWTRefreshKey           string `mapstructure:"JWT_REFRESH_KEY"`
	JWTAccessTokenDuration  string `mapstructure:"JWT_ACCESS_DURATION"`
	JWTRefreshTokenDuration string `mapstructure:"JWT_REFRESH_DURATION"`
}

func GetConfig() *Config {
	once.Do(func() {
		viper.SetConfigName(".env")
		viper.SetConfigType("env")
		viper.AddConfigPath(".")
		viper.AutomaticEnv()

		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("error reading config file, %s", err)
		}

		if err := viper.Unmarshal(&config); err != nil {
			log.Fatalf("error unmarshalling config, %s", err)
		}
	})

	return config
}
