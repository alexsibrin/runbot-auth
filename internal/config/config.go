package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
	"time"
)

const (
	configPath = "."
	configName = "config"
)

type Config struct {
	PostgreSQL
	RestServer
	GRPCServer
	Logger
	Jwt
	Common
}

type PostgreSQL struct {
	Db       string
	Host     string
	Port     string
	User     string
	Password string
	SSLMode  string
}

type RestServer struct {
	Host              string
	Port              string
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
}

type GRPCServer struct {
	Port int
}

type Logger struct {
	Level         int
	Colors        bool
	FullTimestamp bool
}

type Jwt struct {
	Salt      string
	Issuer    string
	Subject   string
	Audience  []string
	ExpiresIn time.Duration
}

type Common struct {
	Version string
	Health  string
}

func New() (*Config, error) {
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

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	for _, s := range os.Environ() {
		replacer := strings.NewReplacer("_", ".")
		envToBind := replacer.Replace(strings.Split(s, "=")[0])
		if err := viper.BindEnv(envToBind); err != nil {
			log.Fatal(err)
		}
	}

	err = viper.Unmarshal(&c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
