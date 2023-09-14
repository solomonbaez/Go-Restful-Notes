package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

// db -> notes_api
const (
	DBUSER     = "mysql"
	DBPASSWORD = "mysql"
	DBNET      = "tcp"
	DBHOST     = "127.0.0.1:3306"
	DBPORT     = "3306"
	DBNAME     = "notes_api"
	DBLIMIT    = 1 // rate limit - default: 1 request / second
)

type Note struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

const (
	MaxTitleLength   = 100
	MaxContentLength = 1000
)

var db *sql.DB

var limiter = time.Tick(DBLIMIT * time.Second)

var cfg = mysql.Config{
	User:   DBUSER,
	Passwd: DBPASSWORD,
	Net:    DBNET,
	Addr:   DBHOST,
	DBName: DBNAME,
}

var cors_cfg = cors.Config{
	AllowAllOrigins:  true,
	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
	AllowHeaders:     []string{"Origin"},
	ExposeHeaders:    []string{"Content-Length"},
	AllowCredentials: true,
	MaxAge:           1 * time.Hour,
}

func main() {
	// database
	var e error
	db, e = initializeDatabase()
	if e != nil {
		log.Fatal(e)
	}

	log.Println("Database Connection: Success!")
	defer db.Close()

	// router
	router := gin.Default()
	router.Use(cors.New(cors_cfg))

	router.POST("/notes", postNote)
	router.GET("/notes", getNotes)

	router.PUT("/notes/:id", updateNote)
	router.GET("/notes/:id", getNote)
	router.DELETE("/notes/:id", deleteNote)

	router.Run(":8000")
}

func initializeDatabase() (*sql.DB, error) {
	db, e := sql.Open("mysql", cfg.FormatDSN())
	if e != nil {
		return nil, e
	}

	p := db.Ping()
	if p != nil {
		db.Close()
		return nil, p
	}

	return db, nil
}

func getNotes(c *gin.Context) {
	rows, e := db.Query("SELECT * FROM notes")
	if e != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notes"})
		return
	}

	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var note Note
		if e := rows.Scan(&note.ID, &note.Title, &note.Content); e != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": e.Error()})
			return
		}
		notes = append(notes, note)
	}

	c.JSON(http.StatusOK, notes)
}

func getNote(c *gin.Context) {
	id := c.Param("id")
	if _, e := strconv.Atoi(id); e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var note Note

	// populate note item if query is successful
	e := db.QueryRow(
		"SELECT * FROM notes WHERE id = ?", id,
	).Scan(&note.ID, &note.Title, &note.Content)
	if e != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": e.Error()})
		return
	}

	c.JSON(http.StatusOK, note)
}

func postNote(c *gin.Context) {
	select {
	case <-limiter:
	default:
		c.Header("Retry-After", strconv.Itoa(DBLIMIT)) // automatic retry
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
		return
	}

	var note Note
	if e := c.ShouldBindJSON(&note); e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": e.Error()})
		return
	}

	if e := validateInputs(note.Title, note.Content); e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": e.Error()})
		return
	}

	_, e := db.Exec(
		"INSERT INTO notes (title, content) VALUES (?, ?)",
		note.Title, note.Content,
	)
	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": e.Error()})
		return
	}

	c.JSON(http.StatusCreated, note)
}

func updateNote(c *gin.Context) {
	select {
	case <-limiter:
	default:
		c.Header("Retry-After", strconv.Itoa(DBLIMIT))
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
		return
	}

	id := c.Param("id")

	var existing_note Note
	e := db.QueryRow("SELECT * FROM notes WHERE id = ?", id).Scan(&existing_note.ID, &existing_note.Title, &existing_note.Content)
	if e != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}

	var updated_note Note
	e = c.ShouldBindJSON(&updated_note)
	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to bind data"})
		return
	}

	if e := validateInputs(updated_note.Title, updated_note.Content); e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": e.Error()})
	}

	_, e = db.Exec(
		"UPDATE notes SET title = ?, content = ? WHERE id = ?",
		updated_note.Title, updated_note.Content, id,
	)
	if e != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update note"})
		return
	}

	c.JSON(http.StatusOK, updated_note)
}

func deleteNote(c *gin.Context) {
	id := c.Param("id")

	_, e := db.Exec(
		"DELETE FROM notes WHERE id = ?",
		id,
	)
	if e != nil {
		if strings.Contains(e.Error(), "no rows in result set") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete note"})
		}
		return
	}

	response := fmt.Sprintf("Note %v deleted", id)
	c.JSON(http.StatusOK, gin.H{"message": response})
}

func validateInputs(title string, content string) error {
	if len(title) > MaxTitleLength {
		return errors.New(
			fmt.Sprintf("Title exceeds maximum length of: %v characters", MaxTitleLength),
		)
	} else if len(content) > MaxContentLength {
		return errors.New(
			fmt.Sprintf("Content exceeds maximum length of: %v characters", MaxTitleLength),
		)
	}

	return nil
}
