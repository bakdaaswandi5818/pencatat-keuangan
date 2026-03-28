package handler

import (
	"crypto/subtle"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// APIKeyAuthMiddleware protects API endpoints using a bearer API key.
// /health is intentionally left open for infrastructure health checks.
func APIKeyAuthMiddleware(apiKey string) echo.MiddlewareFunc {
	expected := []byte(apiKey)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			path := c.Path()
			if path == "" {
				path = c.Request().URL.Path
			}
			if path == "/health" {
				return next(c)
			}

			const bearerPrefix = "Bearer "
			authHeader := c.Request().Header.Get(echo.HeaderAuthorization)
			if !strings.HasPrefix(authHeader, bearerPrefix) {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing or invalid authorization header")
			}

			token := strings.TrimSpace(strings.TrimPrefix(authHeader, bearerPrefix))
			if subtle.ConstantTimeCompare([]byte(token), expected) != 1 {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid api key")
			}

			return next(c)
		}
	}
}
