package usecase

import (
	"context"
	"go-ticket-ms/internal/domain"
	"log"
)

type CreateOrderInputDTO struct {
	UserID   int `json:"userId"`
	Quantity int `json:"quantity"`
}

type CreateOrderUseCase struct {
	orderRepo domain.OrderRepository
}

func NewCreateOrderUseCase(orderRepo domain.OrderRepository) *CreateOrderUseCase {
	return &CreateOrderUseCase{orderRepo: orderRepo}
}

func (uc *CreateOrderUseCase) Execute(ctx context.Context, input CreateOrderInputDTO) error {
	order := domain.NewOrder(input.UserID, input.Quantity)
	if err := uc.orderRepo.Save(ctx, order); err != nil {
		log.Printf("UserId %d: error creating order: %s\n", order.UserID, err.Error())
		return err
	}
	return nil
}
