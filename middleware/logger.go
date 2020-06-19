package middleware

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/timolinn/dns/pkg/web"
)

// Logger logs arbitrary data
func Logger(logger *log.Logger) web.Middleware {
	// define actual middleware
	mid := func(f web.Handler) web.Handler {
		// define handler
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			// start span

			err := f(ctx, w, r)

			logger.Printf("%s : (%d) : %s %s -> %s (%s)",
				"unique-trace-id",
				429,
				r.Method, r.URL.Path,
				r.RemoteAddr,
				time.Since(time.Now()), // use time data from context
			)
			return err
		}
		return h
	}
	return mid
}
