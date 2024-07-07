package handler

import (
	"context"
	"encoding/json"
	"errors"
	"go-ticket-ms/internal/application/usecase"
	"log"
	"net/http"
	"time"
)

type OrderHandler struct {
	createOrderUseCase *usecase.CreateOrderUseCase
}

func NewOrderHandler(createOrderUseCase *usecase.CreateOrderUseCase) *OrderHandler {
	return &OrderHandler{createOrderUseCase: createOrderUseCase}
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	input := usecase.CreateOrderInputDTO{}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func(ctx context.Context) {
		defer cancel()
		err := h.createOrderUseCase.Execute(ctx, input)
		if err != nil {
			log.Printf("Error creating order: %v", err)
		}
	}(ctx)

	select {
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.Canceled) {
			response := map[string]string{"msg": "order created successfully"}
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(response); err != nil {
				http.Error(w, "Error returning response", http.StatusInternalServerError)
			}
		} else {
			http.Error(w, "Request timed out", http.StatusRequestTimeout)
		}
	}
}
