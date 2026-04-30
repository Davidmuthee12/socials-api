package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Davidmuthee12/socials/internal/auth"
	ratelimiter "github.com/Davidmuthee12/socials/internal/rateLimiter"
	"github.com/Davidmuthee12/socials/internal/store"
	cache "github.com/Davidmuthee12/socials/internal/store/cache"
	"go.uber.org/zap"
)

func newTestApplication(t *testing.T, cfg config) *application {
	t.Helper()

	logger := zap.NewNop().Sugar()
	mockStore := store.NewMockStore()
	mockCacheStore := cache.NewMockStore()

	testAuth := &auth.TestAuthenticator{}

	if cfg.rateLimiter.RequestPerTimeFrame <= 0 {
		cfg.rateLimiter.RequestPerTimeFrame = 20
	}

	if cfg.rateLimiter.TimeFrame <= 0 {
		cfg.rateLimiter.TimeFrame = 5 * time.Second
	}

	limiter := ratelimiter.NewFixedWindowLimiter(
		cfg.rateLimiter.RequestPerTimeFrame,
		cfg.rateLimiter.TimeFrame,
	)

	return &application{
		logger:        logger,
		config:        cfg,
		store:         mockStore,
		cacheStorage:  mockCacheStore,
		authenticator: testAuth,
		rateLimiter:   limiter,
	}
}

func executeRequest(req *http.Request, mux http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("expected respinse code %d but got %d", expected, actual)
	}
}
