package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Note struct {
	Text string `json:"note"`
}

func main() {
	router := gin.Default()

	var notes []string

	router.POST("/notes", func(c *gin.Context) {
		var note Note
		if e := c.ShouldBindJSON(&note); e != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": e.Error()})
			return
		}

		notes = append(notes, note.Text)
		c.Status(http.StatusCreated)
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, "hEaLtHy!")
	})

	router.GET("/notes", func(c *gin.Context) {
		c.JSON(http.StatusOK, notes)
	})

	router.Run(":8000")
}
