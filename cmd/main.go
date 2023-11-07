package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"runbot-auth/internal/api/rest/v1"
	"runbot-auth/internal/config"
	"sync"
	"syscall"
)

func main() {

	conf, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	server, err := rest.NewHttpServer()
	if err != nil {
		log.Fatal(err)
	}

	erchan := make(chan error)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		server.Run(ctx)
	}()

}
