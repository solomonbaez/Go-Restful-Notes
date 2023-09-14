package routes

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/solomonbaez/SB-Go-NAPI/api/models"
)

// TODO: utilize cfg.yml structure to centralize globals
const (
	MaxTitleLength   = 100
	MaxContentLength = 1000
)
const DBLIMIT = 1

var limiter = time.Tick(DBLIMIT * time.Second)

// instance
type RouteHandler struct {
	DB *sql.DB
}

// constructor
func NewRouteHandler(db *sql.DB) *RouteHandler {
	return &RouteHandler{
		DB: db,
	}
}

func (rh *RouteHandler) GetNotes(c *gin.Context) {
	rows, e := rh.DB.Query("SELECT * FROM notes")
	if e != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notes"})
		return
	}

	defer rows.Close()

	var notes []models.Note
	for rows.Next() {
		var note models.Note
		if e := rows.Scan(&note.ID, &note.Title, &note.Content); e != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": e.Error()})
			return
		}
		notes = append(notes, note)
	}

	c.JSON(http.StatusOK, notes)
}

func (rh *RouteHandler) GetNote(c *gin.Context) {
	id := c.Param("id")
	if _, e := strconv.Atoi(id); e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var note models.Note

	// populate note item if query is successful
	e := rh.DB.QueryRow(
		"SELECT * FROM notes WHERE id = ?", id,
	).Scan(&note.ID, &note.Title, &note.Content)
	if e != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": e.Error()})
		return
	}

	c.JSON(http.StatusOK, note)
}

func (rh *RouteHandler) PostNote(c *gin.Context) {
	select {
	case <-limiter:
	default:
		c.Header("Retry-After", strconv.Itoa(DBLIMIT)) // automatic retry
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
		return
	}

	var note models.Note
	if e := c.ShouldBindJSON(&note); e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": e.Error()})
		return
	}

	if e := validateInputs(note.Title, note.Content); e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": e.Error()})
		return
	}

	_, e := rh.DB.Exec(
		"INSERT INTO notes (title, content) VALUES (?, ?)",
		note.Title, note.Content,
	)
	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": e.Error()})
		return
	}

	c.JSON(http.StatusCreated, note)
}

func (rh *RouteHandler) UpdateNote(c *gin.Context) {
	select {
	case <-limiter:
	default:
		c.Header("Retry-After", strconv.Itoa(DBLIMIT))
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
		return
	}

	id := c.Param("id")

	var existing_note models.Note
	e := rh.DB.QueryRow("SELECT * FROM notes WHERE id = ?", id).Scan(&existing_note.ID, &existing_note.Title, &existing_note.Content)
	if e != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}

	var updated_note models.Note
	e = c.ShouldBindJSON(&updated_note)
	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to bind data"})
		return
	}

	if e := validateInputs(updated_note.Title, updated_note.Content); e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": e.Error()})
	}

	_, e = rh.DB.Exec(
		"UPDATE notes SET title = ?, content = ? WHERE id = ?",
		updated_note.Title, updated_note.Content, id,
	)
	if e != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update note"})
		return
	}

	c.JSON(http.StatusOK, updated_note)
}

func (rh *RouteHandler) DeleteNote(c *gin.Context) {
	id := c.Param("id")

	_, e := rh.DB.Exec(
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

// utility
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
