package handlers_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/timolinn/dns/cmd/api/handlers"
	"github.com/timolinn/dns/pkg/web"
)

func TestLocation(t *testing.T) {
	var payload = []byte(`{"x":"123.12","z":"789.89","y":"456.56", "vel":"20.0"}`)
	var badPayload = []byte(`{"x":"123,12","z":"789.89","y":"456.56", "vel":"20.0"}`)
	var incompletePayload = []byte(`{"x":"123.12","z":"789.89"}`)
	var success = []byte(`{"loc": 1389.57}`)
	var successShip = []byte(`{"location": 1389.57}`)

	shutdown := make(chan os.Signal, 1)
	logger := log.New(os.Stdout, "TEST : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	t.Run("should return status code 200 OK", func(t *testing.T) {
		buf := bytes.NewReader(payload)
		r := httptest.NewRequest(http.MethodPost, "/v1/locate", buf)
		r.Header.Set("X-System-Type", "drone")
		w := httptest.NewRecorder()

		app := handlers.Register(shutdown, logger)
		app.ServeHTTP(w, r)

		result := w.Result()
		if result.StatusCode != http.StatusOK {
			t.Errorf("Should receive status code %d, got %d", http.StatusOK, result.StatusCode)
		}
	})

	t.Run("should return status code fail for bad payload", func(t *testing.T) {
		buf := bytes.NewReader(badPayload)
		r := httptest.NewRequest(http.MethodPost, "/v1/locate", buf)
		r.Header.Set("X-System-Type", "drone")
		w := httptest.NewRecorder()

		app := handlers.Register(shutdown, logger)
		app.ServeHTTP(w, r)

		result := w.Result()
		if result.StatusCode != http.StatusBadRequest {
			t.Errorf("Should receive status code %d, got %d", http.StatusOK, result.StatusCode)
		}
		type got struct {
			Error  string `json:"error"`
			Fields string `json:"fields"`
		}
		g := got{}
		if err := json.NewDecoder(result.Body).Decode(&g); err != nil {
			t.Errorf("should be able to unmarshal response")
		}

		if g.Error != "malformed request data" {
			t.Errorf("want %s, got %s", "malformed request data", g.Error)
		}
	})

	t.Run("should report validation error for incomplete payload", func(t *testing.T) {
		buf := bytes.NewReader(incompletePayload)
		r := httptest.NewRequest(http.MethodPost, "/v1/locate", buf)
		r.Header.Set("X-System-Type", "drone")
		w := httptest.NewRecorder()

		app := handlers.Register(shutdown, logger)
		app.ServeHTTP(w, r)

		result := w.Result()
		if result.StatusCode != http.StatusUnprocessableEntity {
			t.Errorf("Should receive status code %d, got %d", http.StatusUnprocessableEntity, result.StatusCode)
		}
		type got struct {
			Error  string           `json:"error"`
			Fields []web.FieldError `json:"fields"`
		}
		g := got{}
		if err := json.NewDecoder(result.Body).Decode(&g); err != nil {
			t.Errorf("should be able to unmarshal response")
		}

		if g.Error != "validation error" {
			t.Errorf("want %s, got %s", "validation error", g.Error)
		}

		if g.Fields == nil {
			t.Errorf("missing fields should specified")
		}
	})

	t.Run("should fail when no systemType header is specified", func(t *testing.T) {
		buf := bytes.NewReader(payload)
		r := httptest.NewRequest(http.MethodPost, "/v1/locate", buf)
		w := httptest.NewRecorder()

		app := handlers.Register(shutdown, logger)
		app.ServeHTTP(w, r)

		result := w.Result()
		if result.StatusCode != http.StatusBadRequest {
			t.Errorf("Should receive status code %d, got %d", http.StatusBadRequest, result.StatusCode)
		}
	})

	t.Run("should pass for all supported systemType", func(t *testing.T) {
		st := []string{"drone", "ship"}
		buf := bytes.NewReader(payload)
		r := httptest.NewRequest(http.MethodPost, "/v1/locate", buf)
		w := httptest.NewRecorder()
		app := handlers.Register(shutdown, logger)

		for _, s := range st {
			r.Header.Set("X-System-Type", s)
			app.ServeHTTP(w, r)

			result := w.Result()
			if result.StatusCode != http.StatusOK {
				t.Errorf("Should receive status code %d, got %d", http.StatusOK, result.StatusCode)
			}
		}
	})

	t.Run("should return 'loc' for drone systemType", func(t *testing.T) {
		buf := bytes.NewReader(payload)
		r := httptest.NewRequest(http.MethodPost, "/v1/locate", buf)
		w := httptest.NewRecorder()
		app := handlers.Register(shutdown, logger)

		r.Header.Set("X-System-Type", "drone")
		app.ServeHTTP(w, r)

		result := w.Result()
		type response struct {
			Loc float64 `json:"loc"`
		}
		wt := response{}
		gt := response{}
		json.NewDecoder(bytes.NewReader(success)).Decode(&wt)
		if err := json.NewDecoder(result.Body).Decode(&gt); err != nil {
			t.Errorf("should be able to unmarshal response")
		}

		if gt.Loc != wt.Loc {
			t.Errorf("want %v, got %v", wt.Loc, gt.Loc)
		}
	})

	t.Run("should return 'location' for ship systemType", func(t *testing.T) {
		buf := bytes.NewReader(payload)
		r := httptest.NewRequest(http.MethodPost, "/v1/locate", buf)
		w := httptest.NewRecorder()
		app := handlers.Register(shutdown, logger)

		r.Header.Set("X-System-Type", "ship")
		app.ServeHTTP(w, r)

		result := w.Result()
		type response struct {
			Location float64 `json:"location"`
		}
		wt := response{}
		gt := response{}
		json.NewDecoder(bytes.NewReader(successShip)).Decode(&wt)
		if err := json.NewDecoder(result.Body).Decode(&gt); err != nil {
			t.Errorf("should be able to unmarshal response")
		}

		if gt.Location != wt.Location {
			t.Errorf("want %v, got %v", wt.Location, gt.Location)
		}
	})
}
