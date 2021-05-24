package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var myDb *sql.DB

func main() {
	db, err := sql.Open("mysql", dbConnection)
	if err != nil {
		log.Fatal(err)
	}
	myDb = db
	defer myDb.Close()

	if resetApp {
		setupDb()
	}

	printHeader()
	startWebServer()
}
