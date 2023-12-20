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
	routerrest "runbot-auth/internal/api/rest/v1"
	handlersrest "runbot-auth/internal/api/rest/v1/handlers"
	"runbot-auth/internal/api/rpc"
	handlersrpc "runbot-auth/internal/api/rpc/handlers"
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

	// init REST handlers, middlewares, router
	accounthandlers, err := handlersrest.NewAccount(&handlersrest.DependenciesAccount{
		AccountController: accountcontroller,
	})
	if err != nil {
		log.Fatal(err)
	}

	router, err := routerrest.NewRouter(&routerrest.DependenciesRouter{
		Handlers: &routerrest.Handlers{
			Account: accounthandlers,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// Init http restserver
	restserver, err := rest.NewServer(&rest.DependenciesServer{
		Config: &rest.Config{
			Addr: conf.RestServer.Addr,
		},
		Handler: router,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Init RPC server, handlers
	accountrpchandlers, err := handlersrpc.NewAccount(&handlersrpc.AccountDependencies{
		Controller: accountcontroller,
	})

	grpcserver, err := rpc.NewServer(&rpc.Config{})
	if err != nil {
		log.Fatal(err)
	}
	grpcserver.Add(accountrpchandlers)

	// Init graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		restserver.Run(ctx)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		grpcserver.Run(ctx)
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
