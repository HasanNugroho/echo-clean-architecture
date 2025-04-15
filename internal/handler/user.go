package handler

import (
	"net/http"

	"github.com/HasanNugroho/golang-starter/internal/errs"
	"github.com/HasanNugroho/golang-starter/internal/helper"
	"github.com/HasanNugroho/golang-starter/internal/model"
	"github.com/HasanNugroho/golang-starter/internal/service"
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

// CreateUser godoc
// @Summary      Create an user
// @Description  Create an user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body  UserCreateModel  true  "User Data"
// @Success      201  {object}  shared.Response
// @Failure      400  {object}  shared.Response
// @Failure      404  {object}  shared.Response
// @Failure      500  {object}  shared.Response
// @Router       /users [post]
// @Security ApiKeyAuth
func (c *UserHandler) Create(ctx echo.Context) error {
	var user model.CreateUserRequest
	ctx.Bind(&user)

	if err := c.validate.Struct(user); err != nil {
		return errs.BadRequest("bad request", err)
	}

	if err := c.userService.Create(ctx.Request().Context(), &user); err != nil {
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
// @Success      200     {object}  shared.Response{data=shared.DataWithPagination{items=[]UserModelResponse}}
// @Failure      500     {object}  shared.Response
// @Router       /users [get]
// @Security ApiKeyAuth
func (c *UserHandler) FindAll(ctx echo.Context) error {
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
// @Success      200     {object}  shared.Response{data=UserModel}
// @Failure      500     {object}  shared.Response
// @Router       /users/{id} [get]
// @Security ApiKeyAuth
func (c *UserHandler) FindById(ctx echo.Context) error {
	id := ctx.Param("id")

	if err := c.validate.Var(id, "required"); err != nil {
		return errs.BadRequest("bad request", err)
	}

	user, err := c.userService.FindById(ctx.Request().Context(), id)
	if err != nil {
		return err
	}

	helper.SendSuccess(ctx, http.StatusOK, "User retrieved successfully", user)
	return nil
}

// UpdateUser godoc
// @Summary      Update user
// @Description  Update user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param id path string true "id"
// @Param        user  body  UserUpdateModel  true  "User Data"
// @Success      201  {object}  shared.Response
// @Failure      400  {object}  shared.Response
// @Failure      404  {object}  shared.Response
// @Failure      500  {object}  shared.Response
// @Router       /users/{id} [put]
// @Security ApiKeyAuth
func (c *UserHandler) Update(ctx echo.Context) error {
	id := ctx.Param("id")
	var user model.UpdateUserRequest

	ctx.Bind(&user)

	if err := c.validate.Var(id, "required"); err != nil {
		return errs.BadRequest("bad request", err)
	}

	if err := c.validate.Struct(user); err != nil {
		return errs.BadRequest("bad request", err)
	}

	if err := c.userService.Update(ctx.Request().Context(), id, &user); err != nil {
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
// @Success      200     {object}  shared.Response
// @Failure      500     {object}  shared.Response
// @Router       /users/{id} [delete]
// @Security ApiKeyAuth
func (c *UserHandler) Delete(ctx echo.Context) error {
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
