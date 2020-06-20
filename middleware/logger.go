package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/timolinn/dns/pkg/web"
)

// Logger logs arbitrary data
func Logger(logger *log.Logger) web.Middleware {
	// actual middleware
	mid := func(f web.Handler) web.Handler {
		// define web Handler
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			// start tracer span

			v, ok := ctx.Value(web.KeyValues).(*web.Values)
			if !ok {
				return errors.New("web value missing from context")
			}

			err := f(ctx, w, r)

			logger.Printf("%s : (%d) : %s %s -> %s (%s)",
				"unique-trace-id", // TODO: use uuids for unique ID for tracing
				v.StatusCode,
				r.Method, r.URL.Path,
				r.RemoteAddr,
				time.Since(v.Now),
			)
			return err
		}
		return h
	}
	return mid
}
