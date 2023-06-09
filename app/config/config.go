package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	AppPort    string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
}

func LoadConfig() (config Config, err error) {
	viper.AutomaticEnv() // It reads from environment variables

	config = Config{
		AppPort:    viper.GetString("APP_PORT"),    // MYAPP_SERVER_PORT
		DBHost:     viper.GetString("DB_HOST"),     // MYAPP_DB_HOST
		DBPort:     viper.GetString("DB_PORT"),     // MYAPP_DB_PORT
		DBUser:     viper.GetString("DB_USER"),     // MYAPP_DB_USER
		DBPassword: viper.GetString("DB_PASSWORD"), // MYAPP_DB_PASSWORD
		DBName:     viper.GetString("DB_NAME"),     // MYAPP_DB_NAME
		JWTSecret:  viper.GetString("JWT_SECRET"),  // MYAPP_JWT_SECRET
	}
	return
}
