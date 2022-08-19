package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func main() {
	e := echo.New()

	e.GET("/callback", func(c echo.Context) error {
		c.Set("A", "1")
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return nil
	})

	e.GET("/login", func(c echo.Context) error {
		a := c.Get("A")
		c.String(http.StatusOK, a.(string))
		return nil
	})

	e.Logger.Fatal(e.Start(":3000"))
}
