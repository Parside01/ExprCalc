package middleware

import (
	"github.com/labstack/echo/v4"
)

func CorseDisable() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Add("Access-Control-Allow-Origin", "*")
			c.Response().Header().Add("Access-Control-Allow-Methods", "DELETE, POST, GET, OPTIONS")
			c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
			return next(c)
		}
	}
}
