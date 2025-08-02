package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/cors"

	middleware "concurrent-order-processing-system/middlewares"
	"concurrent-order-processing-system/routes"
)

func main() {
	router := http.NewServeMux()
	routes.InitializeRoutes()

	chainedHandler := middleware.Chain(
		router,
		middleware.LoggingMiddleware,
	)

	handler := cors.AllowAll().Handler(chainedHandler)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	stopChannel := make(chan os.Signal, 1)
	signal.Notify(stopChannel, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		fmt.Printf("Server running on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Could not start server: %v", err)
		}
	}()

	<-stopChannel

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Println("Shutting down server...")
	if err := server.Shutdown(shutdownCtx); err != nil {
		fmt.Printf("Server forced to shutdown: %v", err)
	} else {
		fmt.Println("Server gracefully stopped.")
	}
}
