package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/timolinn/dns/cmd/api/handlers"
)

func main() {
	logger := log.New(os.Stdout, "DNS : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      handlers.Register(shutdown, logger),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     logger,
	}

	// handle errors from request listener
	serverErr := make(chan error, 1)

	// start our server
	go func() {
		log.Println("starting server...")
		serverErr <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErr:
		log.Println(errors.Wrap(err, "server error"))
		os.Exit(1)
	case sig := <-shutdown:
		log.Printf("main %v: Start service shutdown", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			server.Close()
			logger.Println(errors.Wrap(err, "could not stop server gracefully"))
		}
	}
}
