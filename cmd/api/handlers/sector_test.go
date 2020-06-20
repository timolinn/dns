package handlers_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/timolinn/dns/cmd/api/handlers"
)

func TestSolve(t *testing.T) {
	cases := []struct {
		in         handlers.CoordsVelocity
		out        float64
		systemType handlers.System
	}{
		{handlers.CoordsVelocity{2.0, 2.0, 2.0, 2.0}, 8, handlers.Drone},
		{handlers.CoordsVelocity{21.4, 20.3, 223.5, 444.0}, 709.2, handlers.Ship},
		{handlers.CoordsVelocity{21.44, 20.43, 223.75, 444.09}, 709.71, handlers.Drone},
	}

	navigator := handlers.NewSectorNavigator()
	for _, test := range cases {
		t.Run("should return correct computation result", func(t *testing.T) {
			res, err := navigator.Solve(test.in, handlers.Drone)
			if err != nil {
				t.Fatalf("expected nil-err got %s", err)
			}
			if res != test.out {
				t.Errorf("SectorNavigator.Solve(): want=%v :: got=%v", test.out, res)
			}
		})
	}
}

func TestResponse(t *testing.T) {
	cases := []struct {
		in         float64
		out        map[string]float64
		systemType handlers.System
	}{
		{138.89, map[string]float64{"loc": 138.89}, handlers.Drone},
		{138.89, map[string]float64{"location": 138.89}, handlers.Ship},
	}

	navigator := handlers.NewSectorNavigator()
	for _, test := range cases {
		t.Run("should return correct response structure", func(t *testing.T) {
			res := navigator.Response(test.in, test.systemType)
			if !reflect.DeepEqual(res, test.out) {
				t.Errorf("SectorNavigator.Solve(): want=%v :: got=%v", test.out, res)
			}
		})
	}
}

func TestHome(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	shutdown := make(chan os.Signal, 1)
	logger := log.New(os.Stdout, "TEST : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	app := handlers.Register(shutdown, logger)
	app.ServeHTTP(w, r)

	result := w.Result()
	if result.StatusCode != http.StatusOK {
		t.Errorf("should return correct status code: want=%v, got=%v", http.StatusOK, result.StatusCode)
	}
}
