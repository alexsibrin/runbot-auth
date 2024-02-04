package dbpostgres

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"time"
)

var (
	ErrConfigIsNil = errors.New("config is nil")
	ErrDbIsNil     = errors.New("db is nil")
)

const (
	postgresKey = "postgres"
)

type Config struct {
	Db                    string
	Host                  string
	Port                  string
	User                  string
	Password              string
	SSLMode               string
	MaxOpenConnections    int
	MaxIdleConnections    int
	ConnectionMaxLifetime time.Duration
	ConnectionMaxIdletime time.Duration
}

type PostgreSQL struct {
	db *sql.DB
}

func New(c *Config) (*PostgreSQL, error) {
	if c == nil {
		return nil, ErrConfigIsNil
	}

	connevtstring := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", c.User, c.Password, c.Host, c.Port, c.Db, c.SSLMode)

	db, err := sql.Open(postgresKey, connevtstring)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(c.MaxOpenConnections)
	db.SetMaxIdleConns(c.MaxIdleConnections)
	db.SetConnMaxLifetime(c.ConnectionMaxLifetime)
	db.SetConnMaxIdleTime(c.ConnectionMaxIdletime)

	return &PostgreSQL{db}, nil
}

func (p *PostgreSQL) Close() error {
	if p.db != nil {
		return p.db.Close()
	}
	return ErrDbIsNil
}
