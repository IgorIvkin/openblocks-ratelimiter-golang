package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	app := NewApplication()
	router := gin.Default()

	router.GET("/api/v1/rate-limits/:limiterName", func(context *gin.Context) {
		ProcessCheckLimit(app, context)
	})

	router.Run(":8902")
}
