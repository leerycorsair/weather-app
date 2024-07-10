package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	UsersTable       = "users"
	UsersCitiesTable = "users_cities"
	CitiesTable      = "cities"
	ForecastsTable   = "forecasts"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

type PGConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPgConnection(conn PGConfig) (*sqlx.DB, error) {
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		conn.Username,
		conn.Password,
		conn.Host,
		conn.Port,
		conn.DBName,
		conn.SSLMode,
	)

	db, err := sqlx.Open("postgres", connString)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
