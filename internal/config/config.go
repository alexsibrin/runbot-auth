package config

import "errors"

const (
	EnvInitKey  = "env"
	YamlInitKey = "yaml"
)

var (
	ErrWrongInitKey = errors.New("wrong config init key")
)

type Config struct {
	*PostgreSQL
	*Redis
	*HttpServer
}

type PostgreSQL struct {
}

type Redis struct {
}

type HttpServer struct {
	Addr string
}

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
	return nil, nil
}

func initByEnvKey() (*Config, error) {
	return nil, nil
}
