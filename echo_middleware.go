package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func main() {
	e := echo.New()

	e.GET("/hello", func(c echo.Context) error {
		if c.Get("A") != nil {
			c.String(http.StatusOK, "all good")
		} else {
			c.String(http.StatusOK, "all not good")
		}
		return nil
	})

	e.Use(mw1)

	e.Logger.Fatal(e.Start(":3000"))
}

func mw1(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Set("A", "1")
		return next(c)
	}
}
