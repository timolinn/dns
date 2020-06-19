package handlers

import (
	"log"
	"net/http"
	"os"

	"github.com/timolinn/dns/middleware"

	"github.com/timolinn/dns/pkg/web"
)

// Register register request handlers and middlewares
func Register(shutdown chan os.Signal, log *log.Logger) http.Handler {
	app := web.NewApp(shutdown, middleware.Logger(log))

	app.MountHandler(http.MethodPost, "/v1/locate", Locate)

	return app
}
