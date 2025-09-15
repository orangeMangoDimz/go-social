package main

import (
	"net/http"
	"testing"

	"github.com/orangeMangoDimz/go-social/internal/config"
)

func TestGetUser(t *testing.T) {

	app := newTestApplication(t, config.Config{})

	mux := app.Mount("1.0.0")
	testToken, err := app.Authenticator.GenerateToken(nil)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Should not allow unauthenticated request", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := executeRequest(req, mux)
		checkResponseCode(t, http.StatusUnauthorized, rr.Code)

	})

	t.Run("Should allow authorized request", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Authorization", "Bearer "+testToken)

		rr := executeRequest(req, mux)
		checkResponseCode(t, http.StatusOK, rr.Code)
	})
}
