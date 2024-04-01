package utils

import (
	_ "database/sql"
	_ "fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var DB *sqlx.DB

func init() {
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "postgres"
	}
	usr := os.Getenv("DB_USER")
	if usr == "" {
		usr = "test"
	}
	pass := os.Getenv("DB_PASS")
	if pass == "" {
		pass = "test"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}
	uri := "user=" + usr + " dbname=" + dbname + " password=" + pass + " port=" + port + " sslmode=disable"
	db, err := sqlx.Connect("postgres", uri)

	if err != nil {
		log.Fatalln(err)
		panic(err)
	}
	DB = db
}
