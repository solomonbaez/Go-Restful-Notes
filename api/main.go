package main

import (
	"database/sql"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/solomonbaez/SB-Go-NAPI/api/routes"

	cfg "github.com/solomonbaez/SB-Go-NAPI/api/config"
)

var db *sql.DB

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
	rh := routes.NewRouteHandler(db)
	router := initializeRouter(rh)
	router.Run(":8000")
}

func initializeDatabase() (*sql.DB, error) {
	db, e := sql.Open("mysql", cfg.DB.FormatDSN())
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

func initializeRouter(rh *routes.RouteHandler) *gin.Engine {
	router := gin.Default()
	router.Use(cors.New(cfg.CORS))

	router.POST("/notes", rh.PostNote)
	router.GET("/notes", rh.GetNotes)
	router.PUT("/notes/:id", rh.UpdateNote)
	router.GET("/notes/:id", rh.GetNote)
	router.DELETE("/notes/:id", rh.DeleteNote)

	return router
}
