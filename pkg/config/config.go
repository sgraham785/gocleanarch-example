package config

import (
	"github.com/kelseyhightower/envconfig"
)

// PostgresConf is the specification for PostgreSQL configs
type PostgresConf struct {
	PostgresHost     string `required:"true" split_words:"true"`
	PostgresDB       string `required:"true" split_words:"true"`
	PostgresPort     int    `default:"5432" split_words:"true"`
	PostgresUser     string `required:"true" split_words:"true"`
	PostgresPassword string `required:"true" split_words:"true"`
}

// Specification struct for environment variables
type Specification struct {
	PostgresConf          `desc:"PostgreSQL config"`
	PrometheusPushgateway string `required:"true" split_words:"true"`
	APIPort               int    `default:"9000" split_words:"true"`
}

// Load is what loads the config.
func Load() *Specification {
	var cfg Specification
	err := envconfig.Process("gcarch", &cfg)
	if err != nil {
		panic(err)
	}
	return &cfg
}
