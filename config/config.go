package config

import (
	"github.com/caarlos0/env/v6"
)

const (
	EnvLocal string = "local"
)

type Conf struct {
	AppEnv string `env:"APP_ENV" envDefault:"local"`

	HttpPort string `env:"HTTP_PORT" envDefault:"8080"`

	PgHost     string `env:"DB_HOST"`
	PgPort     string `env:"DB_PORT"`
	PgUser     string `env:"DB_USERNAME"`
	PgPassword string `env:"DB_PASSWORD"`
	PgDbName   string `env:"DB_NAME"`

	MaxConnections int `env:"MAX_CONNECTIONS" envDefault:"100"`
}

var Cnf Conf

func NewConf() error {
	if err := env.Parse(&Cnf, env.Options{RequiredIfNoDef: true}); err != nil {
		return err
	}

	return nil
}
