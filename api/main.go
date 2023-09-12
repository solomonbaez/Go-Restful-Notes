package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()

	var notes []string

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, "hEaLtHy!")
	})

	router.GET("/notes", func(c *gin.Context) {
		c.JSON(http.StatusOK, notes)
	})

	router.Run(":8000")
}
