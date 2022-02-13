package server

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Serve() {
	fmt.Print("hello vero")
	e := echo.New()
	e.GET("/", IndexHandler)
	e.GET("/hello", HelloHandler)
	e.Logger.Fatal(e.Start(":1323"))
}

// curl http://localhost:1323/
func IndexHandler(c echo.Context) error {
	return c.String(http.StatusOK, "hello vero 🌙\n")
}

// curl http://localhost:1323/hello
func HelloHandler(c echo.Context) error {
	trace("")
	return c.String(http.StatusOK, "hello hello 👋\n")
}

func trace(string) {
	// no-op
}
