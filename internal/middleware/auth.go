package middleware

import (
	"github.com/HasanNugroho/golang-starter/internal/errs"
	"github.com/HasanNugroho/golang-starter/internal/helper"
	"github.com/HasanNugroho/golang-starter/internal/service/account"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type AuthMiddleware struct {
	userService account.IUserService
	logger      *zerolog.Logger
}

func NewAuthMiddleware(logger *zerolog.Logger, userService account.IUserService) *AuthMiddleware {
	return &AuthMiddleware{userService: userService, logger: logger}
}

func (m *AuthMiddleware) AuthRequired() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenString := c.Request().Header.Get("Authorization")
			if tokenString == "" {
				m.logger.Error().Msg("missing authorization header")
				return errs.Unauthorized("Unauthorized", nil)
			}

			claims, err := helper.ParseToken(tokenString)
			if err != nil {
				m.logger.Error().Err(err).Msg("invalid or expired token")
				return errs.Unauthorized("Unauthorized", err)
			}

			data, ok := claims["data"].(map[string]interface{})
			if !ok {
				m.logger.Error().Msg("invalid token payload")
				return errs.Unauthorized("Unauthorized", nil)
			}

			userID, ok := data["user_id"].(string)
			if !ok {
				m.logger.Error().Msg("invalid token payload")
				return errs.Unauthorized("Unauthorized", nil)
			}

			// Ambil informasi IP dan Device (User-Agent) dari request
			ipAddress := c.Request().RemoteAddr
			if forwardedIP := c.Request().Header.Get("X-Forwarded-For"); forwardedIP != "" {
				ipAddress = forwardedIP
			}

			device := c.Request().Header.Get("User-Agent")

			user, err := m.userService.FindById(c.Request().Context(), userID)
			if err != nil {
				m.logger.Error().Err(err).Str("user_id", userID).Str("ip_address", ipAddress).Str("device", device).Msg("user not found")
				return errs.Unauthorized("Unauthorized", err)
			}

			// m.logger.Info().
			// 	Str("user_id", userResponse.ID).
			// 	Str("ip_address", ipAddress).
			// 	Str("device", device).
			// 	Msg("User access successfully")

			c.Set("user", user)

			return next(c)
		}
	}
}
