package main

import (
	"Gee/gee"
	"net/http"
)

func main() {
	r := gee.New()
	r.Static("/assets", "./static")
	r.Use(gee.Logger())
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "css.tmpl", nil)
	})

	r.Run(":9999")
}
