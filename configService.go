package main

import (
	"CazaAlquileres/config"
	"CazaAlquileres/database"
	"github.com/sirupsen/logrus"
)

var cfg Config

type Config struct {
	DBConfig 			database.DBConfig
	SenderEmail			string
	SenderPassword		string
}

func LoadConfig() *Config {

	var err error
	cfg = Config{}

	provider, err := load()

	if err != nil {
		logrus.Fatalf("Error loading Cfg: %s", err)
	}

	err = provider.Unmarshal(&cfg)

	if err != nil {
		logrus.Fatalf("Error parsing Cfg: %s", err)
	}

	return &cfg
}

func load() (config.Provider, error) {
	cfg := config.NewConfig()

	cfg.WithProperty("sender-email", true).Alias("SenderEmail").EnvAlias("SENDER_EMAIL")
	cfg.WithProperty("sender-password", true).Alias("SenderPassword").EnvAlias("SENDER_PASSWORD")

	database.Load(cfg)

	cfg.SetFileName("config")
	cfg.SetFilePath(".")

	return cfg.Load(true)
}

func GetConfig() Config {
	return cfg
}
