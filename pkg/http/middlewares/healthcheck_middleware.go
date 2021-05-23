package middlewares

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func HealthCheckMiddleware(path string, defaultResponse string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if c.Request().URL.Path == path {
				return c.JSON(http.StatusOK, defaultResponse)
			}
			return next(c)
		}
	}
}

