package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/orangeMangoDimz/go-social/internal/auth"
	"github.com/orangeMangoDimz/go-social/store"
	"github.com/orangeMangoDimz/go-social/store/cache"
	"go.uber.org/zap"
)

func newTestApplication(t *testing.T) *application {
	t.Helper()

	logger := zap.NewNop().Sugar() // Ignore logs
	// logger := zap.Must(zap.NewProduction()).Sugar() // Show logs

	mockStore := store.NewMockStore()
	mockCacheStore := cache.NewMockStore()
	testAuth := &auth.TestAuthenticator{}

	return &application{
		logger:        logger,
		store:         mockStore,
		cacheStorage:  mockCacheStore,
		authenticator: testAuth,
	}

}

func executeRequest(req *http.Request, mux http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d, but we got %d", expected, actual)
	}
}
