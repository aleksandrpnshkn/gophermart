package config

import (
	"flag"
	"os"
)

type Config struct {
	RunAddress   string
	LogLevel     string
	JwtSecretKey string
	DatabaseURI  string
}

func New() *Config {
	config := Config{
		RunAddress:   "localhost:8081",
		LogLevel:     "info",
		JwtSecretKey: "secret",
		DatabaseURI:  "postgres://admin:qwerty@localhost:5433/gophermart?sslmode=disable",
	}

	flag.StringVar(&config.RunAddress, "a", config.RunAddress, "net address host:port")
	flag.StringVar(&config.JwtSecretKey, "s", config.JwtSecretKey, "jwt secret key")
	flag.StringVar(&config.DatabaseURI, "d", config.DatabaseURI, "database URI")

	flag.Parse()

	envRunAddress, ok := os.LookupEnv("RUN_ADDRESS")
	if ok {
		config.RunAddress = envRunAddress
	}

	envJwtSecretKey, ok := os.LookupEnv("JWT_SECRET_KEY")
	if ok {
		config.JwtSecretKey = envJwtSecretKey
	}

	envDatabaseURI, ok := os.LookupEnv("DATABASE_URI")
	if ok {
		config.DatabaseURI = envDatabaseURI
	}

	return &config
}
