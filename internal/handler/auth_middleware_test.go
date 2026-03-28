package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestAPIKeyAuthMiddleware_HealthIsPublic(t *testing.T) {
	e := echo.New()
	e.Use(APIKeyAuthMiddleware("secret"))
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
	e.Use(APIKeyAuthMiddleware("secret"))
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
	e.Use(APIKeyAuthMiddleware("secret"))
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
	e.Use(APIKeyAuthMiddleware("secret"))
	e.GET("/summary", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/summary", nil)
	req.Header.Set(echo.HeaderAuthorization, "Bearer secret")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}
