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
