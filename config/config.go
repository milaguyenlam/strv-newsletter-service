package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	ServerPort   int    `mapstructure:"SERVER_PORT"`
	DatabaseDSN  string `mapstructure:"DATABASE_DSN"`
	AwsRegion    string `mapstructure:"AWS_REGION"`
	JWTSecret    string `mapstructure:"JWT_SECRET"`
	FirebaseJSON []byte `mapstructure:"FIREBASE_JSON"`
}

func LoadConfig() (*Config, error) {
	viper.AutomaticEnv()
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("Loading config from env: %w", err)

	}
	return &config, nil
}
