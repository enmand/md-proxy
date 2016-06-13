package main

import (
	"github.com/labstack/echo"
	"net/http"
	"strings"
)

func initRoutes(g ...*echo.Group) {
	for _, group := range g {
		group.GET("/meta_data.json", func(c echo.Context) error {
			remoteAddr := strings.Split(c.Request().RemoteAddress(), ":")

			if md, ok := serverMetadata[remoteAddr[0]]; !ok {
				return c.String(http.StatusNotFound, "Server not found")
			} else {
				return c.JSON(http.StatusOK, md)
			}
		})
	}
}
