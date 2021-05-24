package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

const resetApp = true

// rest api server settings
const serverPort = "8080"

// api key - keep value <= 50
const keyLength = 50

// db settings
var (
	dbHostname   string
	dbPort       string
	dbUsername   string
	dbPassword   string
	dbDatabase   string
	dbConnection string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	} else {
		dbHostname = os.Getenv("MYSQL_HOSTNAME")
		dbPort = os.Getenv("MYSQL_PORT")
		dbUsername = os.Getenv("MYSQL_USERNAME")
		dbPassword = os.Getenv("MYSQL_PASSWORD")
		dbDatabase = os.Getenv("MYSQL_DATABASE")

		dbConnection = dbUsername + ":" + dbPassword + "@tcp(" + dbHostname + ":" + dbPort + ")/" + dbDatabase
	}
}
