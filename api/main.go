package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Note struct {
	ID   int    `json:"id"`
	Text string `json:"note"`
}

var notes []Note
var id_counter = 0

func main() {
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, "hEaLtHy!")
	})

	router.GET("/notes", func(c *gin.Context) {
		c.JSON(http.StatusOK, notes)
	})

	router.POST("/notes", func(c *gin.Context) {
		var note Note
		if e := c.ShouldBindJSON(&note); e != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": e.Error()})
			return
		}

		note.ID = id_counter
		id_counter++
		notes = append(notes, note)
		c.Status(http.StatusCreated)
	})

	router.Run(":8000")
}
