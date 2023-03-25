package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	app := NewApplication()
	router := gin.Default()

	router.GET("/api/v1/rate-limits/:limiterName", func(context *gin.Context) {
		limiterName := context.Param("limiterName")
		context.JSON(http.StatusOK, gin.H{
			"message": checkLimit(app, limiterName),
		})
	})

	router.Run(":8902")
}

func checkLimit(app *Application, limiterName string) bool {
	bucket, ok := app.Buckets[limiterName]
	if ok {
		token := bucket.GetToken()
		return token >= 0
	} else {
		log.Printf("Cannot retreive bucket by name '%s'", limiterName)
		return false
	}
}
