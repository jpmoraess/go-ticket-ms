package persistence

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"go-ticket-ms/internal/domain"
	"log"
)

type OrderRepository struct {
	pool *pgxpool.Pool
}

func NewOrderRepository(pool *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{pool: pool}
}

func (or *OrderRepository) Save(ctx context.Context, order *domain.Order) error {
	conn, err := or.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		log.Printf("User %d: Error starting transaction: %v\n", order.UserID, err)
		return err
	}

	var available int
	err = tx.QueryRow(ctx, `SELECT available FROM tickets WHERE id = 1 FOR UPDATE`).Scan(&available)
	if err != nil {
		err = tx.Rollback(ctx)
		if err != nil {
			return err
		}
		log.Printf("User %d: Error querying tickets: %v\n", order.UserID, err)
		return err
	}

	if available > 0 && available >= order.Quantity {
		_, err = tx.Exec(ctx, "UPDATE tickets SET available = available - $1 WHERE id = 1", order.Quantity)
		if err != nil {
			err = tx.Rollback(ctx)
			if err != nil {
				return err
			}
			log.Printf("User %d: Error updating tickets: %v\n", order.UserID, err)
			return err
		}

		_, err = tx.Exec(ctx, "INSERT INTO orders (id, user_id, quantity) VALUES($1, $2, $3)", order.ID, order.UserID, order.Quantity)
		if err != nil {
			err = tx.Rollback(ctx)
			if err != nil {
				return err
			}
			log.Printf("User %d: Error insert order into database for user: %v\n", order.UserID, err)
			return err
		}

		err = tx.Commit(ctx)
		if err != nil {
			log.Printf("User %d: Error commiting transaction: %v\n", order.UserID, err)
			return err
		}
	} else {
		err = tx.Rollback(ctx)
		if err != nil {
			return err
		}
		log.Printf("User %d found no tickets.\n", order.UserID)
	}
	return nil
}
