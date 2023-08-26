package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nil-nil/ticket/internal/services/api"
	"github.com/nil-nil/ticket/internal/services/config"
)

func main() {
	configFilePath := flag.String("config", "config.yaml", "Configuration file")
	flag.Parse()

	config, err := config.ReadAndParseConfigFile(*configFilePath)
	if err != nil {
		log.Fatal(err)
	}

	apiServer := api.NewApi()
	authProvider, err := api.NewJwtAuthProvider(
		func(userID uint64) (user api.User, err error) {
			return api.User{Id: 999}, nil
		},
		[]byte(config.Auth.JWT.PublicKey),
		[]byte(config.Auth.JWT.PrivateKey),
		api.GetJWTProtocol(config.Auth.JWT.SigningMethod),
		config.Auth.JWT.TokenLifetime,
	)
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Use(middleware.Logger())
	e.Use(middleware.Gzip())
	api.RegisterHandlers(e, api.NewStrictHandler(apiServer, []runtime.StrictEchoMiddlewareFunc{
		api.AuthMiddleware(authProvider),
	}))

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
		e.Shutdown(ctx)
		log.Println("shutdown")
	}()

	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%d", config.HTTP.ListenAddress, config.HTTP.Port)))
}
