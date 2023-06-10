package config

import (
	"errors"

	"github.com/spf13/viper"
)

type Config struct {
	ServerPort              string
	DatabaseDSN             string
	AwsRegion               string
	FirebaseCredentialsFile string
	JWTSecret               string
	SESSender               string
}

func LoadConfig() (Config, error) {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("APP")

	var config Config

	// Define all required keys
	requiredKeys := []string{
		"ServerPort",
		"DatabaseDSN",
		"AwsRegion",
		"FirebaseCredentialsFile",
		"JWTSecret",
		"SESSender",
	}

	// Check all required keys
	for _, key := range requiredKeys {
		viper.BindEnv(key)

		if !viper.IsSet(key) {
			return Config{}, errors.New("required key " + key + " missing value")
		}
	}

	config.ServerPort = viper.GetString("ServerPort")
	config.DatabaseDSN = viper.GetString("DatabaseDSN")
	config.AwsRegion = viper.GetString("AwsRegion")
	config.FirebaseCredentialsFile = viper.GetString("FirebaseCredentialsFile")
	config.JWTSecret = viper.GetString("JWTSecret")
	config.SESSender = viper.GetString("SESSender")

	return config, nil
}
