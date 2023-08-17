package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/wehmoen/logart/database"
	"log"
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
				log.Printf("Path %s not allowed", c.Path())
				return echo.ErrBadRequest
			}

			apiKey := c.Request().Header.Get(HeaderXAPIKey)
			db := database.GetDatabaseFromContext(c)
			if db == nil {
				log.Printf("Could not get database from context")
				return echo.ErrInternalServerError
			}

			user, err := db.UserByApiKey(apiKey)

			if err != nil {
				log.Printf("Could not get user by api key (%s): %s", apiKey, err.Error())
				return echo.ErrUnauthorized
			}

			if user == nil {
				log.Printf("User not found")
				return echo.ErrUnauthorized
			}

			c.Set("user", user)
			return next(c)
		}
	}
}
