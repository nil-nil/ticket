package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/nil-nil/ticket/internal/frontend"
)

func main() {
	server := frontend.NewServer()

	// Graceful shutdown
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)
		<-quit
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	log.Printf("starting http server on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
