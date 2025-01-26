package middleware

import (
	"net/http"
	"fmt"

	"github.com/labstack/echo/v4"
)

type ETagValidationConfig struct {
	Skipper func(c echo.Context) bool
}

var DefaultETagValidationConfig = ETagValidationConfig{
	Skipper: func(c echo.Context) bool {
		return false
	},
}

func ETagValidation(config ETagValidationConfig) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = DefaultETagValidationConfig.Skipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			requestETag := c.Request().Header.Get("If-None-Match")
			fmt.Println("Request: ", requestETag)

			err := next(c)
			if err != nil {
				return err
			}

			responseETag := c.Response().Header().Get("etag")
			fmt.Println("Response: ", responseETag)

			if requestETag != "" && requestETag == responseETag {
				fmt.Println("Triggered")

				return c.NoContent(http.StatusNotModified)
			}
			fmt.Println("Not triggered")
			return nil
		}
	}
}
