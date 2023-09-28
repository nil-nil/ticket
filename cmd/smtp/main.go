package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/nil-nil/ticket/internal/domain"
	"github.com/nil-nil/ticket/internal/infrastructure/gosmtpmail"
	"github.com/nil-nil/ticket/internal/infrastructure/ristrettocache"
	"github.com/nil-nil/ticket/internal/infrastructure/ticketeventbus"
)

func main() {
	// configFilePath := flag.String("config", "config.yaml", "Configuration file")
	// flag.Parse()
	// config, err := config.ReadAndParseConfigFile(*configFilePath)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	cache, err := ristrettocache.NewCache(nil)
	if err != nil {
		log.Fatal(err)
	}

	bus, err := ticketeventbus.NewBus(":")
	if err != nil {
		log.Fatal(err)
	}

	server := gosmtpmail.NewServer(nil, cache, bus, func(username, password string) (domain.User, error) { return domain.User{}, nil })

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
