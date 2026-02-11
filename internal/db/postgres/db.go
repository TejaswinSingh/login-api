package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DbConfig struct {
	DbUser     string
	DbPassword string
	DbHost     string
	DbPort     int
	DbName     string
	DbSslMode  string
}

func NewDbConnPool(config DbConfig) (*pgxpool.Pool, error) {
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		config.DbUser,
		config.DbPassword,
		config.DbHost,
		config.DbPort,
		config.DbName,
		config.DbSslMode,
	)
	dbpool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		return nil, err
	}
	return dbpool, nil
}
