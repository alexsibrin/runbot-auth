package config

type Config struct {
	PostgreSQL
	Redis
	HttpServer
}

type PostgreSQL struct {
}

type Redis struct {
}

type HttpServer struct {
	Addr string
}

func NewConfig() (*Config, error) {
	return nil, nil
}
