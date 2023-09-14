package config

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/go-sql-driver/mysql"
)

// EDIT THESE CONSTANTS
const (
	DBUSER     = "mysql"
	DBPASSWORD = "mysql"
	DBNET      = "tcp"
	DBHOST     = "127.0.0.1:3306"
	DBPORT     = "3306"
	DBNAME     = "notes_api"

	RATELIMIT = 1
)

// INCREASE RATELIMIT to DECREASE requests / second
var Limiter = time.Tick(RATELIMIT * time.Second)

var (
	CORSORIGINS = []string{"http://localhost:3000"} // I'm using NEXT.js
	CORSMETHODS = []string{"GET", "POST", "PUT", "DELETE"}
	CORSHEADERS = []string{"Origin"}
	CORSEXPOSED = []string{"Content-Length"}
	CORSCRD     = true
	CORSMAXAGE  = 1 * time.Hour
)

// DATABASE configuration
var DB = mysql.Config{
	User:   DBUSER,
	Passwd: DBPASSWORD,
	Net:    DBNET,
	Addr:   DBHOST,
	DBName: DBNAME,
}

// CORS configuration
var CORS = cors.Config{
	AllowOrigins:     CORSORIGINS,
	AllowMethods:     CORSMETHODS,
	AllowHeaders:     CORSHEADERS,
	ExposeHeaders:    []string{"Content-Length"},
	AllowCredentials: true,
	MaxAge:           1 * time.Hour,
}
