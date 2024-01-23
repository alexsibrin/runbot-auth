package dbpostgres

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
)

var (
	ErrConfigIsNil = errors.New("config is nil")
	ErrDbIsNil     = errors.New("db is nil")
)

const (
	postgresKey = "postgres"
)

type Config struct {
	Db       string
	Host     string
	Port     string
	User     string
	Password string
	SSLMode  bool
}

type PostgreSQL struct {
	db *sql.DB
}

func New(c *Config) (*PostgreSQL, error) {
	// TODO: Add checking
	if c == nil {
		return nil, ErrConfigIsNil
	}

	connevtstring := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", c.User, c.Password, c.Host, c.Port, c.Db)

	db, err := sql.Open(postgresKey, connevtstring)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &PostgreSQL{db}, nil
}
