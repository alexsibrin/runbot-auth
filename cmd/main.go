package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/alexsibrin/runbot-auth/internal/api/controllers"
	"github.com/alexsibrin/runbot-auth/internal/api/rest"
	restv1 "github.com/alexsibrin/runbot-auth/internal/api/rest/v1"
	handlersrest "github.com/alexsibrin/runbot-auth/internal/api/rest/v1/handlers"
	"github.com/alexsibrin/runbot-auth/internal/api/rpc"
	handlersrpc "github.com/alexsibrin/runbot-auth/internal/api/rpc/handlers"
	"github.com/alexsibrin/runbot-auth/internal/config"
	"github.com/alexsibrin/runbot-auth/internal/usecases"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	log.Println("App is starting the initialization...")

	// Init config
	conf, err := config.New(config.YamlInitKey)
	if err != nil {
		log.Fatal(err)
	}

	// init logger
	defer log.Println("App is stopped.")

	// init db, cache

	// init usecases
	accountusecase, err := usecases.NewAccount(&usecases.AccountDependencies{})
	if err != nil {
		log.Fatal(err)
	}

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

	router, err := restv1.NewRouter(&restv1.DependenciesRouter{
		Handlers: &restv1.Handlers{
			Account: accounthandlers,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// Init http restserver
	restserver, err := rest.NewServer(&rest.DependenciesServer{
		Config: &rest.Config{
			Host: conf.RestServer.Host,
			Port: conf.RestServer.Port,
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
