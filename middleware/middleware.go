package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/wehmoen/logart/database"
	"net/http"
)

const (
	HeaderXAPIKey = "X-API-Key"
)

func ValidateRequest() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			if c.Request().Method != http.MethodPost {
				return echo.ErrMethodNotAllowed
			}

			if c.Path() != "/log" {
				return echo.ErrBadRequest
			}

			apiKey := c.Request().Header.Get(HeaderXAPIKey)
			db := database.GetDatabaseFromContext(c)
			if db == nil {
				return echo.ErrInternalServerError
			}

			user, err := db.UserByApiKey(apiKey)

			if err != nil {
				return echo.ErrUnauthorized
			}

			if user == nil {

				return echo.ErrUnauthorized
			}

			c.Set("user", user)
			return next(c)
		}
	}
}
