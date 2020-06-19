package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/timolinn/dns/cmd/api/handlers"
)

var addr string
var readtimeout, writetimeout time.Duration

func main() {
	flag.StringVar(&addr, "host", ":8080", "define server address")
	flag.DurationVar(&readtimeout, "read timeout", 5, "sets the read timeout in seconds")
	flag.DurationVar(&writetimeout, "write timeout", 10, "sets the write timeout in seconds")

	logger := log.New(os.Stdout, "DNS : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	if err := run(logger); err != nil {
		logger.Println(err)
		os.Exit(1)
	}
}

func run(logger *log.Logger) error {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	server := &http.Server{
		Addr:         addr,
		Handler:      handlers.Register(shutdown, logger),
		ReadTimeout:  readtimeout * time.Second,
		WriteTimeout: writetimeout * time.Second,
		ErrorLog:     logger,
	}

	// handle errors from request listener
	serverErr := make(chan error, 1)

	// start our server
	go func() {
		log.Println("started running")
		serverErr <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErr:
		return errors.Wrap(err, "server error")

	case sig := <-shutdown:
		log.Printf("main %v: Start service shutdown", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			server.Close()
			return errors.Wrap(err, "could not stop server gracefully")
		}
	}
	return nil
}
