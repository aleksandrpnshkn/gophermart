package config

import (
	"flag"
	"os"
)

type Config struct {
	RunAddress           string
	LogLevel             string
	JwtSecretKey         string
	DatabaseURI          string
	AccrualSystemAddress string
}

func New() *Config {
	config := Config{
		RunAddress:           "localhost:8081",
		LogLevel:             "info",
		JwtSecretKey:         "secret",
		DatabaseURI:          "postgres://admin:qwerty@localhost:5433/gophermart?sslmode=disable",
		AccrualSystemAddress: "http://localhost:8083",
	}

	envLogLevel, ok := os.LookupEnv("LOG_LEVEL")
	if ok {
		config.LogLevel = envLogLevel
	}
	flag.StringVar(&config.RunAddress, "a", config.RunAddress, "net address host:port")

	envRunAddress, ok := os.LookupEnv("RUN_ADDRESS")
	if ok {
		config.RunAddress = envRunAddress
	}
	flag.StringVar(&config.LogLevel, "l", config.LogLevel, "log level")

	envJwtSecretKey, ok := os.LookupEnv("JWT_SECRET_KEY")
	if ok {
		config.JwtSecretKey = envJwtSecretKey
	}
	flag.StringVar(&config.JwtSecretKey, "s", config.JwtSecretKey, "jwt secret key")

	envDatabaseURI, ok := os.LookupEnv("DATABASE_URI")
	if ok {
		config.DatabaseURI = envDatabaseURI
	}
	flag.StringVar(&config.DatabaseURI, "d", config.DatabaseURI, "database URI")

	envAccrualSystemAddress, ok := os.LookupEnv("ACCRUAL_SYSTEM_ADDRESS")
	if ok {
		config.AccrualSystemAddress = envAccrualSystemAddress
	}
	flag.StringVar(&config.AccrualSystemAddress, "r", config.AccrualSystemAddress, "accrual system address")

	flag.Parse()

	return &config
}
