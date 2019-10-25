package httputils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func generateTestRun(t *testing.T, restrict bool, login string, password string, expectedStatus int) {
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	})
	middleware := NewBasicAuthMiddleware(map[string]string{
		"user": "$2y$12$AH6cua0416hf09Bd02A5GOUqlFm28Uu6b5STwQoIcu4z0.J4UVmxy", // plain: "something"
	}, restrict)

	req, err := http.NewRequest("GET", "/health-check", nil)
	assert.Nil(t, err)
	req.SetBasicAuth(login, password)

	rr := httptest.NewRecorder()
	handler := middleware(testHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, expectedStatus, rr.Code)
}

func TestNewBasicAuthMiddleware_NoRestrict(t *testing.T) {
	generateTestRun(t, false, "user", "test", http.StatusCreated)
}

func TestNewBasicAuthMiddleware_MissingUser(t *testing.T) {
	generateTestRun(t, true, "no_user", "something", http.StatusUnauthorized)
}

func TestNewBasicAuthMiddleware_InvalidPassword(t *testing.T) {
	generateTestRun(t, true, "user", "non_something", http.StatusUnauthorized)
}

func TestNewBasicAuthMiddleware_SuccessAuth(t *testing.T) {
	generateTestRun(t, true, "user", "something", http.StatusCreated)
}
