package middleware

import (
	"net/http"

	"github.com/HasanNugroho/golang-starter/internal/errs"
	"github.com/HasanNugroho/golang-starter/internal/helper"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func ErrorHandler() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					log.Error().Msgf("Recovered from panic: %v", r)
					helper.SendError(c, http.StatusInternalServerError, "Internal Server Error", nil)
				}
			}()

			err := next(c)
			if err == nil {
				return nil
			}

			if customErr, ok := err.(*errs.CustomError); ok {
				helper.SendError(c, customErr.StatusCode(), customErr.MessageText(), nil)
				return nil
			}

			log.Error().Err(err).Msg("Unhandled error")
			helper.SendError(c, http.StatusInternalServerError, "Internal Server Error", nil)
			return nil
		}
	}
}
