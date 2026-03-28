package handler

import (
	"crypto/subtle"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

// APIKeyAuthMiddleware protects API endpoints using a bearer API key.
// /health is intentionally left open for infrastructure health checks.
func APIKeyAuthMiddleware(apiKey string) echo.MiddlewareFunc {
	return apiKeyAuthMiddlewareWithClock(apiKey, time.Now)
}

func apiKeyAuthMiddlewareWithClock(apiKey string, now func() time.Time) echo.MiddlewareFunc {

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
			expected := []byte(dynamicAPIKey(apiKey, now()))
			if subtle.ConstantTimeCompare([]byte(token), expected) != 1 {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid api key")
			}

			return next(c)
		}
	}
}

func dynamicAPIKey(baseAPIKey string, current time.Time) string {
	return baseAPIKey + "-" + current.UTC().Format("20060102")
}
