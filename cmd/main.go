package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runbot-auth/internal/api/rest"
	v1 "runbot-auth/internal/api/rest/v1/handlers"
	"runbot-auth/internal/config"
	"runbot-auth/internal/services"
	"sync"
	"syscall"
)

func main() {
	log.Println("App is starting the initialization...")
	// Init config
	conf, err := config.NewConfig(config.EnvInitKey)
	if err != nil {
		log.Fatal(err)
	}

	// <- init logger
	defer log.Println("App is stopped.")

	// <- init db, cache

	// <- init services
	authservice, err := services.NewAuth(&services.DependenciesAuth{})

	// <-- init handlers, middlewares, router
	authhandler, err := v1.NewAuth(&v1.DependenciesAuth{})

	// Init http server
	server, err := rest.NewHttpServer(&rest.DependenciesHttpServer{
		Config: conf.HttpServer,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Init graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		server.Run(ctx)
		wg.Done()
	}()

	log.Println("App is running.")
	<-ctx.Done()

	if !errors.Is(ctx.Err(), context.Canceled) {
		err := context.Cause(ctx)
		log.Fatal(fmt.Sprintf("App is crushed: %s", err.Error()))
	}

	log.Println("Services are stopping. Please wait...")

	wg.Wait()
}
