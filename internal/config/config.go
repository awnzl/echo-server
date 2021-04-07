package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	Port     string
	LogLevel string
}

func Get() Config {
	_ = godotenv.Load()

	return Config{
		Port:     os.Getenv("SERVER_PORT"),
		LogLevel: os.Getenv("LOG_LEVEL"),
	}
}
