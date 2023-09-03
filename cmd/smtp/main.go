package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	mailgosmtp "github.com/nil-nil/ticket/internal/infrastructure/mail-go-smtp"
	"github.com/nil-nil/ticket/internal/services/config"
)

func main() {
	configFilePath := flag.String("config", "config.yaml", "Configuration file")
	flag.Parse()

	config, err := config.ReadAndParseConfigFile(*configFilePath)
	if err != nil {
		log.Fatal(err)
	}

	server := mailgosmtp.NewServer()

	// Shutdown the app on signal
	ctx := context.Background()
	// Listen for SIGINT to gracefully shutdown.
	nctx, stop := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer stop()

	go func() {
		<-nctx.Done()
		log.Println("shutdown initiated")
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		server.Shutdown(ctx)
		log.Println("shutdown")
	}()

	server.ListenAndServe()
}
