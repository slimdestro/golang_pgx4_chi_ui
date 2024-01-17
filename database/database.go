// Upwork test task. will be deleted soon...

package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

func InitDB() *pgxpool.Pool {
	connectionString := "postgresql://postgres:admin@localhost/postgres?sslmode=disable"

	config, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		log.Fatal(err)
	}

	db, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Fatal(err)
	}

	return db
}
