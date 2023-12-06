package dbpostgres

type Config struct {
}

type PostgreSQL struct {
}

func NewPostgreSQL(c *Config) (*PostgreSQL, error) {
	return &PostgreSQL{}, nil
}
