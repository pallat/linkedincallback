package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", hello)
	e.GET("/auth/linkedin/callback", linkedinCallback)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

func linkedinCallback(c echo.Context) error {
	m := map[string]interface{}{}
	err := c.Bind(&m)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, m)
	return nil
}