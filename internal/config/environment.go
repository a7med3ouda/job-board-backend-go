package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port  string
	DbConnectionString string
}

func Load() (*Config, error) {
	_ = godotenv.Load()
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbConnectionString := os.Getenv("DB_CONNECTION_STRING")
	if dbConnectionString == "" {
		return nil, fmt.Errorf("error: DB_CONNECTION_STRING environment variable is required")
	}

	return &Config{
		Port:  port,
		DbConnectionString: dbConnectionString,
	}, nil
}