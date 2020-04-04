package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"

	"github.com/cuirixin/phoenix_gocore/middleware/gin/revision"
)

func rootHandler(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"Hello": "world",
	})
}

func main() {
	port := os.Getenv("PORT")
	router := gin.Default()
	if port == "" {
		port = "8000"
	}
	router.Use(revision.Middleware())
	router.GET("/", rootHandler)
	router.Run(":" + port)
}
