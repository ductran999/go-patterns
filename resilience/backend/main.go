package main

import "github.com/gin-gonic/gin"

func main() {
	router := gin.Default()
	router.GET("/v1/chat/completions", func(c *gin.Context) {
		c.JSON(500, gin.H{
			"message": "model is down",
		})
	})

	router.Run("localhost:8080")
}
