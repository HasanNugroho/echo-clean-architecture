package helper

import (
	"github.com/HasanNugroho/golang-starter/internal/model"
	"github.com/labstack/echo/v4"
)

// SendSuccess mengirim response sukses
func SendSuccess(c echo.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, model.WebResponse{
		Status:  statusCode,
		Message: message,
		Data:    data,
	})
}

// SendError mengirim response error
func SendError(c echo.Context, statusCode int, message string, err interface{}) {
	c.JSON(statusCode, model.WebResponse{
		Status:  statusCode,
		Message: message,
		Data:    err,
	})
}
