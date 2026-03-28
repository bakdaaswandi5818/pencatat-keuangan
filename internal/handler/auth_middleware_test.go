package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
)

func TestAPIKeyAuthMiddleware_HealthIsPublic(t *testing.T) {
	e := echo.New()
	e.Use(apiKeyAuthMiddlewareWithClock("secret", fixedNow))
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestAPIKeyAuthMiddleware_RequiresAuthorization(t *testing.T) {
	e := echo.New()
	e.Use(apiKeyAuthMiddlewareWithClock("secret", fixedNow))
	e.GET("/summary", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/summary", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}
}

func TestAPIKeyAuthMiddleware_InvalidToken(t *testing.T) {
	e := echo.New()
	e.Use(apiKeyAuthMiddlewareWithClock("secret", fixedNow))
	e.GET("/summary", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/summary", nil)
	req.Header.Set(echo.HeaderAuthorization, "Bearer wrong")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}
}

func TestAPIKeyAuthMiddleware_ValidToken(t *testing.T) {
	e := echo.New()
	e.Use(apiKeyAuthMiddlewareWithClock("secret", fixedNow))
	e.GET("/summary", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/summary", nil)
	req.Header.Set(echo.HeaderAuthorization, "Bearer secret-20260328")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestDynamicAPIKey(t *testing.T) {
	got := dynamicAPIKey("secret", time.Date(2026, 3, 28, 15, 4, 5, 0, time.FixedZone("WIB", 7*60*60)))
	if got != "secret-20260328" {
		t.Fatalf("expected dynamic key %q, got %q", "secret-20260328", got)
	}
}

func fixedNow() time.Time {
	return time.Date(2026, 3, 28, 0, 0, 0, 0, time.UTC)
}
