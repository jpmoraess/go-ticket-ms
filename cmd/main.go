package main

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4/pgxpool"
	"go-ticket-ms/internal/application/usecase"
	"go-ticket-ms/internal/infrastructure/handler"
	"go-ticket-ms/internal/infrastructure/persistence"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	connStr = "postgres://backend:backend@localhost:5432/postgres?sslmode=disable"
)

func main() {
	pool, err := pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()

	// init repositories
	orderRepo := persistence.NewOrderRepository(pool)

	// init use cases
	createOrderUseCase := usecase.NewCreateOrderUseCase(orderRepo)

	// handler
	orderHandler := handler.NewOrderHandler(createOrderUseCase)

	r := http.NewServeMux()
	r.HandleFunc("POST /order", orderHandler.CreateOrder)

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Channel for listening to operating system signals
	idleConnClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
		<-sigint

		// Interrupt signal received
		log.Println("Interrupt signal received, starting graceful shutdown...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Graceful shutdown error: %v\n", err)
		}
		close(idleConnClosed)
	}()

	// Starting HTTP server
	log.Println("HTTP server running on port 8080")
	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Erro ao iniciar o servidor HTTP: %v\n", err)
	}

	<-idleConnClosed
	log.Println("Finalized HTTP server successfully")
}
