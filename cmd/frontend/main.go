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

	// Shutdown the app on signal
	ctx := context.Background()
	// Listen for SIGINT to gracefully shutdown.
	nctx, stop := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer stop()

	// Graceful shutdown
	go func() {
		<-nctx.Done()
		log.Println("shutdown initiated")
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		server.Shutdown(ctx)
		log.Println("shutdown")
	}()

	log.Printf("starting http server on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
