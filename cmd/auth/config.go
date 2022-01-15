package main

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"

	"store/pkg/common/infrastructure/mysql"
)

func parseEnv() (*config, error) {
	c := new(config)
	if err := envconfig.Process(appID, c); err != nil {
		return nil, errors.Wrap(err, "failed to parse env")
	}
	return c, nil
}

type config struct {
	JWTSecret            string `envconfig:"jwt_secret" default:"secret"`
	DBHost               string `envconfig:"db_host"`
	DBName               string `envconfig:"db_name"`
	DBUser               string `envconfig:"db_user"`
	DBPassword           string `envconfig:"db_password"`
	DBMaxConn            int    `envconfig:"db_max_conn" default:"0"`
	DBConnectionLifetime int    `envconfig:"db_conn_lifetime" default:"0"`
	AMQPHost             string `envconfig:"amqp_host"`
	AMQPPort             string `envconfig:"amqp_port"`
	AMQPUser             string `envconfig:"amqp_user" default:"guest"`
	AMQPPassword         string `envconfig:"amqp_password" default:"guest"`
	MigrationsDir        string `envconfig:"migrations_dir"`
}

func (c *config) dsn() mysql.DSN {
	return mysql.DSN{
		Host:     c.DBHost,
		Database: c.DBName,
		User:     c.DBUser,
		Password: c.DBPassword,
	}
}
