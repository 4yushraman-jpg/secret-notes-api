package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRateLimiter(t *testing.T) {
	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	protectedHandler := RateLimitMiddleware(dummyHandler)

	for i := 1; i <= 5; i++ {
		req := httptest.NewRequest("GET", "http://localhost:8080/notes", nil)
		req.RemoteAddr = "192.168.1.1:1234"
		w := httptest.NewRecorder()

		protectedHandler.ServeHTTP(w, req)

		if i <= 3 {
			if w.Code != http.StatusOK {
				t.Errorf("Request %d should have passed, but got %d", i, w.Code)
			}
		} else {
			if w.Code != http.StatusTooManyRequests {
				t.Errorf("Request %d should have failed with 429, but got %d", i, w.Code)
			}
		}
	}
}
