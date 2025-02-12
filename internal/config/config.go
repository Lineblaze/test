package config

import (
	"fmt"
	"os"
)

type Config struct {
	ServiceName string `json:"serviceName"`

	Postgres struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		DBName   string `json:"dbName"`
		SSLMode  string `json:"sslMode"`
		PgDriver string `json:"pgDriver"`
	} `json:"postgres"`

	Server struct {
		Port string `json:"port"`
	} `json:"server"`

	Logger struct {
		MinLevel int8 `json:"minLevel"`
	} `json:"logger"`

	Auth struct {
		Secret string `json:"secret"`
	} `json:"auth"`
}

func LoadConfig() (*Config, error) {
	requiredEnvVars := []string{
		"POSTGRES_HOST", "POSTGRES_PORT", "POSTGRES_USER", "POSTGRES_PASSWORD", "POSTGRES_DB",
		"SERVER_PORT", "JWT_SECRET_KEY",
	}
	for _, env := range requiredEnvVars {
		if os.Getenv(env) == "" {
			return nil, fmt.Errorf("%s is not set in the environment", env)
		}
	}

	cfg := &Config{
		ServiceName: "Avito Test",
		Postgres: struct {
			Host     string `json:"host"`
			Port     string `json:"port"`
			User     string `json:"user"`
			Password string `json:"password"`
			DBName   string `json:"dbName"`
			SSLMode  string `json:"sslMode"`
			PgDriver string `json:"pgDriver"`
		}{
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_PORT"),
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			DBName:   os.Getenv("POSTGRES_DB"),
			SSLMode:  "disable",
			PgDriver: "pgx",
		},
		Server: struct {
			Port string `json:"port"`
		}{
			Port: os.Getenv("SERVER_PORT"),
		},
		Logger: struct {
			MinLevel int8 `json:"minLevel"`
		}{
			MinLevel: -4,
		},
		Auth: struct {
			Secret string `json:"secret"`
		}{
			Secret: os.Getenv("JWT_SECRET_KEY"),
		},
	}

	return cfg, nil
}
