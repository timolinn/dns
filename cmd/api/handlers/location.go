package handlers

import (
	"context"
	"net/http"

	"github.com/timolinn/dns/pkg/web"
)

// Locate calculates complex maths
func locate(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// start trace span

	data := CoordsVelocity{}
	err := web.Decode(r, &data)
	if err != nil {
		return web.RespondError(ctx, w, err)
	}
	systemType := System(r.Header.Get("X-System-Type"))

	system := NewSectorNavigator()
	result, err := system.Solve(data, systemType)
	if err != nil {
		er := &web.Error{Err: err, Status: http.StatusBadRequest}
		return web.RespondError(ctx, w, er)
	}
	return web.Respond(ctx, w, system.Response(result, systemType), http.StatusOK)
}
