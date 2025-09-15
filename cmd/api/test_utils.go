package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/orangeMangoDimz/go-social/internal/auth"
	"github.com/orangeMangoDimz/go-social/internal/config"
	"github.com/orangeMangoDimz/go-social/internal/ratelimiter"
	httpserver "github.com/orangeMangoDimz/go-social/internal/server/http"
	"github.com/orangeMangoDimz/go-social/internal/storage/cache"
	"github.com/orangeMangoDimz/go-social/internal/storage/postgres"
	"go.uber.org/zap"
)

func newTestApplication(t *testing.T, cfg config.Config) *httpserver.Application {
	t.Helper()

	logger := zap.NewNop().Sugar()
	// Uncomment to enable logs
	// logger := zap.Must(zap.NewProduction()).Sugar()
	mockStore := postgres.NewMockStore()
	mockCacheStore := cache.NewMockStore()

	testAuth := &auth.TestAuthenticator{}

	// Rate limiter
	rateLimiter := ratelimiter.NewFixedWindowLimiter(
		cfg.RateLimiter.RequestPerTimeFrame,
		cfg.RateLimiter.TimeFrame,
	)

	return &httpserver.Application{
		Logger:        logger,
		Store:         mockStore,
		CacheStorage:  mockCacheStore,
		Authenticator: testAuth,
		Config:        cfg,
		RateLimiter:   rateLimiter,
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
