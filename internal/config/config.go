package config

import "errors"

const (
	EnvInitKey  = "env"
	YamlInitKey = "yaml"
)

var (
	ErrWrongInitKey = errors.New("wrong config init key")
	ErrConfigInit   = errors.New("error while config initialization")
)

type Config struct {
	Auth
	PostgreSQL
	ClickHouse
	RestServer
}

type PostgreSQL struct {
}

type ClickHouse struct {
}

type RestServer struct {
	Addr string
}

type Auth struct {
	TokenKey string
}

// TODO: Use flags here OR maybe in the main

func NewConfig(key string) (config *Config, err error) {
	switch key {
	case EnvInitKey:
		config, err = initByEnvKey()
	case YamlInitKey:
		config, err = initByYamlKey()
	default:
		config, err = nil, ErrWrongInitKey
	}
	return config, err
}

func initByYamlKey() (*Config, error) {
	return nil, ErrConfigInit
}

func initByEnvKey() (*Config, error) {
	return nil, ErrConfigInit
}
