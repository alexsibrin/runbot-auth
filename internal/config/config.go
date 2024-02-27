package config

import (
	"errors"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
	"time"
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
	Host string
	Port string
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

/*func New(key string) (config *Config, err error) {
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
*/

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

// --- Some old stuff
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

// FIXME: Waiting when viper will get ability to map env vars to a struct directly without kostils
func initByEnvKey() (*Config, error) {
	var cfg Config
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	for _, s := range os.Environ() {
		replacer := strings.NewReplacer("_", ".")
		envToBind := replacer.Replace(strings.Split(s, "=")[0])
		if err := viper.BindEnv(envToBind); err != nil {
			log.Fatal(err)
		}
	}

	err := viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
