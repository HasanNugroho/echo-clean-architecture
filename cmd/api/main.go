package main

import (
	"github.com/HasanNugroho/golang-starter/internal"
	"github.com/labstack/echo/v4"
)

func main() {
	router := echo.New()
	internal.Init(router)
}
