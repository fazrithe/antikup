package config

import (
	"github.com/BurntSushi/toml"
	"log"
)

type Config struct {
	AppName string `toml:"APP_NAME"`

	DbHost     string `toml:"DB_HOST"`
	DbUsername string `toml:"DB_USERNAME"`
	DbPassword string `toml:"DB_PASSWORD"`
	DbDatabase string `toml:"DB_DATABASE"`

	ServerAddr string `toml:"SERVER_ADDR"`
	CertFile   string `toml:"CERT_FILE"`
	KeyFile    string `toml:"KEY_FILE"`

	AppKey string `toml:"APP_KEY"`
}

func (c *Config) Read() *Config {
	if _, err := toml.DecodeFile("config.toml", c); err != nil {
		log.Fatalf("Error reading config file: %v", err.Error())
	}
	return c
}
