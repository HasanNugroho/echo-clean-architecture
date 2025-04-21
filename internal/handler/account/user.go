package handler

import (
	"net/http"

	"github.com/HasanNugroho/golang-starter/internal/errs"
	"github.com/HasanNugroho/golang-starter/internal/helper"
	"github.com/HasanNugroho/golang-starter/internal/model"
	"github.com/HasanNugroho/golang-starter/internal/model/account"
	service "github.com/HasanNugroho/golang-starter/internal/service/account"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService service.IUserService
	validate    *validator.Validate
}

func NewUserHandler(us service.IUserService) *UserHandler {
	return &UserHandler{
		userService: us,
		validate:    validator.New(),
	}
}

// GetCurrentUser godoc
// @Summary      Get current user
// @Description  Get current authenticated user
// @Tags         users
// @Produce      json
// @Success      200  {object}  model.WebResponse
// @Failure      401  {object}  model.WebResponse
// @Router       /users/me [get]
// @Security     ApiKeyAuth
func (c *UserHandler) GetCurrentUser(ctx echo.Context) error {
	user, ok := ctx.Get("user").(*account.User)
	if !ok || user == nil || !user.IsHasAccess([]string{"users:read"}) {
		return errs.Forbidden("Forbidden", nil)
	}

	resp := user.ToUserResponse()
	helper.SendSuccess(ctx, http.StatusOK, "success", resp)
	return nil
}

// CreateUser godoc
// @Summary      Create an user
// @Description  Create an user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body  account.CreateUserRequest  true  "User Data"
// @Success      201  {object}  model.WebResponse
// @Failure      400  {object}  model.WebResponse
// @Failure      404  {object}  model.WebResponse
// @Failure      500  {object}  model.WebResponse
// @Router       /users [post]
// @Security ApiKeyAuth
func (c *UserHandler) Create(ctx echo.Context) error {
	user, ok := ctx.Get("user").(*account.User)
	if !ok || user == nil || !user.IsHasAccess([]string{"users:create"}) {
		return errs.Forbidden("Forbidden", nil)
	}

	var payload account.CreateUserRequest
	ctx.Bind(&payload)

	if err := c.validate.Struct(payload); err != nil {
		return errs.BadRequest("bad request", err)
	}

	if err := c.userService.Create(ctx.Request().Context(), &payload); err != nil {
		return err
	}

	helper.SendSuccess(ctx, http.StatusCreated, "users created successfully", nil)
	return nil
}

// FindAllUsers godoc
// @Summary      Get all users
// @Description  Retrieve a list of all users
// @Tags         users
// @Accept       json
// @Produce      json
// @Param limit query int false "total data per-page" minimum(1) default(10)
// @Param page query int false "page" minimum(1) default(1)
// @Param search query string false "keyword"
// @Success      200     {object}  model.WebResponse{data=model.DataWithPagination{items=[]account.UserResponse}}
// @Failure      500     {object}  model.WebResponse
// @Router       /users [get]
// @Security ApiKeyAuth
func (c *UserHandler) FindAll(ctx echo.Context) error {
	user, ok := ctx.Get("user").(*account.User)
	if !ok || user == nil || !user.IsHasAccess([]string{"users:read"}) {
		return errs.Forbidden("Forbidden", nil)
	}

	var filter model.PaginationFilter

	// Binding query parameters
	if err := ctx.Bind(&filter); err != nil {
		return errs.BadRequest("bad request", err)
	}

	users, totalItem, err := c.userService.FindAll(ctx.Request().Context(), &filter)
	if err != nil {
		return err
	}

	paginate := helper.BuildPagination(&filter, totalItem)
	result := model.DataWithPagination{
		Items:  users,
		Paging: paginate,
	}

	helper.SendSuccess(ctx, http.StatusOK, "users retrieved successfully", result)
	return nil
}

// FindUser godoc
// @Summary      Get all users
// @Description  Retrieve a user by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param id path string true "id"
// @Success      200     {object}  model.WebResponse{data=account.UserResponse}
// @Failure      500     {object}  model.WebResponse
// @Router       /users/{id} [get]
// @Security ApiKeyAuth
func (c *UserHandler) FindById(ctx echo.Context) error {
	user, ok := ctx.Get("user").(*account.User)
	if !ok || user == nil || !user.IsHasAccess([]string{"users:read"}) {
		return errs.Forbidden("Forbidden", nil)
	}

	id := ctx.Param("id")

	if err := c.validate.Var(id, "required"); err != nil {
		return errs.BadRequest("bad request", err)
	}

	user, err := c.userService.FindById(ctx.Request().Context(), id)
	if err != nil {
		return err
	}

	helper.SendSuccess(ctx, http.StatusOK, "User retrieved successfully", user.ToUserResponse())
	return nil
}

// UpdateUser godoc
// @Summary      Update user
// @Description  Update user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param id path string true "id"
// @Param        user  body  account.UpdateUserRequest  true  "User Data"
// @Success      201  {object}  model.WebResponse
// @Failure      400  {object}  model.WebResponse
// @Failure      404  {object}  model.WebResponse
// @Failure      500  {object}  model.WebResponse
// @Router       /users/{id} [put]
// @Security ApiKeyAuth
func (c *UserHandler) Update(ctx echo.Context) error {
	user, ok := ctx.Get("user").(*account.User)
	if !ok || user == nil || !user.IsHasAccess([]string{"users:update"}) {
		return errs.Forbidden("Forbidden", nil)
	}

	id := ctx.Param("id")
	var payload account.UpdateUserRequest

	ctx.Bind(&payload)

	if err := c.validate.Var(id, "required"); err != nil {
		return errs.BadRequest("bad request", err)
	}

	if err := c.validate.Struct(user); err != nil {
		return errs.BadRequest("bad request", err)
	}

	if err := c.userService.Update(ctx.Request().Context(), id, &payload); err != nil {
		return err
	}

	helper.SendSuccess(ctx, http.StatusOK, "users updated successfully", nil)
	return nil
}

// DeleteUser godoc
// @Summary      Delete user
// @Description  Delete user by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param id path string true "id"
// @Success      200     {object}  model.WebResponse
// @Failure      500     {object}  model.WebResponse
// @Router       /users/{id} [delete]
// @Security ApiKeyAuth
func (c *UserHandler) Delete(ctx echo.Context) error {
	user, ok := ctx.Get("user").(*account.User)
	if !ok || user == nil || !user.IsHasAccess([]string{"users:delete"}) {
		return errs.Forbidden("Forbidden", nil)
	}

	id := ctx.Param("id")

	if err := c.validate.Var(id, "required"); err != nil {
		return errs.BadRequest("bad request", err)
	}

	err := c.userService.Delete(ctx.Request().Context(), id)
	if err != nil {
		return err
	}

	helper.SendSuccess(ctx, http.StatusOK, "user deleted successfully", nil)
	return nil
}
