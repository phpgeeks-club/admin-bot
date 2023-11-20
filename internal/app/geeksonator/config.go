package geeksonator

import (
	"fmt"

	"github.com/caarlos0/env/v10"
)

// Config represents application configuration.
type Config struct {
	TgBotToken       string `env:"GEEKSONATOR_TELEGRAM_BOT_TOKEN"`
	TgTimeoutSeconds int    `env:"GEEKSONATOR_TELEGRAM_TIMEOUT_SECONDS" envDefault:"15"`
	DebugMode        bool   `env:"GEEKSONATOR_DEBUG_MODE"`
	DebugTgBotToken  string `env:"GEEKSONATOR_DEBUG_TELEGRAM_BOT_TOKEN"`
}

// LoadConfig loads application configuration.
func LoadConfig() (*Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("env.Parse(): %v", err)
	}

	return &cfg, nil
}
