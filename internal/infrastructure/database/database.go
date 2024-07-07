package database

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
)

const (
	dsn                  = "postgres://backend:backend@localhost:5432/postgres?sslmode=disable"
	initTicketsAvailable = 500
)

func InitDB() (*pgxpool.Pool, error) {
	pool, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	// create tickets table
	_, err = pool.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS tickets (id SERIAL PRIMARY KEY, available int)`)
	if err != nil {
		return nil, err
	}

	// Adding UUID EXTENSION
	_, err = pool.Exec(context.Background(), `CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)
	if err != nil {
		return nil, err
	}

	// create order table
	_, err = pool.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS orders (id uuid PRIMARY KEY, user_id INT, quantity INT)`)
	if err != nil {
		return nil, err
	}

	// start tickets table with available tickets
	_, err = pool.Exec(context.Background(), `INSERT INTO tickets (id, available) VALUES (1, $1) ON CONFLICT (id) DO UPDATE SET available = EXCLUDED.available`, initTicketsAvailable)
	if err != nil {
		return nil, err
	}

	return pool, nil
}
