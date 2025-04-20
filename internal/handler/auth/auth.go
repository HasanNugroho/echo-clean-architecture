package handler

import (
	"net/http"

	"github.com/HasanNugroho/golang-starter/internal/errs"
	"github.com/HasanNugroho/golang-starter/internal/helper"
	model "github.com/HasanNugroho/golang-starter/internal/model/auth"
	"github.com/HasanNugroho/golang-starter/internal/service/auth"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService auth.IAuthService
	validate    *validator.Validate
}

func NewAuthHandler(rs auth.IAuthService) *AuthHandler {
	return &AuthHandler{
		authService: rs,
		validate:    validator.New(),
	}
}

// Login godoc
// @Summary      User login
// @Description  Authenticate user and return access token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        login  body  auth.LoginRequest  true  "Login credentials"
// @Success      200  {object}  model.WebResponse
// @Failure      400  {object}  model.WebResponse
// @Failure      401  {object}  model.WebResponse
// @Failure      500  {object}  model.WebResponse
// @Router       /auth/login [post]
// @Security     ApiKeyAuth
func (c *AuthHandler) Login(ctx echo.Context) error {
	var request model.LoginRequest
	if err := ctx.Bind(&request); err != nil {
		return errs.BadRequest("invalid request format", err)
	}

	if err := c.validate.Struct(request); err != nil {
		return errs.BadRequest("validation error", err)
	}

	resp, err := c.authService.Login(ctx.Request().Context(), request)
	if err != nil {
		return err
	}

	helper.SendSuccess(ctx, http.StatusOK, "login successful", resp)
	return nil
}
