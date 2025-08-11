package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	ApiHost      string
	TemporalHost string
}

func Load() (*Config, error) {
	envPath := filepath.Join(".env")
	if err := godotenv.Load(envPath); err != nil {
		log.Printf("Warning: could not load .env file: %s", err)
	}

	apiPort := os.Getenv("API_PORT")
	if apiPort == "" {
		return nil, fmt.Errorf("API_PORT env variable is not set")
	}
	temporalPort := os.Getenv("TEMPORAL_HOST_PORT")
	if temporalPort == "" {
		return nil, fmt.Errorf("TEMPORAL_HOST_PORT env variable is not set")
	}

	return &Config{
		ApiHost:      apiPort,
		TemporalHost: temporalPort,
	}, nil
}
