package config

import (
	"errors"

	"github.com/spf13/viper"
)

type Config struct {
	ServerPort              int
	DatabaseDSN             string
	AwsRegion               string
	FirebaseCredentialsFile string
	JWTSecret               string
}

func LoadConfig() (Config, error) {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("APP")

	var config Config

	// Define all required keys
	requiredKeys := []string{
		"SERVER_PORT",
		"DATABASE_DSN",
		"AWS_REGION",
		"FIREBASE_CREDENTIALS_FILE",
		"JWT_SECRET",
	}

	// Check all required keys
	for _, key := range requiredKeys {
		viper.BindEnv(key)

		if !viper.IsSet(key) {
			return Config{}, errors.New("required key " + key + " missing value")
		}
	}

	config.ServerPort = viper.GetInt("SERVER_PORT")
	config.DatabaseDSN = viper.GetString("DATABASE_DSN")
	config.AwsRegion = viper.GetString("AWS_REGION")
	config.FirebaseCredentialsFile = viper.GetString("FIREBASE_CREDENTIALS_FILE")
	config.JWTSecret = viper.GetString("JWT_SECRET")

	return config, nil
}
