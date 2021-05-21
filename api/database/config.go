package database

import "github.com/jenpet/plooral/config"

type dbConfig struct {
	PostgresURI string `required:"true" envconfig:"POSTGRES_URI"`
}

func parseConfig() dbConfig {
	cfg := dbConfig{}
	config.Parse(&cfg)
	return cfg
}
