package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runbot-auth/internal/api/controllers"
	"runbot-auth/internal/api/rest"
	routerv1 "runbot-auth/internal/api/rest/v1"
	handlersv1 "runbot-auth/internal/api/rest/v1/handlers"
	"runbot-auth/internal/api/rpc"
	"runbot-auth/internal/config"
	"runbot-auth/internal/usecases"
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

	// init logger
	defer log.Println("App is stopped.")

	// init db, cache

	// init usecases
	accountusecase, err := usecases.NewAccount(&usecases.AccountDependencies{})

	// init controllers
	accountcontroller, err := controllers.NewAccount(&controllers.AccountDependencies{
		Usecase: accountusecase,
	})

	// init handlers, middlewares, router
	accounthandlers, err := handlersv1.NewAccount(&handlersv1.DependenciesAccount{
		AccountController: accountcontroller,
	})
	if err != nil {
		log.Fatal(err)
	}

	router, err := routerv1.NewRouter(&routerv1.DependenciesRouter{
		Handlers: &routerv1.Handlers{
			Account: accounthandlers,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// Init grpcserver
	grpcserver, err := rpc.NewServer(&rpc.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Init http restserver
	restserver, err := rest.NewServer(&rest.DependenciesServer{
		Config: &rest.Config{
			Addr: conf.Server.Addr,
		},
		Handler: router,
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
		restserver.Run(ctx)
		wg.Done()
	}()

	log.Println("App is running.")
	ctx.Done()

	if !errors.Is(ctx.Err(), context.Canceled) {
		err := context.Cause(ctx)
		log.Fatal(fmt.Sprintf("App is crushed: %s", err.Error()))
	}

	log.Println("Services are stopping. Please wait...")

	wg.Wait()
}
