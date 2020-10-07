package main

import (
	"fmt"
	"os"
)

type Config struct {
	Addr        string
	DatabaseURL string
	SigningKey  string
}

func (c Config) validate() error {
	if c.Addr == "" || c.Addr == ":" {
		return fmt.Errorf("Config.validate: addr is required")
	}
	if c.DatabaseURL == "" {
		return fmt.Errorf("Config.validate: database url is required")
	}
	if c.SigningKey == "" {
		return fmt.Errorf("Config.validate: signing key is required")
	}
	return nil
}

func loadConfig() (Config, error) {
	env := os.Getenv("ENV")
	if env == "production" {
		return loadProductionConfig()
	}
	return loadDevelopmentConfig()
}

func loadDevelopmentConfig() (Config, error) {
	cfg := Config{
		Addr:        ":8080",
		DatabaseURL: os.Getenv("DATABASE_URL"),
		SigningKey:  "statickey",
	}
	return cfg, cfg.validate()
}

func loadProductionConfig() (Config, error) {
	cfg := Config{
		Addr:        ":" + os.Getenv("PORT"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		SigningKey:  os.Getenv("JWT_SIGNING_KEY"),
	}
	return cfg, cfg.validate()
}
