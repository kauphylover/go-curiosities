package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func main() {
	e := echo.New()

	routes(e)

	e.Logger.Fatal(e.Start(":3000"))
}

func routes(e *echo.Echo) {
	e.GET("*", h1)
	e.GET("/callback", h2)
}

func h1(c echo.Context) error {
	c.String(http.StatusOK, "h1")
	return nil
}

func h2(c echo.Context) error {
	c.String(http.StatusOK, "h2")
	return nil
}
