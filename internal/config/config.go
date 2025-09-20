package config

import (
	"flag"
	"os"
)

type Config struct {
	RunAddress string
	LogLevel   string
}

func New() *Config {
	config := Config{
		RunAddress: "localhost:8081",
		LogLevel:   "info",
	}

	flag.StringVar(&config.RunAddress, "a", config.RunAddress, "net address host:port")

	flag.Parse()

	envRunAddress, ok := os.LookupEnv("RUN_ADDRESS")
	if ok {
		config.RunAddress = envRunAddress
	}

	return &config
}
