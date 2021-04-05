package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	Port     string
	LogLevel string
}

func Get() (Config, error) {
	if err := godotenv.Load(); err != nil {
		return Config{}, err
	}

	return Config{
		Port:     os.Getenv("SERVER_PORT"),
		LogLevel: os.Getenv("LOG_LEVEL"),
	}, nil
}
