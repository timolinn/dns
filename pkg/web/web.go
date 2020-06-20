// Package web provides functions and concrete types
// for configuring and instrumenting out web application
package web

import (
	"context"
	"log"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

type key int

// KeyValues is how request values or stored/retrieved.
const KeyValues key = 1

// Values represent state for each request.
type Values struct {
	Now        time.Time
	StatusCode int
}

// A Handler handles http requests
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// App bootstraps the application, it's the entrypoint to our service
type App struct {
	*mux.Router
	shutdown chan os.Signal
	mw       []Middleware
}

// NewApp constructs an App
func NewApp(shutdown chan os.Signal, mw ...Middleware) *App {
	app := &App{
		Router:   mux.NewRouter(),
		shutdown: shutdown,
		mw:       mw,
	}

	return app
}

// MountHandler mounts a http handler on the router
func (a *App) MountHandler(verb, path string, handler Handler, mw ...Middleware) {
	handler = wrapMiddleware(mw, handler)

	// wrap application level middlewares
	handler = wrapMiddleware(a.mw, handler)

	h := func(w http.ResponseWriter, r *http.Request) {
		// TODO: start tracer span

		// add relevant values the context for propagation
		v := Values{
			Now: time.Now(),
		}
		ctx := context.WithValue(r.Context(), KeyValues, &v)

		if err := handler(ctx, w, r); err != nil {
			log.Println(err)
			a.Shutdown()
			return
		}
	}
	a.HandleFunc(path, h).Methods(verb)
}

// Shutdown sends a sigterm signal to the app to shutdown gracefully
func (a *App) Shutdown() {
	a.shutdown <- syscall.SIGTERM
}
