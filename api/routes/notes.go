package routes

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	cfg "github.com/solomonbaez/SB-Go-NAPI/api/config"
	"github.com/solomonbaez/SB-Go-NAPI/api/models"
)

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
		response := "Failed to fetch notes"
		log.Error().
			Str("Error", e.Error()).
			Msg(response)

		c.JSON(http.StatusInternalServerError, gin.H{"error": response})
		return
	}

	defer rows.Close()

	var notes []models.Note
	for rows.Next() {
		var note models.Note
		if e := rows.Scan(&note.ID, &note.Title, &note.Content); e != nil {
			response := fmt.Sprintf("Failed to fetch note (ID: %v)", note.ID)
			log.Error().
				Str("Error", e.Error()).
				Msg(response)

			c.JSON(http.StatusInternalServerError, gin.H{"error": response})
			return
		}
		notes = append(notes, note)
	}

	c.JSON(http.StatusOK, notes)
}

func (rh *RouteHandler) GetNote(c *gin.Context) {
	id := c.Param("id")
	if _, e := strconv.Atoi(id); e != nil {
		response := fmt.Sprintf("Invalid ID format (ID: %v)", id)
		log.Error().
			Str("Error", e.Error()).
			Msg(response)

		c.JSON(http.StatusBadRequest, gin.H{"error": response})
		return
	}

	var note models.Note

	// populate note item if query is successful
	e := rh.DB.QueryRow(
		"SELECT * FROM notes WHERE id = ?", id,
	).Scan(&note.ID, &note.Title, &note.Content)
	if e != nil {
		response := fmt.Sprintf("Failed to fetch note (ID: %v)", id)
		log.Error().
			Str("Error", e.Error()).
			Msg(response)

		c.JSON(http.StatusNotFound, gin.H{"error": response})
		return
	}

	c.JSON(http.StatusOK, note)
}

func (rh *RouteHandler) PostNote(c *gin.Context) {
	select {
	case <-cfg.Limiter:
	default:
		c.Header("Retry-After", strconv.Itoa(cfg.RATELIMIT)) // automatic retry
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
		return
	}

	var note models.Note
	if e := c.ShouldBindJSON(&note); e != nil {
		response := "Failed to parse JSON"
		log.Error().
			Str("Error", e.Error()).
			Msg(response)

		c.JSON(http.StatusBadRequest, gin.H{"error": response})
		return
	}

	if e := validateInputs(note.Title, note.Content); e != nil {
		response := "Failed to validate data"
		log.Error().
			Str("Error", e.Error()).
			Msg(response)

		c.JSON(http.StatusBadRequest, gin.H{"error": response})
	}

	_, e := rh.DB.Exec(
		"INSERT INTO notes (title, content) VALUES (?, ?)",
		note.Title, note.Content,
	)
	if e != nil {
		response := "Failed to insert data"
		log.Error().
			Str("Error", e.Error()).
			Msg(response)

		c.JSON(http.StatusBadRequest, gin.H{"error": response})
		return
	}

	log.Info().
		Msg("Note created")
	c.JSON(http.StatusCreated, note)
}

func (rh *RouteHandler) UpdateNote(c *gin.Context) {
	select {
	case <-cfg.Limiter:
	default:
		c.Header("Retry-After", strconv.Itoa(cfg.RATELIMIT))
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
		return
	}

	id := c.Param("id")

	var existing_note models.Note
	e := rh.DB.QueryRow("SELECT * FROM notes WHERE id = ?", id).Scan(&existing_note.ID, &existing_note.Title, &existing_note.Content)
	if e != nil {
		response := fmt.Sprintf("Failed to fetch note (ID: %v)", id)
		log.Error().
			Str("Error", e.Error()).
			Msg(response)

		c.JSON(http.StatusInternalServerError, gin.H{"error": response})
		return
	}

	var updated_note models.Note
	e = c.ShouldBindJSON(&updated_note)
	if e != nil {
		response := fmt.Sprintf("Failed to bind data (ID: %v)", id)
		log.Error().
			Str("Error", e.Error()).
			Msg(response)

		c.JSON(http.StatusInternalServerError, gin.H{"error": response})
		return
	}

	if e := validateInputs(updated_note.Title, updated_note.Content); e != nil {
		response := "Failed to validate data"
		log.Error().
			Str("Error", e.Error()).
			Msg(response)

		c.JSON(http.StatusBadRequest, gin.H{"error": response})
	}

	_, e = rh.DB.Exec(
		"UPDATE notes SET title = ?, content = ? WHERE id = ?",
		updated_note.Title, updated_note.Content, id,
	)
	if e != nil {
		response := "Failed to insert data"
		log.Error().
			Str("Error", e.Error()).
			Msg(response)

		c.JSON(http.StatusBadRequest, gin.H{"error": response})
		return
	}

	log.Info().
		Str("ID", id).
		Msg("Note updated")

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
			response := fmt.Sprintf("Note not found (ID: %v)", id)
			log.Error().
				Str("Error", e.Error()).
				Msg(response)

			c.JSON(http.StatusBadRequest, gin.H{"error": response})
		} else {
			response := fmt.Sprintf("Failed to delete note (ID: %v)", id)
			log.Error().
				Str("Error", e.Error()).
				Msg(response)

			c.JSON(http.StatusBadRequest, gin.H{"error": response})
		}
		return
	}

	response := fmt.Sprintf("Note %v deleted", id)
	log.Info().
		Str("ID", id).
		Msg(response)

	c.JSON(http.StatusOK, gin.H{"message": response})
}

// utility
func validateInputs(title string, content string) error {
	if len(title) > cfg.MaxTitleLength {
		return errors.New(
			fmt.Sprintf("Title exceeds maximum length of: %v characters", cfg.MaxTitleLength),
		)
	} else if len(content) > cfg.MaxContentLength {
		return errors.New(
			fmt.Sprintf("Content exceeds maximum length of: %v characters", cfg.MaxTitleLength),
		)
	}

	return nil
}
