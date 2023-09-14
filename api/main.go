package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/solomonbaez/SB-Go-NAPI/api/routes"
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

var db *sql.DB

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
	rh := routes.NewRouteHandler(db)
	router := initializeRouter(rh)
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

func initializeRouter(rh *routes.RouteHandler) *gin.Engine {
	router := gin.Default()
	router.Use(cors.New(cors_cfg))

	router.POST("/notes", rh.PostNote)
	router.GET("/notes", rh.GetNotes)
	router.PUT("/notes/:id", rh.UpdateNote)
	router.GET("/notes/:id", rh.GetNote)
	router.DELETE("/notes/:id", rh.DeleteNote)

	return router
}
