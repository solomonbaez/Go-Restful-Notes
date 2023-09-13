package main

import (
	"database/sql"
	// "fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	// "github.com/go-sql-driver/mysql"
)

// db -> notes_api
const (
	DBUSER     = "mysql"
	DBPASSWORD = "mysql"
	DBHOST     = "localhost"
	DBPORT     = "3306"
	DBNAME     = "notes_api"
)

type Note struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Text  string `json:"content"`
}

var db *sql.DB

func main() {
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, "hEaLtHy!")
	})

	router.GET("/notes", func(c *gin.Context) {
		c.JSON(http.StatusOK, notes)
	})

	router.GET("/notes/:id", func(c *gin.Context) {
		id_str := c.Param("id")
		id, e := strconv.Atoi(id_str)
		if e != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		for _, note := range notes {
			if note.ID == id {
				c.JSON(http.StatusOK, note)
				return
			}
		}

		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
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
