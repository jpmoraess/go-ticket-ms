package domain

import "context"

type OrderRepository interface {
	Save(ctx context.Context, order *Order) error
}
