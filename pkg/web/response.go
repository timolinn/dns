package web

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

// Respond sends successful request processing response to the client
func Respond(ctx context.Context, w http.ResponseWriter, data interface{}, statusCode int) error {
	// start tracer span

	// If the context is missing this value, request the service
	// will be shutdown gracefully.
	v, ok := ctx.Value(KeyValues).(*Values)
	if !ok {
		return errors.New("web value missing from context")
	}

	// Set the statusCode for the http request logger middleware.
	v.StatusCode = statusCode

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// Sends the result back to the client.
	if _, err := w.Write(jsonData); err != nil {
		return err
	}
	return nil
}

// RespondError sends error response to the client
func RespondError(ctx context.Context, w http.ResponseWriter, err error) error {
	// check if err is a web error
	if webErr, ok := (err).(*Error); ok {
		er := ErrorResponse{
			Error:  webErr.Err.Error(),
			Fields: webErr.Fields,
		}
		if err := Respond(ctx, w, er, webErr.Status); err != nil {
			return err
		}
		return nil
	}

	er := ErrorResponse{
		Error: http.StatusText(http.StatusInternalServerError),
	}

	if err := Respond(ctx, w, er, http.StatusInternalServerError); err != nil {
		return err
	}
	return nil
}
