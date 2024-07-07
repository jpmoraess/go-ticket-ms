package domain

import "github.com/google/uuid"

type Order struct {
	ID       uuid.UUID
	UserID   int
	Quantity int
}

func NewOrder(userID int, quantity int) *Order {
	return &Order{
		ID:       uuid.New(),
		UserID:   userID,
		Quantity: quantity,
	}
}
