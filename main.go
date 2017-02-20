package main

import (
	"fmt"
	"net/http"

	"os"

	"net/url"

	"io/ioutil"

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
	port := ":1323"
	if p := os.Getenv("PORT"); p != "" {
		port = ":" + p
	}
	e.Logger.Fatal(e.Start(port))
}

func linkedinCallback(c echo.Context) error {
	m := map[string]interface{}{}
	err := c.Bind(&m)
	if err != nil {
		return err
	}

	u, _ := url.ParseRequestURI(m["uri"].(string))
	code := u.Query().Get("code")
	params := `grant_type=authorization_code&code=` + code + `&redirect_uri=https%3A%2F%2Flinkedincallback.herokuapp.com%2Fauth%2Flinkedin%2Fcallback&client_id=81fz2e3avl91e1&client_secret=eT7BJdFihOW1gtvA`

	res, err := http.Post("https://www.linkedin.com/oauth/v2/accessToken?"+params, "application/x-www-form-urlencoded", nil)
	if err != nil {
		return err
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(b))

	c.JSON(http.StatusOK, m)
	return nil
}
