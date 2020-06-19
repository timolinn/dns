package web

import (
	"context"
	"log"
	"net/http"
	"os"
	"syscall"

	"github.com/gorilla/mux"
)

// A Handler handles http requests
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

type App struct {
	*mux.Router
	shutdown chan os.Signal
	mw       []Middleware
}

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

	// wrap app based middlewares
	handler = wrapMiddleware(a.mw, handler)

	h := func(w http.ResponseWriter, r *http.Request) {
		// start tracer span

		if err := handler(r.Context(), w, r); err != nil {
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
