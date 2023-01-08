package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/bosamatheus/order-book/order/api/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/segmentio/kafka-go"
)

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Logger.SetLevel(log.INFO)

	// Producer
	topic := "orderbook.orders"
	partition := 0
	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:29092", topic, partition)
	if err != nil {
		e.Logger.Fatal("failed to create kafka conn:", err)
	}
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

	// Handler
	orderHandler := handler.NewOrderHandler(conn)

	// Routes
	g := e.Group("api/v1/")
	g.POST("orders", func(c echo.Context) error {
		return orderHandler.SendOrder(c)
	})

	// Start server
	go func() {
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
