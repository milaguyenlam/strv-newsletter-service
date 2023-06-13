package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	ServerPort   int    `mapstructure:"SERVER_PORT"`
	DatabaseURL  string `mapstructure:"DATABASE_URL"`
	AwsRegion    string `mapstructure:"AWS_REGION"`
	JWTSecret    string `mapstructure:"JWT_SECRET"`
	FirebaseJSON string `mapstructure:"FIREBASE_JSON"`
}

var requiredVariables = []string{
	"SERVER_PORT", "DATABASE_URL", "AWS_REGION", "JWT_SECRET", "FIREBASE_JSON",
}

func LoadConfig() (*Config, error) {
	viper.AutomaticEnv()

	for _, v := range requiredVariables {
		err := viper.BindEnv(v)
		if err != nil {
			return nil, fmt.Errorf("Loading config from env: %w", err)
		}
	}

	var config Config

	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("Loading config from env: %w", err)

	}
	return &config, nil
}
