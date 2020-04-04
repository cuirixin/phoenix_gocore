package main

import (
	"github.com/gin-gonic/gin"

	"github.com/cuirixin/phoenix_gocore/middleware/gin/location"
)

func main() {
	router := gin.Default()

	// configure to automatically detect scheme and host with
	// fallback to https://foo.com/base
	// - use https when default scheme cannot be determined
	// - use foo.com when default host cannot be determined
	// - include /base as the path
	router.Use(location.New(location.Config{
		Scheme: "https",
		Host:   "foo.com",
		Base:   "/base",
		Headers: location.Headers{Scheme: "X-Forwarded-Proto", Host: "X-Forwarded-For"},
	}))

	router.GET("/", func(c *gin.Context) {
		url := location.Get(c)
		c.String(200, url.String())
	})

	router.Run()
}
