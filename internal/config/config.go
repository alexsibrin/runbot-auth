package config

import (
	"errors"
	"github.com/spf13/viper"
)

const (
	EnvInitKey = "env"

	YamlInitKey = "yaml"
	configPath  = "."
	configName  = "config"
)

var (
	ErrWrongInitKey = errors.New("wrong config init key")
	ErrConfigInit   = errors.New("error while config initialization")
)

type Config struct {
	PostgreSQL
	RestServer
	Logger
}

type PostgreSQL struct {
	Db       string
	Host     string
	Port     string
	User     string
	Password string
	SSLMode  bool
}

type RestServer struct {
	Host string
	Port string
}

type Logger struct {
	Level         int
	Colors        bool
	FullTimestamp bool
}

// TODO: Use flags here OR maybe in the main

func New(key string) (config *Config, err error) {
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
	viper.AddConfigPath(configPath)
	viper.SetConfigName(configName)

	var c Config

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func initByEnvKey() (*Config, error) {
	// TODO: Complete me
	return nil, ErrConfigInit
}
