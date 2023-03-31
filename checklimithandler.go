package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Обрабатывает логику запроса по проверке доступности ресурса
// через рейт лимитер с именем из параметра "limiterName"
func ProcessCheckLimit(app *Application, context *gin.Context) {
	limiterName := context.Param("limiterName")
	bucket, ok := app.Buckets[limiterName]
	if ok {
		token := bucket.GetToken()
		context.JSON(http.StatusOK, token >= 0)
	} else {
		context.JSON(http.StatusNotFound, gin.H{
			"message": fmt.Sprintf("Cannot retreive bucket by name '%s'", limiterName),
		})
	}
}
