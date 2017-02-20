package main

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"

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
	// u, _ := url.ParseRequestURI(m["uri"].(string))
	// code := u.Query().Get("code")
	// params := `grant_type=authorization_code&code=` + code + `&redirect_uri=https%3A%2F%2Flinkedincallback.herokuapp.com%2Fauth%2Flinkedin%2Fcallback&client_id=81fz2e3avl91e1&client_secret=eT7BJdFihOW1gtvA`

	// res, err := http.Post("https://www.linkedin.com/oauth/v2/accessToken?"+params, "application/x-www-form-urlencoded", nil)
	// if err != nil {
	// 	return err
	// }

	code := c.FormValue("code")
	c.Logger().Print(code)
	apiUrl := "https://www.linkedin.com"
	resource := "/oauth/v2/accessToken/"
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Add("code", code)
	data.Add("redirect_uri", "https%3A%2F%2Flinkedincallback.herokuapp.com%2Fauth%2Flinkedin%2Fcallback")
	data.Add("client_id", os.Getenv("CLIENTID"))
	data.Add("client_secret", os.Getenv("CLIENTSECRET"))
	c.Logger().Print(os.Getenv("CLIENTID"))
	c.Logger().Print(os.Getenv("CLIENTSECRET"))

	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	urlStr := fmt.Sprintf("%v", u) // "https://api.com/user/"

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode())) // <-- URL-encoded payload
	// r.Header.Add("Authorization", "auth_token=\"XXXXXXX\"")
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	res, err := client.Do(r)
	if err != nil {
		return err
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	c.String(res.StatusCode, string(b))
	return nil
}
